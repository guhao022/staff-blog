package curl

import (
	"fmt"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

// curl.Print 设置
type PrintOps struct {
	Header   bool
	Footer   bool
	LeftEnd  string
	RightEnd string
	Fill     string
	Arrow    string
	Empty    string
}

// 设置PrintOps默认值
var Options = PrintOps{true, true, "[", "]", "=", ">", "_"}

// Print 输出下载开始信息
func header(dl *Download) {
	if Options.Header {
		fmt.Printf("开始下载 [%v]......\n", strings.Join((*dl).GetValues("Title"), ", "))
	}
}

// Print 输出下载结束信息
func footer() {
	if Options.Footer {
		fmt.Println("\r\n下载完成....")
	}
}

// 显示下载情况
func progressbar(title string, start time.Time, i int, suffix string) {
	h := Options.LeftEnd + strings.Repeat(Options.Fill, i) + Options.Arrow + strings.Repeat(Options.Empty, 50-i) + Options.RightEnd
	d := time.Now().Sub(start)
	s := fmt.Sprintf("%v %.0f%% %s %v", safeTitle(title), float32(i)/50*100, h, time.Duration(d.Seconds())*time.Second)
	l := utf8.RuneCountInString(s)
	if l > 80 {
		l = l - 80
	} else {
		l = 80 - l
	}
	e := strings.Repeat(" ", l)
	fmt.Printf("\r%v%v%v", s, e, suffix)
}

func parseArgs(args ...interface{}) (int, Download) {
	dl := Download{}
	if len(args) == 0 {
		panic(CurlError{"curl.New()", -6, "curl.New() parameter type error."})
	} else {
		switch args[0].(type) {
		case string:
			url, title, name, dst := safeArgs(args...)
			dl.AddTask(Task{url, title, name, dst, 0})
		case Task:
			for _, v := range args {
				dl.AddTask(v.(Task))
			}
		case Download:
			dl = args[0].(Download)
		}
	}
	return len(dl), dl
}

func safeArgs(args ...interface{}) (url, title, name, dst string) {
	url = args[0].(string)
	switch len(args) {
	case 1:
		names := strings.Split(url, "/")
		title = names[len(names)-1:][0]
		name = title
		dst, _ = os.Getwd()
	case 2:
		title = args[1].(string)
		name = title
		dst, _ = os.Getwd()
	case 3:
		title, name = args[1].(string), args[2].(string)
		dst, _ = os.Getwd()
	case 4:
		title, name, dst = args[1].(string), args[2].(string), args[3].(string)
	}
	return
}
