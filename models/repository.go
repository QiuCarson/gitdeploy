package models

import "gitdeploy/libs"

type Repository struct{}

// 克隆git仓库
func (this *Repository) CloneRepo(url string, dst string) error {
	out, stderr, err := libs.ExecCmd("git", "clone", url, dst)
	debug("out", out)
	debug("stderr", stderr)
	debug("err", err)
	if err != nil {
		return concatenateError(err, stderr)
	}
	return nil
}
