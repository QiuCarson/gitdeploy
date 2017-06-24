package libs

import "os"

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
