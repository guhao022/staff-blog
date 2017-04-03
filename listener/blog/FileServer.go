package blog

import (
	"fmt"
	"github.com/num5/axiom"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"
)

type FileHandler struct {
	tplPath  string
	savePath string
	ctx      *axiom.Context
}

func newFileHandler(tpl, save string, ctx *axiom.Context) *FileHandler {
	return &FileHandler{
		tplPath:  tpl,
		savePath: save,
		ctx:      ctx,
	}
}

func (fh *FileHandler) Http() {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(fh.tplPath+"/assets/"))))
	http.HandleFunc("/", fh.index)
	http.HandleFunc("/upload", fh.upload)
	fh.ctx.Reply("文件上传服务监听端口：%d", 8800)
	err := http.ListenAndServe(":8800", nil)
	if err != nil {
		fh.ctx.Reply("开启文件上传服务器错误：%s", err.Error())
	}
}

func (fh *FileHandler) index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(fh.tplPath + "/index.html")
	if err != nil {
		fh.ctx.Reply("解析主页模版失败：%s", err)
	}
	err = t.Execute(w, "上传文件")
	if err != nil {
		fh.ctx.Reply("解析主页模版失败：%s", err)
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

		r.ParseMultipartForm(32 << 20) //在使用r.MultipartForm前必须先调用ParseMultipartForm方法，参数为最大缓存

		file, handler, err := r.FormFile("file")
		if err != nil {
			fh.ctx.Reply("未找到上传文件：%s", err)
			return
		}

		filename := handler.Filename

		save := fh.savePath + "/" + filename

		//检查文件是否存在
		if !Exist(fh.savePath) {
			os.MkdirAll(fh.savePath, os.ModePerm)
		} else {
			if Exist(save) {
				fh.ctx.Reply("博客《%s》文件已经存在", filename)
				return
			}
		}

		//结束文件
		of, err := handler.Open()
		if err != nil {
			fh.ctx.Reply("文件处理错误： %s", err)
			return
		}
		defer file.Close()

		//保存文件
		f, err := os.Create(save)
		if err != nil {
			fh.ctx.Reply("创建文件失败： %s", err)
			return
		}
		defer f.Close()
		io.Copy(f, of)

		//获取文件状态信息
		fstat, _ := f.Stat()

		//打印接收信息
		fh.ctx.Reply("上传时间:%s, Size: %dKB,  Name:%s\n", time.Now().Format("2006-01-02 15:04:05"), fstat.Size()/1024, filename)

		w.Write([]byte("1"))
		return
	}
}

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

