package model

type Repository struct {
	ID              string          `json:"id"`
	RestAPIID       *int64          `json:"rest_api_id"`
	Name            *string         `json:"name"`
	FullName        *string         `json:"fullName"`
	CollaboratorsID *int            `json:"collaborators"`
}
func (Repository) IsNode() {}

type GraphID struct {
	Typename string
	ID       string
}
