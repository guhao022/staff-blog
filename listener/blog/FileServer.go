package blog

import (
	"net/http"
	"os"
	"io"
	"time"
	"fmt"
	"github.com/num5/axiom"
	"html/template"
)

type FileHandler struct {
	tplPath string
	savePath string
	ctx *axiom.Context
}

func newFileHandler(tpl,save string, ctx *axiom.Context) *FileHandler {
	return &FileHandler{
		tplPath: tpl,
		savePath: save,
		ctx: ctx,
	}
}

func (fh *FileHandler) Http() {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(fh.tplPath + "/assets/"))))
	http.HandleFunc("/", fh.index)
	http.HandleFunc("/upload", fh.upload)
	fh.ctx.Reply("监听端口：%d", 8888)
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		fh.ctx.Reply("开启文件上传服务器错误：%s", err.Error())
	}
}

func (fh *FileHandler) index(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles(fh.tplPath + "/index.html")
	t.Execute(w, "上传文件")
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

		file, handler, err := r.FormFile("file")
		if err != nil {
			fh.ctx.Reply("未找到上传文件：%s", err)
			return
		}

		filename := handler.Filename

		save := fh.savePath + "/" +filename

		//检查文件是否存在
		_, err = os.Stat(filename)
		if err != nil && !os.IsExist(err) {
			fh.ctx.Reply("博客文件已经存在： %s", err.Error())
		}

		//结束文件
		of, err := handler.Open()
		if err != nil {
			fh.ctx.Reply("文件处理错误： %s", err.Error())
		}
		defer file.Close()
		//保存文件

		f, err := os.Create(save)
		defer f.Close()
		io.Copy(f, of)

		//获取文件状态信息
		fstat,_ := f.Stat()

		//打印接收信息
		fh.ctx.Reply("%s Size: %d KB  Name：%s\n", time.Now().Format("2006-01-02 15:04:05"), fstat.Size()/1024, filename)

		return
	}
}