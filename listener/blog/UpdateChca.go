package blog

import (
	"net/url"
	"os"
	"staff/tools/curl"
	"strings"
	"os/exec"
	"runtime"
	"path"
	"staff/tools/targz"

	"github.com/num5/axiom"
)

// 更新chca
func (b *BlogListener) updateChca(ctx *axiom.Context, m string) {

	// 下载chca
	downUrl := b.ChcaUrl
	ctx.Reply("下载CHCA，下载链接【" + downUrl + "】...")

	ctx.Reply("文件下载中，请稍后...")

	tarFile, err := b.download(downUrl)
	if err != nil {
		ctx.Reply("下载失败，错误信息：" + err.Error())
		return
	}

	ctx.Reply("下载完成，开始解压缩...")

	// 解压chca
	err = targz.Extract(tarFile, b.WorkerDir)
	if err != nil {
		ctx.Reply("解压失败，错误信息：" + err.Error())
		return
	}
	ctx.Reply("解压缩完成，复制文件...\n")

	if runtime.GOOS != "windows" {
		cmd :=exec.Command("chmod", "777", 	b.WorkerDir+"/chca")
		if err := cmd.Start(); err != nil {
			ctx.Reply("修改文件权限失败,请手动修改\n")
			return
		}
	}

	os.Remove(tarFile)

	theme := path.Join(b.WorkerDir, "theme", b.Theme)
	if !Exist(theme) {
		ctx.Reply("检测到博客模板不存在，下载默认模板...")

		tarTheme, err := b.download(b.ThemeUrl)
		if err != nil {
			ctx.Reply("下载失败，错误信息：" + err.Error())
			return
		}

		ctx.Reply("模板下载完成，开始解压缩...")

		// 解压
		err = targz.Extract(tarTheme, b.WorkerDir + "/theme/")
		if err != nil {
			ctx.Reply("解压模板文件失败，错误信息：" + err.Error())
			return
		}
		ctx.Reply("解压模板文件完成，复制文件...\n")
		os.Remove(tarTheme)
	}

	ctx.Reply(m + "成功")

}

// 下载
func (b *BlogListener) download(downUrl string) (string, error) {
	fileUrl, err := url.Parse(downUrl)
	if err != nil {
		return "", err
	}

	filePath := fileUrl.Path

	fbs := strings.Split(filePath, "/")

	fileName := fbs[len(fbs)-1]

	_, cerr := curl.New(downUrl)

	if cerr != nil {
		return "", cerr[0]
	}

	return fileName, nil
}
