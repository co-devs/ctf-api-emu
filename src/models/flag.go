package models

type Flag struct {
	Identifier	string		`json:"identifier"`
	TeamID		int			`json:"team_id"`
	ServiceID	int			`json:"service_id"`
	Tick		int 		`json:"tick"`
	Expiration	string	`json:"expiration"`
	Hostname	string		`json:"hostname"`
}
