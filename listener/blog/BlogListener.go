package blog

import (
	"github.com/num5/axiom"
	"os"
	"strings"
	"regexp"
)

type BlogListener struct{}

var (
	WORKER_DIR = os.Getenv("CHCA_WORKER_DIR")
	BLOG_HOST = os.Getenv("BLOG_HOST")
)

func (b *BlogListener) Handle() []*axiom.Listener {

	return []*axiom.Listener{
		{
			// 编译博客
			Regex: "编译博客|博客编译|更新博客|博客更新|编译markdown|编译MARKDOWN|markdown编译|MARKDOWN编译",
			HandlerFunc: func(ctx *axiom.Context) {
				b.compileBlog(ctx)
			},
		},{
			// 开启chca内部webserver
			Regex: "开启博客|开启webserver|开启服务器|打开博客服务器|打开web|打开web服务器",
			HandlerFunc: func(ctx *axiom.Context) {
				var port string = "9900"
				regexp := regexp.MustCompile(`(端口：|端口:|port：|port:)(\d+)`)
				matches := regexp.FindStringSubmatch(ctx.Message.Text)

				if len(matches) >= 3 {
					port = matches[2]
				}
				go func() {
					b.blogserver(ctx, port)
				}()

			},
		},{
			// 更新博客生成器
			Regex: "更新chca|更新博客生成器|下载chca|下载博客生成器",
			HandlerFunc: func(ctx *axiom.Context) {
				var m string
				if strings.Contains(ctx.Matches[0], "更新") {
					m = "更新"
				}
				if strings.Contains(ctx.Matches[0], "下载") {
					m = "下载"
				}
				b.updateChca(ctx, m)
			},
		},{
			// 上传博客
			Regex: "上传博客|上传博客文件",
			HandlerFunc: func(ctx *axiom.Context) {
				UPLOAD_TEMPLATE := os.Getenv("UPLOAD_TEMPLATE")
				go func() {

					fh.Http()
				}()

			},
		},/*{
			Regex: "",
			HandlerFunc: func(ctx *axiom.Context) {
				ctx.Reply("未识别命令，so so so sorry ~ ~ ~ ")
			},
		},*/
	}
}

// 创建文件
func createFile(name string) (*os.File, error) {
	err := os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name, "/")]), 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}
