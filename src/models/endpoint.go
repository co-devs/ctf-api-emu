package models

type Endpoint struct {
	TeamID		int		`json:"team_id"`
	ServiceID	int		`json:"service_id"`
	ServiceName	string	`json:"service_name"`
	Hostname	string	`json:"hostname"`
}
