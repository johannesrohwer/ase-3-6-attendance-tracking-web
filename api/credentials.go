package api

type Credentials struct {
	ID               string
	PermissionGroups []string
}

func NewCredentials(ID string, permissionGroups []string) *Credentials {
	return &Credentials{ID: ID, PermissionGroups: permissionGroups}
}
