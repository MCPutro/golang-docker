package web

type UserCreateRequest struct {
	Username string `json:"username"`
	FullName string `json:"fullname"`
	Password string `json:"password,omitempty"`
}
