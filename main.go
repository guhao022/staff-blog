package main

import (
	"staff/listener/blog"
	"staff/tools/env"

	"github.com/num5/axiom"
)

func main() {
	env, err := env.Load()
	if err != nil {
		panic(err)
	}

	blogListener := &blog.BlogListener{
		WorkerDir: env.Get("CHCA_WORKER_DIR"),
		MarkdownDir: env.Get("BLOG_MARKDOWN_DIR"),
		Host: env.Get("BLOG_HOST"),
		UploadTpl: env.Get("UPLOAD_TEMPLATE"),
		HtmlDir: env.Get("BLOG_HTML_DIR"),
		ChcaUrl: env.Get("CHCA_DOWNLOAD_URL"),
	}

	b := axiom.New(env.Get("BOT_NAME"))
	b.AddAdapter(axiom.NewShell(b))
	b.Register(blogListener)

	b.Start()
}
