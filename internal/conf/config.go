package conf

import (
	"log"
	"net"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	Port      string
	Whitelist []string
}

var Cfg *Config

func init() {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath("conf")

	if err := v.ReadInConfig(); err != nil {
		log.Fatal("读取配置失败:", err)
	}

	Cfg = new(Config)
	if err := v.Unmarshal(Cfg); err != nil {
		log.Fatal("解析配置失败:", err)
	}
}

func (s *ServerConfig) IsAccessAllowed(clientIP string) bool {
	if len(s.Whitelist) == 0 {
		return true
	}
	ip := net.ParseIP(clientIP)
	if ip == nil {
		return false
	}
	for _, entry := range s.Whitelist {
		if _, cidr, err := net.ParseCIDR(entry); err == nil {
			if cidr.Contains(ip) {
				return true
			}
			continue
		}
		if net.ParseIP(entry) != nil && ip.Equal(net.ParseIP(entry)) {
			return true
		}
	}
	return false
}
