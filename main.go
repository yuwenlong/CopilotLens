package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"copilotlens/internal/conf"
	"copilotlens/internal/handler"
)

func main() {
	r := gin.Default()
	r.Use(handler.IPWhitelist())
	r.LoadHTMLGlob("web/*.html")
	r.Static("/static", "./web/static")

	h := &handler.Handler{DataDir: "data"}
	h.RegisterRoutes(r)

	addr := fmt.Sprintf(":%s", conf.Cfg.Server.Port)
	log.Printf("服务启动，监听端口 %s", conf.Cfg.Server.Port)
	r.Run(addr)
}
