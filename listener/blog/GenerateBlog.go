package blog

import (
	"github.com/num5/axiom"
	"os/exec"
)

func (b *BlogListener) compileBlog(ctx *axiom.Context) {
	cmd := exec.Command(WORKER_DIR + "/chca", "compile")
	if err := cmd.Start(); err != nil {
		ctx.Reply("博客编译错误：%s", err.Error())
	}

	ctx.Reply("编译成功，请登录 http://" + BLOG_HOST + " 查看")
}

func (b *BlogListener) blogserver(ctx *axiom.Context, port string) {
	cmd := exec.Command(WORKER_DIR + "/chca", "http", port)

	if err := cmd.Start(); err != nil {
		ctx.Reply("开启内置web服务失败：%s", err.Error())
	}

	ctx.Reply("内置web服务开启成功，监听端口：%s", port)

}
