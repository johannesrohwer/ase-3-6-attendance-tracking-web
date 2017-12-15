package api

type Credentials struct {
	ID          string   `json:"id"`
	Permissions []string `json:"permissions"`
}

func NewCredentials(ID string, permissions []string) *Credentials {
	return &Credentials{ID: ID, Permissions: permissions}
}
