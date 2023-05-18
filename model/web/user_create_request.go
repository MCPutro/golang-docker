package web

type UserCreateRequest struct {
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Password string `json:"password,omitempty"`
}
