package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"copilotlens/internal/client"
	"copilotlens/internal/config"
	"copilotlens/internal/github"
	"copilotlens/internal/handler"
	"copilotlens/tasks"
)

var conf = config.Config()

func main() {
	token := conf.GitHub.Token
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	if token == "" || conf.GitHub.Org == "" {
		log.Fatal("GitHub API 未配置，请在 conf/config.toml 的 [github] 段设置 token 和 org")
	}
	client.GitHubClient = github.NewClient(token, conf.GitHub.Org)
	log.Printf("GitHub API 已启用，组织: %s", conf.GitHub.Org)
	tasks.Init(client.GitHubClient.Cache)

	r := gin.Default()
	r.Use(handler.IPWhitelist())
	r.LoadHTMLGlob("web/*.html")
	r.Static("/static", "./web/static")

	h := &handler.Handler{DataDir: "data"}
	h.RegisterRoutes(r)

	addr := fmt.Sprintf(":%s", conf.Server.Port)
	log.Printf("服务启动，监听端口 %s", conf.Server.Port)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		tasks.Stop()
		log.Println("定时任务已停止")
		os.Exit(0)
	}()

	r.Run(addr)
}
