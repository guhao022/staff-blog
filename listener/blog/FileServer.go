package blog

import (
	"net/http"
	"os"
	"io"
	"time"
	"fmt"
)

type FileServer struct {

}

func indexHandle(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("获取页面失败")
		}
	}()
	http.Request{}

	// 上传页面
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(200)
	html := `
		<html>
	    <head>
	        <title>Golang Upload Files</title>
	    </head>
	    <body>
	        <form id="uploadForm"  enctype="multipart/form-data" action="/upload" method="POST">
	            <p>Golang Upload</p> <br/>
	            <input type="file" id="file1" name="userfile" multiple />	<br/>
	            <input type="submit" value="Upload">
	        </form>
	   	</body>
		</html>`
	io.WriteString(w, html)
}
// 上传文件接口
func upload(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("文件上传异常")
		}
	}()

	if "POST" == r.Method {

		r.ParseMultipartForm(32 << 20)	//在使用r.MultipartForm前必须先调用ParseMultipartForm方法，参数为最大缓存
		// fmt.Println(r.MultipartForm)
		// fmt.Println(r.MultipartReader())
		if r.MultipartForm != nil && r.MultipartForm.File != nil {
			fhs := r.MultipartForm.File["file"]		//获取所有上传文件信息
			num := len(fhs)

			fmt.Printf("总文件数：%d 个文件", num)

			//循环对每个文件进行处理
			for n, fheader := range fhs {
				//获取文件名
				filename := fheader.Filename

				//结束文件
				file,err := fheader.Open()
				if err != nil {
					fmt.Println(err)
				}

				//保存文件
				defer file.Close()
				f, err := os.Create(filename)
				defer f.Close()
				io.Copy(f, file)

				//获取文件状态信息
				fstat,_ := f.Stat()

				//打印接收信息
				fmt.Fprintf(w, "%s  NO.: %d  Size: %d KB  Name：%s\n", time.Now().Format("2006-01-02 15:04:05"), n, fstat.Size()/1024, filename)
				fmt.Printf("%s  NO.: %d  Size: %d KB  Name：%s\n", time.Now().Format("2006-01-02 15:04:05"), n, fstat.Size()/1024, filename)

			}
		} else {
			indexHandle(w, r)
		}

		return
	}
}