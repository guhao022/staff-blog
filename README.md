# Staff 静态博客部署助手

### TODO

- [x] 下载/更新 chca
- [x] 使用chca自动部署博客
- [ ] Docker部署项目
- [ ] websocket适配器

### staff 说明
staff是用来帮助用户部署静态博客的小助手，基于 [Axiom](https://github.com/num5/axiom) 运维机器人框架开发

### 安装

```go
go get github.com/num5/staff
```

### 使用

安装完成以后首先需要创建配置文件，在staff工作目录创建.env文件，文件配置如下：
```env
# chca下载目录，tar.gz压缩
# linux 平台
CHCA_DOWNLOAD_URL=https://github.com/num5/chca/releases/download/v1.1/chca-linux_amd64.tar.gz
# windows 平台
#CHCA_DOWNLOAD_URL=https://github.com/num5/chca/releases/download/v1.1/chca-windows_amd64.tar.gz
# mac 平台
#CHCA_DOWNLOAD_URL=https://github.com/num5/chca/releases/download/v1.1/chca-darwin_amd64.tar.gz

# chca工作目录
CHCA_WORKER_DIR=/chca/

# markdown存放目录(基于chca工作目录)
BLOG_MARKDOWN_DIR=markdown

# blog 编译后的html存放目录
BLOG_HTML_DIR=/data/www/blog/

# 机器人名称
BOT_NAME=staff

# 博客地址(不要加 http)
BLOG_HOST = www.golune.com
```

运行(内置shell适配器)
```bash
staff
Axiom> 下载chca
......
Axiom> 编译博客
......
Axiom> 打开web服务器
......
```

命令
1. 更新chca|更新博客生成器|下载chca|下载博客生成器
> 这条命令用于下载或者更新chca博客生成器

2. 编译博客|博客编译|更新博客|博客更新|编译markdown|编译MARKDOWN|markdown编译|MARKDOWN编译
> 这条命令用于将markdown文件编译成博客html文件，编译好的可以用nginx或者apache部署

3. 开启博客|开启webserver|开启服务器|打开博客服务器|打开web|打开web服务器|打开服务器
> 开启博客内部服务器，默认端口9900，可以自定义端口，例如：
```bash
// 以8080端口运行内置web服务器
Axiom> 开启webserver 端口:8080
```



