package libs

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var emailPattern = regexp.MustCompile("[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?")

func RealPath(filePath string) string {
	return os.ExpandEnv(filePath)
}

func IsFile(filePath string) bool {
	f, e := os.Stat(filePath)
	if e != nil {
		return false
	}
	return !f.IsDir()
}
func IsDir(dir string) bool {
	f, e := os.Stat(dir)
	if e != nil {
		return false
	}
	return f.IsDir()
}

// 版本对比 v1比v2大返回1，小于返回-1，等于返回0
func VerCompare(ver1, ver2 string) int {
	ver1 = strings.TrimLeft(ver1, "ver") // 清除v,e,r
	ver2 = strings.TrimLeft(ver2, "ver") // 清除v,e,r
	p1 := strings.Split(ver1, ".")
	p2 := strings.Split(ver2, ".")

	ver1 = ""
	for _, v := range p1 {
		iv, _ := strconv.Atoi(v)
		ver1 = fmt.Sprintf("%s%04d", ver1, iv)
	}

	ver2 = ""
	for _, v := range p2 {
		iv, _ := strconv.Atoi(v)
		ver2 = fmt.Sprintf("%s%04d", ver2, iv)
	}

	return strings.Compare(ver1, ver2)
}

// 换行符换成<br />
func Nl2br(s string) string {
	s = strings.Replace(s, "\r\n", "\n", -1)
	s = strings.Replace(s, "\r", "\n", -1)
	s = strings.Replace(s, "\n", "<br />", -1)
	return s
}
func IsEmail(b []byte) bool {
	return emailPattern.Match(b)
}
