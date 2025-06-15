package app

type Team struct {
	ID          int    `json:"id"`
	TeamName    string `json:"team_name"`
	TeamCount   int    `json:"team_count"`
	MemberNames string `json:"member_names"`
	SchoolName  string `json:"school_name"`
	Password    string `json:"password"`
}

type AdminUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
