package dto

type RoleDto struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}
