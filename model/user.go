package model

type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	RootDir  string `json:"root_dir"`
	WorkDir  string `json:"work_dir"`
}
