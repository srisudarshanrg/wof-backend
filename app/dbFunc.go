package app

import "log"

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
