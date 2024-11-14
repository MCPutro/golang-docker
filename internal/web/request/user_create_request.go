package request

type UserCreate struct {
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Password string `json:"password,omitempty"`
}
