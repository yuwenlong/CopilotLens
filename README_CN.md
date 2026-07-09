<div align="center">

# 🚀 CopilotLens

**企业级 AI 用量监控与分析平台**

*实时监控 GitHub Copilot AI Credits 消耗，自动同步数据*

[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go&logoColor=white)](https://golang.org/)
[![Gin](https://img.shields.io/badge/Gin-v1.12-00E676?style=flat&logo=gin&logoColor=white)](https://gin-gonic.com/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**中文版** | [English](README.md)

</div>

---

## 📋 目录

- [✨ 功能特性](#-功能特性)
- [🛠️ 技术栈](#️-技术栈)
- [📁 项目结构](#-项目结构)
- [🚀 快速开始](#-快速开始)
- [📊 数据来源](#-数据来源)
- [🔌 API 文档](#-api-文档)
- [🌐 国际化](#-国际化)
- [🔒 IP 白名单](#-ip-白名单)
- [📜 脚本说明](#-脚本说明)
- [📄 许可证](#-许可证)

---

## ✨ 功能特性

| 功能 | 说明 |
|------|------|
| 📊 **月度总用量** | Hero 渐变卡片 + 数字动画，直观展示本月消耗 |
| 👥 **用户用量** | 柱状图（用户维度）+ 模型明细表（交替背景色、hover 联动） |
| 🤖 **模型用量** | 柱状图 + 饼图，分析各模型消耗占比 |
| 🔌 **GitHub API** | 自动通过 GitHub Billing API 同步数据（无需手动导出 CSV） |
| ⚡ **并发获取** | 并行 API 调用，快速加载数据 |
| 💾 **智能缓存** | 2 小时 TTL 缓存 + 自动清理 + 每小时整点预热 |
| 🎨 **加载动画** | 全屏遮罩 + CSS3 三环旋转动画 |
| 🌐 **国际化** | 中英文一键切换，localStorage 持久化 |
| 🔒 **IP 白名单** | 支持精确 IP 和 CIDR 网段，灵活访问控制 |
| 📜 **进程管理** | PID 文件管理，支持 start/stop/restart/reload |

---

## 🛠️ 技术栈

| 层 | 技术 |
|---|---|
| 🖥️ 后端 | [Go](https://golang.org/) 1.26 + [Gin](https://gin-gonic.com/) v1.12 |
| ⚙️ 配置 | [Viper](https://github.com/spf13/viper) v1.21 + TOML |
| 🎨 前端 | 原生 HTML/CSS/JS + [jQuery](https://jquery.com/) 4.0 + [ECharts](https://echarts.apache.org/) 6.1 |
| 🌐 国际化 | 自研 i18n（`data-i18n` 属性 + JS 翻译文件） |
| 🔌 API | GitHub Billing API（AI Credit 用量） |

---

## 📁 项目结构

```
CopilotLens/
├── main.go                          # 🚀 入口文件 + GitHub 客户端初始化
│
├── domain/dto/                      # 📦 数据传输对象
│   ├── copilot.go                   #    CopilotRecord
│   ├── monthly_total.go             #    MonthlyTotalResponse
│   ├── monthly_user.go              #    UserUsage, UserModel, MonthlyUserResponse
│   └── monthly_model.go             #    ModelUsage, MonthlyModelResponse
│
├── internal/                        # 🔧 内部包
│   ├── client/                      #    数据加载工具
│   │   └── client.go                #    LoadUsers, LoadUsernameMap, Round2
│   ├── config/                      #    配置管理（init() 自动加载）
│   │   └── config.go                #    AppConfig + Config() getter
│   ├── github/                      #    GitHub API 客户端
│   │   ├── client.go                #    Billing API 获取器
│   │   └── cache.go                 #    CacheManager（2h TTL, sync.RWMutex）
│   └── handler/                     #    HTTP 处理器
│       ├── router.go                #    路由注册 + IPWhitelist 中间件
│       ├── total.go                 #    月度总用量 handler
│       ├── user.go                  #    用户用量 + 日用量 handler
│       └── model.go                 #    模型用量 + 日用量 handler
│
├── tasks/                           # ⏱️ 后台任务
│   ├── common.go                    #    ITask 接口 + Init/Stop
│   ├── cache_clean.go               #    缓存过期清理（10分钟）
│   └── cache_warm.go                #    每小时整点缓存预热
│
├── web/                             # 🎨 前端资源
│   ├── index.html                   #    首页（Hero Banner + 功能卡片）
│   ├── monthly-total.html           #    月度总用量（Hero 卡片 + 数字动画）
│   ├── monthly-user.html            #    用户用量（图表 + 明细表）
│   ├── monthly-model.html           #    模型用量（柱状图 + 饼图）
│   └── static/
│       ├── css/style.css            #    全局样式 + 加载动画
│       └── i18n/
│           ├── zh.js                #    中文翻译
│           └── en.js                #    英文翻译
│
├── data/                            # 📊 数据文件
│   └── username.csv                 #    账号→显示名称映射
│
├── toml/config_template.toml        # ⚙️ 配置模板
├── build.sh / build.ps1             # 🔨 构建脚本
├── run.sh / run.ps1                 # ▶️ 进程管理
├── LICENSE                          # 📄 MIT 许可证
└── .gitignore
```

---

## 🚀 快速开始

### 环境要求

- Go 1.26 或更高版本
- GitHub Personal Access Token（需 `copilot` 权限）

### 1️⃣ 克隆 & 构建

```bash
# 克隆仓库
git clone https://github.com/yuwenlong/CopilotLens.git
cd CopilotLens

# 构建（选择你的平台）
# Windows
.\build.ps1

# Linux/macOS
chmod +x build.sh
./build.sh
```

### 2️⃣ 配置

```bash
# 编辑配置文件
vim bin/conf/config.toml
```

```toml
[server]
port = "8080"
whitelist = []  # 空数组 = 允许所有 IP

[github]
token = ""      # GitHub PAT（或设置 GITHUB_TOKEN 环境变量）
org = "your-org"  # GitHub 组织名称
```

### 3️⃣ 运行

```bash
# Windows
.\run.ps1 start

# Linux/macOS
chmod +x run.sh
./run.sh start
```

### 4️⃣ 访问

打开浏览器访问：**http://localhost:8080**

---

## 📊 数据来源

### GitHub Billing API（自动同步）

CopilotLens 通过 GitHub Billing API 自动获取 AI Credits 用量数据：

1. **认证方式**：使用具有 `copilot` 权限的 GitHub Personal Access Token
2. **数据来源**：`GET /organizations/{org}/settings/billing/ai_credit/usage`
3. **自动同步**：缓存未命中时实时获取数据
4. **缓存机制**：结果缓存 2 小时，减少 API 调用
5. **缓存预热**：定时任务在每个整点预拉取当前月/当前日数据，使大部分请求直接命中缓存，不再调用 GitHub API

### 配置方式

通过配置文件或环境变量设置 GitHub Token：

```toml
[github]
token = "ghp_xxxxxxxxxxxx"
org = "your-org"
```

或使用环境变量：
```bash
export GITHUB_TOKEN="ghp_xxxxxxxxxxxx"
```

### 用户名映射

`data/username.csv` 文件用于将 GitHub 账号映射为显示名称：

```csv
账号,姓名
xxxxxx,xx
yyyyyy,yy
```

**用途：**
- 将 GitHub 账号映射为更易读的显示名称
- 支持任意语言（中文、英文等）
- 在图表和表格中提升可读性
- 格式：`账号,显示名称`（第一行为表头）

---

## 🔌 API 文档

### 基础 URL

```
http://localhost:8080
```

### 接口列表

#### 1️⃣ 月度总用量

```http
GET /api/monthly-total?month=YYYY-MM
```

**返回示例：**
```json
{
  "month": "2026-06",
  "total": 185944.03,
  "daily": [
    {"date": "2026-06-30", "amount": 6429.48},
    {"date": "2026-06-29", "amount": 6528.43}
  ]
}
```

#### 2️⃣ 用户用量（月度）

```http
GET /api/monthly-user?month=YYYY-MM
```

**返回示例：**
```json
{
  "month": "2026-06",
  "users": [
    {
      "username": "aaaaa",
      "name": "xx",
      "total": 8542.33,
      "cost": 85.42,
      "models": [
        { "model": "GPT-5.4", "total": 4200.12, "cost": 42.00 },
        { "model": "Claude Sonnet 4.6", "total": 4342.21, "cost": 43.42 }
      ]
    }
  ]
}
```

#### 3️⃣ 用户用量（日度）

```http
GET /api/daily-user?date=YYYY-MM-DD
```

**返回示例：** 结构同月度，但为单日数据。

#### 4️⃣ 模型用量（月度）

```http
GET /api/monthly-model?month=YYYY-MM
```

**返回示例：**
```json
{
  "month": "2026-06",
  "models": [
    { "model": "GPT-5.4", "total": 15234.56, "cost": 152.35 },
    { "model": "Claude Sonnet 4.6", "total": 12456.78, "cost": 124.57 }
  ]
}
```

#### 5️⃣ 模型用量（日度）

```http
GET /api/daily-model?date=YYYY-MM-DD
```

**返回示例：** 结构同月度，但为单日数据。

---

## 🌐 国际化

### 切换语言

点击导航栏右上角的 **"中 | EN"** 按钮即可切换语言。

### 工作原理

- 默认语言：**中文（zh）**
- 切换时设置 `localStorage('lang')` 为 `'en'` 或 `'zh'`
- 页面重新加载，加载对应的 i18n 文件
- 所有 `data-i18n` 属性自动翻译

### 添加翻译

编辑 `web/static/i18n/zh.js` 和 `web/static/i18n/en.js`：

```javascript
window.i18n = {
    title: {
        index: 'Copilot AI Credits 用量分析系统',
        // ... 更多 key
    }
};
```

---

## 🔒 IP 白名单

### 配置方式

编辑 `conf/config.toml`：

```toml
[server]
port = "8080"

# 空数组 = 允许所有 IP
whitelist = []

# 或指定允许的 IP/网段
whitelist = [
    "127.0.0.1",
    "::1",
    "192.168.1.0/24",
    "10.0.0.0/8"
]
```

### 支持格式

| 格式 | 示例 | 说明 |
|------|------|------|
| 精确 IP | `192.168.1.100` | 单个 IP 地址 |
| IPv6 | `::1` | 本地回环地址 |
| CIDR | `192.168.1.0/24` | 子网范围（256 个地址） |
| CIDR | `10.0.0.0/8` | 大型子网（1600 万地址） |

### 行为

- ✅ **空白名单**：允许所有 IP 访问
- ❌ **非空白名单**：仅允许列表中的 IP/网段访问
- 🚫 **被拒绝**：返回 `403 Forbidden` + `{"error": "ip not allowed"}`

---

## 📜 脚本说明

### 构建脚本

| 脚本 | 平台 | 用途 |
|------|------|------|
| `build.ps1` | Windows | 编译 Go 二进制 + 复制资源到 `bin/` |
| `build.sh` | Linux/macOS | 同上 |

### 运行脚本

| 命令 | 说明 |
|------|------|
| `.\run.ps1 start` | 后台启动服务 |
| `.\run.ps1 stop` | 停止服务 |
| `.\run.ps1 restart` | 重启服务（不重新构建） |
| `.\run.ps1 reload` | 重新构建 + 重启 |

### PID 管理

- 服务 PID 存储在 `bin/.pid`
- `start` 检查服务是否已运行
- `stop` 优雅终止进程
- `reload` 触发完整构建后重启

---

## 📄 许可证

本项目基于 MIT 许可证开源 - 详见 [LICENSE](LICENSE) 文件了解详情。

---

<div align="center">

**使用 Go + Gin + ECharts + GitHub API 构建 ❤️**

</div>
