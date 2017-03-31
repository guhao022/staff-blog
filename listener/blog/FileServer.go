package blog

import (
	"net/http"
	"os"
	"io"
	"time"
	"fmt"
	"github.com/num5/axiom"
	"github.com/robertkrimen/otto/file"
)

type FileHandler struct {
	savePath string
	ctx *axiom.Context
}

func newFileHandler(save string, ctx *axiom.Context) *FileHandler {
	return FileHandler{
		savePath: save,
		ctx: ctx,
	}
}

// 上传文件接口
func (fh *FileHandler) upload(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("文件上传异常")
		}
	}()

	if "POST" == r.Method {

		r.ParseMultipartForm(32 << 20)	//在使用r.MultipartForm前必须先调用ParseMultipartForm方法，参数为最大缓存

		if r.MultipartForm != nil && r.MultipartForm.File != nil {
			fhs := r.MultipartForm.File["file"]		//获取所有上传文件信息

			fh.ctx.Reply("总文件数：%d 个文件", len(fhs))

			//循环对每个文件进行处理
			for n, fheader := range fhs {
				//获取文件名
				filename := fheader.Filename

				save := fh.savePath + "/" +filename

				//检查文件是否存在
				_, err := os.Stat(filename)
				if err != nil && !os.IsExist(err) {
					fh.ctx.Reply("博客文件已经存在： %s", err.Error())
					continue
				}

				//结束文件
				file,err := fheader.Open()
				if err != nil {
					fh.ctx.Reply("文件处理错误： %s", err.Error())
				}

				//保存文件
				defer file.Close()

				f, err := os.Create(save)
				defer f.Close()
				io.Copy(f, file)

				//获取文件状态信息
				fstat,_ := f.Stat()

				//打印接收信息
				fh.ctx.Reply("%s  NO.: %d  Size: %d KB  Name：%s\n", time.Now().Format("2006-01-02 15:04:05"), n, fstat.Size()/1024, filename)

			}
		}

		return
	}
}