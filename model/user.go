package model

type User struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	FullName     string `json:"full_name"`
	Password     string `json:"password,omitempty"`
	CreationDate string `json:"creation_date"`
}
