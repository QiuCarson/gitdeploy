# gitdeploy

go语言写的代码发布系统，程序通过git下载到跳板机，程序必须打上标签才能同步，通过rsync同步到正式服务器上，在执行生成的sh脚本


必须安装rsync<br>
yum -y install rsync xinetd

登录用户<br>
admin<br>
admin888<br>

要发布的服务器必须要跟跳板机授权<br>
[两台服务器ssh授权](https://www.phpsong.com/2169.html)

程序使用的包<br>
go get github.com/astaxie/beego 【beego 框架】<br>
go get golang.org/x/crypto/ssh<br>
go get github.com/pkg/sftp<br>
go get gopkg.in/gomail.v2<br>
