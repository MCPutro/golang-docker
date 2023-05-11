package web

type UserResponse struct {
	Id           int    `json:"id,omitempty"`
	Username     string `json:"username,omitempty"`
	Fullname     string `json:"fullname,omitempty"`
	Token        string `json:"token,omitempty"`
	CreationDate string `json:"creation_date,omitempty"`
}
