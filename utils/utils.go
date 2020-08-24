package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
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

func ForEachFile(path string, handle func(filename string)) {
	fileinfo, _ := os.Stat(path)

	if !fileinfo.IsDir() {
		file, _ := os.Open(path)
		handle(file.Name())
		return
	}

	fileinfoList, _ := ioutil.ReadDir(path)

	for _, fileinfo = range fileinfoList {
		ForEachFile(filepath.Join(path, fileinfo.Name()), handle)
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
