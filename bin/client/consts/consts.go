package consts

type Command = string
const (
	LOGIN Command = "LOGIN"
	GET   Command = "GET"
	GETA  Command = "GETA"
	SEND  Command = "SEND"
	SENDA Command = "SENDA"
	LIST  Command = "LIST"
	RCD   Command = "RCD"

	QUIT Command = "QUIT"

	Delim = '\n'
)

var (
	BasePath string
	BaseUrl string
)