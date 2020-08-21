package model

type FileInfo struct {
	FileMode   uint32 `json:"file_mode"`
	ModifyTime int64 `json:"modify_time"`
	FilePath   string `json:"file_path"`
}

type File struct {
	FileInfo
	Data []byte `json:"data"`
	IsCompressed bool `json:"is_compressed"`
}
