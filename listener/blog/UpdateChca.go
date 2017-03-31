package blog

import (
	"github.com/num5/axiom"
	"net/url"
	"os"
	"staff/tools/archive"
	"staff/tools/curl"
	"strings"
)

// 更新chca
func (b *BlogListener) updateChca(ctx *axiom.Context, m string) {

	// 下载chca
	downUrl := os.Getenv("CHCA_DOWNLOAD_URL")
	ctx.Reply("下载CHCA，下载链接【" + downUrl + "】...")

	ctx.Reply("文件下载中，请稍后...")

	tarFile, err := b.download(downUrl)
	if err != nil {
		ctx.Reply("下载失败，错误信息：" + err.Error())
		return
	}

	ctx.Reply("下载完成，开始解压缩...")

	// 解压chca
	err = archive.UnTarGz(tarFile, WORKER_DIR)
	if err != nil {
		ctx.Reply("解压失败，错误信息：" + err.Error())
		return
	}
	ctx.Reply("解压缩完成，复制文件...")

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
