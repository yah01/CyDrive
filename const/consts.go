package _const

const(
	ListenPort = ":6454"
	UserDataDir = "user_data"
)

type Command = string
const (
	LOGIN Command = "LOGIN"
	GET   Command = "GET"
	GETA  Command = "GETA"
	SEND  Command = "SEND"
	SENDA Command = "SENDA"
	LIST  Command = "LIST"
	DONE  Command = "DONE"
	RCD   Command = "RCD"

	QUIT Command = "QUIT"

	Delim = '\n'
)

type CdpStatus = int
const (
	StatusOk CdpStatus = iota
	StatusAuthError
	StatusIoError
	StatusInternalError
)
