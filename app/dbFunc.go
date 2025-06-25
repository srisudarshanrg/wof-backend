package app

import (
	"database/sql"
	"errors"
	"log"
)

func (app *Application) AddTeam(team Team) (string, error) {
	if team.TeamCount > 5 {
		return "Team can have a maximum of 5 members", nil
	}
	checkTeamExistsQuery := `select * from teams where team_name=$1`
	result, err := app.DB.Exec(checkTeamExistsQuery, team.TeamName)
	if err != nil {
		return "", err
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected > 0 {
		return "This team name already exists, pick another one!", nil
	}

	hashPassword, err := app.HashPassword(team.Password)
	if err != nil {
		return "", err
	}

	addTeamQuery := `insert into teams(team_name, team_count, member_names, school_name, password) values($1, $2, $3, $4, $5)`
	_, err = app.DB.Exec(addTeamQuery, team.TeamName, team.TeamCount, team.MemberNames, team.SchoolName, hashPassword)
	if err != nil {
		return "", err
	}

	return "Your account has been created successfully!", nil
}

func (app *Application) LoginTeam(team Team) (bool, Team, string, error) {
	query := `select * from teams where team_name=$1`
	result, err := app.DB.Exec(query, team.TeamName)
	if err != nil {
		return true, Team{}, "", err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return false, Team{}, "Either team name or password is incorrect", nil
	}

	row := app.DB.QueryRow(query, team.TeamName)
	var teamRecieved Team
	var createdAt, updatedAt interface{}
	err = row.Scan(&teamRecieved.ID, &teamRecieved.TeamName, &teamRecieved.TeamCount, &teamRecieved.MemberNames, &teamRecieved.SchoolName, &teamRecieved.Password, &createdAt, &updatedAt)
	log.Println(teamRecieved)
	if err != nil {
		return true, Team{}, "", err
	}

	err = app.ComparePasswordWithHash(team.Password, teamRecieved.Password)
	if err != nil {
		return false, Team{}, "Either team name or password is incorrect", err
	}

	return true, teamRecieved, "Logged in successfully!", nil
}

func (app *Application) AdminLoginFunc(credential, password string) (AdminUser, error) {
	query := `select * from admin_credentials where username=$1 or email=$1`
	row := app.DB.QueryRow(query, credential)

	var adminUser AdminUser
	var createdAt, updatedAt interface{}

	err := row.Scan(&adminUser.ID, &adminUser.Username, &adminUser.Email, &adminUser.Password, &createdAt, &updatedAt)
	if err != nil && err == sql.ErrNoRows {
		return AdminUser{}, errors.New("username or password is incorrect")
	} else if err != nil {
		return AdminUser{}, err
	}
	check := app.ComparePasswordWithHash(password, adminUser.Password)
	if check != nil {
		return AdminUser{}, errors.New("username or password is incorrect")
	}
	return adminUser, nil
}

func (app *Application) GetDataForAdmin() ([]Team, []ProjectTeamJoin, error) {
	queryGetTeams := `select * from teams`
	rowsTeam, err := app.DB.Query(queryGetTeams)
	if err != nil {
		return []Team{}, []ProjectTeamJoin{}, err
	}
	defer rowsTeam.Close()

	var teams []Team
	for rowsTeam.Next() {
		var team Team
		var createdAt, updatedAt interface{}

		err = rowsTeam.Scan(&team.ID, &team.TeamName, &team.TeamCount, &team.MemberNames, &team.SchoolName, &team.Password, &createdAt, &updatedAt)
		if err != nil {
			return []Team{}, []ProjectTeamJoin{}, err
		}

		teams = append(teams, team)
	}

	queryGetProjects := `select * from projects`
	rowsProjects, err := app.DB.Query(queryGetProjects)
	if err != nil {
		return []Team{}, []ProjectTeamJoin{}, err
	}
	defer rowsProjects.Close()

	var projects []ProjectTeamJoin
	for rowsProjects.Next() {
		var project ProjectTeamJoin
		var createdAt, updatedAt interface{}

		err = rowsProjects.Scan(&project.ID, &project.TeamName, &project.ProjectRepo, &project.ImageLink, &createdAt, &updatedAt)
		if err != nil {
			return []Team{}, []ProjectTeamJoin{}, nil
		}

		getSchoolNameQuery := `select school_name from teams where team_name=$1`
		row := app.DB.QueryRow(getSchoolNameQuery, &project.TeamName)
		err = row.Scan(&project.SchoolName)
		if err != nil {
			return []Team{}, []ProjectTeamJoin{}, err
		}

		projects = append(projects, project)
	}

	return teams, projects, nil
}

func (app *Application) SubmitProject(teamName, projectName, projectRepo, imageLink string) (Project, string, error) {
	queryCheckTeamExists := `select * from teams where team_name=$1`
	result, err := app.DB.Exec(queryCheckTeamExists, teamName)
	if err != nil {
		return Project{}, "", err
	}
	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		return Project{}, "the team you are submitting project for does not exist", nil
	}

	queryCheckProjectExists := `select id from projects where team_name=$1`
	row := app.DB.QueryRow(queryCheckProjectExists, teamName)
	var id int
	err = row.Scan(&id)
	if err != nil && err == sql.ErrNoRows {
		queryAddProject := `insert into projects(team_name, project_name, project_repo, image_link) values($1, $2, $3)`
		_, err = app.DB.Exec(queryAddProject, teamName, projectName, projectRepo, imageLink)

		projectRow := app.DB.QueryRow(`select * from projects where team_name=$1`, teamName)
		var project Project
		var createdAt, updatedAt interface{}
		err = projectRow.Scan(&project.ID, &project.TeamName, &project.ProjectRepo, &project.ImageLink, &createdAt, &updatedAt)
		if err != nil {
			return Project{}, "", err
		}

		return project, "Project added to database", nil
	} else if err != nil && err != sql.ErrNoRows {
		return Project{}, "", err
	}

	queryUpdateProject := `update projects set project_repo=$1, image_link=$2, project_name=$3 where id=$3`
	_, err = app.DB.Exec(queryUpdateProject, projectRepo, imageLink, projectName, id)
	if err != nil {
		return Project{}, "", err
	}

	projectRow := app.DB.QueryRow(`select * from projects where team_name=$1`, teamName)
	var project Project
	var createdAt, updatedAt interface{}
	err = projectRow.Scan(&project.ID, &project.TeamName, &project.ProjectName, &project.ProjectRepo, &project.ImageLink, &createdAt, &updatedAt)
	if err != nil {
		return Project{}, "", err
	}

	return project, "Project updated to database", nil
}

func (app *Application) GetProjects(teamName string) (Project, error) {
	query := `select * from projects where team_name=$1`
	row := app.DB.QueryRow(query, teamName)

	var project Project
	var createdAt, updatedAt interface{}
	err := row.Scan(&project.ID, &project.TeamName, &project.ProjectName, &project.ProjectRepo, &project.ImageLink, &createdAt, &updatedAt)
	if err != nil && err == sql.ErrNoRows {
		return Project{}, errors.New("no projects yet")
	}
	return project, nil
}
