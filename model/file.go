package model

type FileInfo struct {
	FileMode     uint32 `json:"file_mode"`
	ModifyTime   int64  `json:"modify_time"`
	FilePath     string `json:"file_path"`
	Size         int64  `json:"size"`
	IsDir        bool   `json:"is_dir"`
	IsCompressed bool   `json:"is_compressed"`
}
