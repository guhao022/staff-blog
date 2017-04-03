package curl

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"
)

var (
	wg         sync.WaitGroup
	curLine    int = -1
	maxNameLen int
	mutex      *sync.RWMutex = new(sync.RWMutex)
	count      int           = 0
)

type (
	Task struct {
		Url   string
		Title string
		Name  string
		Dst   string
		Code  int
	}

	Download []Task

	// 逐行获取内容和行数
	processFunc func(content string, line int) bool
)

// 接收参数并返回新的任务结构
func (this Task) New(args ...interface{}) Task {
	if len(args) == 0 {
		panic(CurlError{"curl.New()", -6, "curl.New() parameter type error."})
	} else {
		this.Url, this.Title, this.Name, this.Dst = safeArgs(args...)
	}
	return this
}

// 将task填加到downlaod结构体中
func (this *Download) AddTask(ts Task) {
	*this = append(*this, ts)
}

// 根据key值返回download的值
func (this Download) GetValues(key string) []string {
	var arr []string
	for i := 0; i < len(this); i++ {
		v := reflect.ValueOf(this[i]).FieldByName(key)
		arr = append(arr, v.String())
	}
	return arr
}

// 下载Get方法
func Get(url string) (code int, res *http.Response, err error) {

	// get res
	res, err = http.Get(url)

	if err != nil {
		return -5, res, CurlError{url, -5, err.Error()}
	}

	// check state code
	if res.StatusCode != 200 {
		s := fmt.Sprintf("%v an [%v] error occurred.", url, res.StatusCode)
		return -1, res, CurlError{url, -1, s}
	}

	return 0, res, err
}

func ReadLine(body io.ReadCloser, process processFunc) error {

	var content string
	var err error
	var line int = 1

	buff := bufio.NewReader(body)

	for {
		content, err = buff.ReadString('\n')

		if line > 1 && (err != nil || err == io.EOF) {
			break
		}

		if ok := process(content, line); ok {
			break
		}

		line++
	}

	return err
}

func New(args ...interface{}) (dl Download, errStack []CurlError) {
	curLine = -1
	count, dl = parseArgs(args...)
	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(CurlError); ok {
				errStack = append(errStack, v)
			} else {
				errStack = append(errStack, CurlError{"curl.New()", -5, err})
			}
		}
	}()

	maxNameLen = maxTitleLength(dl.GetValues("Title"))
	header(&dl)

	wg.Add(count)
	for i := 0; i < count; i++ {
		progressbar(dl[i].Title, time.Now(), 0, "\n")
		go func(dl Download, num int) {
			download(&dl[num], num, count, &errStack)
			wg.Done()
		}(dl, i)
	}
	wg.Wait()

	curDown(count - curLine)
	footer()

	return
}

func download(ts *Task, line, max int, errStack *[]CurlError) {
	url, title, name, dst := ts.Url, ts.Title, ts.Name, safeDst(ts.Dst)
	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(CurlError); ok {
				*errStack = append(*errStack, v)
				ts.Code = v.code
			} else {
				*errStack = append(*errStack, CurlError{url, -5, err})
				ts.Code = -5
			}
			curMove(line, max)
			msg := fmt.Sprintf("%v download error.", safeTitle(title))
			empty := strings.Repeat(" ", 80-len(msg))
			fmt.Printf("\r%v%v", msg, empty)
		}
	}()

	// get url
	code, res, err := Get(url)
	if code == -1 {
		panic(err)
	}
	defer res.Body.Close()

	// create dst
	if !isDirExist(dst) {
		if err := os.Mkdir(dst, 0777); err != nil {
			panic(CurlError{url, -2, "Create folder error, Error: " + err.Error()})
		}
	}

	// create file
	file, createErr := os.Create(dst + name)
	if createErr != nil {
		panic(CurlError{url, -2, "Create file error, Error: " + createErr.Error()})
	}
	defer file.Close()

	// verify content length
	if res.ContentLength == -1 && isBodyBytes(res.Header.Get("Content-Type")) {
		panic(CurlError{url, -4, "Download content length is -1."})
	}

	start := time.Now()
	if isBodyBytes(res.Header.Get("Content-Type")) {
		buf := make([]byte, res.ContentLength)
		var m float32
		for {
			n, err := res.Body.Read(buf)
			if n == 0 && err.Error() == "EOF" {
				break
			}
			if err != nil && err.Error() != "EOF" {
				panic(CurlError{url, -7, "Download size error, Error: ." + err.Error()})
			}
			m = m + float32(n)
			i := int(m / float32(res.ContentLength) * 50)
			file.WriteString(string(buf[:n]))

			func(title string, start time.Time, i, line, max int) {
				curMove(line, max)
				progressbar(title, start, i, "")
			}(title, start, i, line, max)
		}

		// valid download exe
		fi, err := file.Stat()
		if err == nil {
			if fi.Size() != res.ContentLength {
				panic(CurlError{url, -3, "Downlaod size verify error, please check your network."})
			}
		}
	} else {
		if bytes, err := ioutil.ReadAll(bufio.NewReader(res.Body)); err != nil {
			panic(CurlError{url, -8, err.Error()})
		} else {
			file.Write(bytes)
			curMove(line, max)
			progressbar(title, start, 50, "")
		}
	}
}

func curMove(line, max int) {
	mutex.Lock()
	switch {
	case curLine == -1:
		curReset(max - line)
	case line < curLine:
		curUp(curLine - line)
	case line > curLine:
		curDown(line - curLine)
	}
	if curLine != line {
		curLine = line
	}
	mutex.Unlock()
}

func isDirExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		return true
	}
}

func isBodyBytes(content string) (isBytes bool) {
	if strings.Index(content, "json") != -1 {
		isBytes = false
	} else if strings.Index(content, "text") != -1 {
		isBytes = false
	} else if strings.Index(content, "application") != -1 {
		isBytes = true
	}
	return
}

func maxTitleLength(titles []string) int {
	max := 0
	for _, v := range titles {
		if len(v) > max {
			max = len(v)
		}
	}
	if max > 15 {
		max = 15
	}
	return max
}

func safeTitle(title string) string {
	h := ""
	if len(title) > 15 {
		title = title[:12] + "..."
	} else if len(title) <= maxNameLen {
		h = strings.Repeat(" ", maxNameLen-len(title))
	}
	return h + title + ":"
}

func safeDst(dst string) string {
	if !strings.HasSuffix(dst, "/") {
		dst += "/"
	}
	return dst
}

func curReset(i int) {
	fmt.Printf("\r\033[%dA", i)
}

func curUp(i int) {
	fmt.Printf("\r\033[%dA", i)
}

func curDown(i int) {
	fmt.Printf("\r\033[%dB", i)
}
