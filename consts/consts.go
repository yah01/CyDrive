package consts

const(
	ListenPort = ":6454"
	UserDataDir = "user_data"
)

type CdpStatus = int
const (
	StatusOk CdpStatus = 0
	StatusAuthError = (1<<iota)/2
	StatusNoParameterError
	StatusSessionError
	StatusFileTooLargeError
	StatusIoError
	StatusInternalError
)

const (
	// The size of file must be not greater than 1GB
	FileSizeLimit int64 = 1<<30
)