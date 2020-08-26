package model

import (
	"github.com/yah01/CyDrive/consts"
	"os"
	"reflect"
)

type FileInfo struct {
	FileMode     uint32 `json:"file_mode"`
	ModifyTime   int64  `json:"modify_time"`
	FilePath     string `json:"file_path"`
	Size         int64  `json:"size"`
	IsDir        bool   `json:"is_dir"`
	IsCompressed bool   `json:"is_compressed"`
}

func NewFileInfo(fileInfo os.FileInfo, path string) FileInfo {
	return FileInfo{
		FileMode:     uint32(fileInfo.Mode()),
		ModifyTime:   fileInfo.ModTime().Unix(),
		FilePath:     path,
		Size:         fileInfo.Size(),
		IsDir:        fileInfo.IsDir(),
		IsCompressed: fileInfo.Size() > consts.CompressBaseline,
	}
}

func NewFileInfoFromMap(infoMap map[string]interface{}) *FileInfo {
	fileInfo := FileInfo{}
	value := reflect.ValueOf(&fileInfo)
	typeOf := reflect.TypeOf(fileInfo)
	for i := 0; i < typeOf.NumField(); i++ {
		field := value.Elem().Field(i)
		tag := infoMap[typeOf.Field(i).Tag.Get("json")]
		newValue := reflect.ValueOf(tag).Convert(field.Type())
		field.Set(newValue)
	}
	return &fileInfo
}
