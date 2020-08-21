package model

import (
	"time"
)

type User struct {
	Id        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	// absolute path
	RootDir   string    `json:"root_dir"`
	// relative path to RootDir
	WorkDir   string    `json:"work_dir"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
