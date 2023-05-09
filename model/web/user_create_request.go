package web

type UserCreateRequest struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Password string `json:"password,omitempty"`
}
