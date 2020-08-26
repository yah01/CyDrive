package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/yah01/CyDrive/consts"
	"github.com/yah01/CyDrive/model"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	TimeFormat = "2006-01-02 15:04:05.999999999 -0700 MST"
)

func Md5Hash(password []byte) []byte {
	md5Value := md5.Sum(password)
	return md5Value[:]
}

func Sha256Hash(password []byte) []byte {
	sha256Value := sha256.Sum256(password)
	return sha256Value[:]
}

func PasswordHash(password string) string {
	bytes := Sha256Hash(Md5Hash([]byte(password)))
	var res string
	for _, v := range bytes {
		res += fmt.Sprint(v)
	}
	return res
}

func PackRange(begin, end int64) string {
	return fmt.Sprintf("bytes=%v-%v", begin, end)
}

func UnpackRange(rangeStr string) (int64, int64) {
	rangeStr = strings.TrimPrefix(rangeStr, "bytes=")
	tuple := strings.Split(rangeStr, "-")
	begin, _ := strconv.ParseInt(tuple[0], 10, 64)
	end, _ := strconv.ParseInt(tuple[1], 10, 64)

	return begin, end
}

func ShouldCompressed(fileInfo os.FileInfo) bool {
	return fileInfo.Size() > consts.CompressBaseline
}

func GetResp(resp *http.Response) *model.Resp {
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	res := model.Resp{}
	if err = json.Unmarshal(bytes, &res); err != nil {
		return nil
	}
	return &res
}

func ForEachFile(path string, handle func(file *os.File)) {
	fileinfo, err := os.Stat(path)
	if err != nil {
		fmt.Println(err, path)
		return
	}
	if !fileinfo.IsDir() {
		file, _ := os.Open(path)
		handle(file)
		file.Close()
		return
	}

	fileinfoList, _ := ioutil.ReadDir(path)

	for _, fileinfo = range fileinfoList {
		ForEachFile(filepath.Join(path, fileinfo.Name()), handle)
	}
}

func ForEachRemoteFile(path string,
	getFileInfo func(path string) *model.FileInfo,
	readDir func(path string) []*model.FileInfo,
	handle func(file *model.FileInfo)) {

	fileInfo := getFileInfo(path)
	if fileInfo == nil {
		fmt.Println("can't get file info:", path)
		return
	}
	if !fileInfo.IsDir {
		handle(fileInfo)
		return
	}

	fileinfoList := readDir(path)

	for _, fileInfo = range fileinfoList {
		ForEachRemoteFile(fileInfo.FilePath, getFileInfo, readDir, handle)
	}
}

func FilterEmptyString(strList []string) []string {
	res := []string{}
	for _, str := range strList {
		if len(str) > 0 {
			res = append(res, str)
		}
	}
	return res
}
