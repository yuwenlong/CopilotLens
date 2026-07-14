package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"copilotlens/internal/client"
	"copilotlens/internal/config"
	"copilotlens/internal/github"
	"copilotlens/internal/handler"
	"copilotlens/tasks"
)

type dailyWriter struct {
	dir     string
	curDate string
	file    *os.File
}

func (w *dailyWriter) Write(p []byte) (n int, err error) {
	today := time.Now().Format("2006-01-02")
	if today != w.curDate {
		if w.file != nil {
			w.file.Close()
		}
		_ = os.MkdirAll(w.dir, 0755)
		f, err := os.OpenFile(filepath.Join(w.dir, "copilotlens-"+today+".log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return 0, err
		}
		w.file = f
		w.curDate = today
	}
	return w.file.Write(p)
}

var conf = config.Config()

func main() {
	dw := &dailyWriter{dir: "logs"}
	log.SetOutput(io.MultiWriter(os.Stderr, dw))
	gin.DefaultWriter = io.MultiWriter(os.Stdout, dw)
	token := conf.GitHub.Token
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	if token == "" || conf.GitHub.Org == "" {
		log.Fatal("GitHub API 未配置，请在 conf/config.toml 的 [github] 段设置 token 和 org")
	}
	client.GitHubClient = github.NewClient(token, conf.GitHub.Org)
	log.Printf("GitHub API 已启用，组织: %s", conf.GitHub.Org)

	// 缓存清理任务和周期性预热任务先启动（Run() 内部首次立即预热，不阻塞 HTTP）
	tasks.Init(client.GitHubClient.Cache, client.GitHubClient)

	r := gin.Default()
	r.Use(handler.IPWhitelist())
	r.LoadHTMLGlob("web/*.html")
	r.Static("/static", "./web/static")

	h := handler.NewHandler("data", client.GitHubClient)
	h.RegisterRoutes(r)

	addr := fmt.Sprintf(":%s", conf.Server.Port)
	log.Printf("服务启动，监听端口 %s", conf.Server.Port)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP 服务启动失败: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("收到终止信号，正在关闭服务...")

	tasks.Stop()
	log.Println("定时任务已停止")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("HTTP 服务关闭异常: %v", err)
	}
	log.Println("HTTP 服务已停止")
}
