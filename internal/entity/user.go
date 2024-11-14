package entity

type User struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	Fullname     string `json:"fullname"`
	Password     string `json:"password,omitempty"`
	CreationDate string `json:"creation_date,omitempty"`
}
