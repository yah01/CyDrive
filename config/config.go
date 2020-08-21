package config

import "fmt"

type Config struct {
	// "mysql" or "mem"
	UserStoreType string

	// mysql host or json filepath
	Path     string

	User     string
	Password string
	Database string
}

func (config Config) PackDSN() string {
	return fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Path,
		config.Database)
}
