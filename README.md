<div align="center">

# 🚀 CopilotLens

**Enterprise AI Usage Monitoring & Analytics Platform**

*Real-time GitHub Copilot AI Credits monitoring with automatic data synchronization*

[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go&logoColor=white)](https://golang.org/)
[![Gin](https://img.shields.io/badge/Gin-v1.12-00E676?style=flat&logo=gin&logoColor=white)](https://gin-gonic.com/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

[中文版](README_CN.md) | **English**

</div>

---

## 📋 Table of Contents

- [✨ Features](#-features)
- [🛠️ Tech Stack](#️-tech-stack)
- [📁 Project Structure](#-project-structure)
- [🚀 Quick Start](#-quick-start)
- [📊 Data Source](#-data-source)
- [🔌 API Documentation](#-api-documentation)
- [🌐 Internationalization](#-internationalization)
- [🔒 IP Whitelist](#-ip-whitelist)
- [📜 Scripts](#-scripts)
- [📄 License](#-license)

---

## ✨ Features

| Feature | Description |
|---------|-------------|
| 📊 **Monthly Total** | Hero card with animated number showing total AI Credits consumed |
| 👥 **User Usage** | Bar chart by user + model breakdown table with hover effects |
| 🤖 **Model Usage** | Bar chart + pie chart for model distribution analysis |
| 🔌 **GitHub API** | Automatic data sync via GitHub Billing API (no manual CSV export) |
| ⚡ **Concurrent Fetching** | Parallel API calls for fast data loading |
| 💾 **Smart Caching** | 2-hour TTL cache with automatic cleanup and hourly warm-up |
| 🎨 **Loading Animation** | Full-screen overlay with CSS3 ring animation |
| 🌐 **i18n Support** | Chinese/English toggle with localStorage persistence |
| 🔒 **IP Whitelist** | Access control with exact IP and CIDR subnet support |
| 📜 **Process Management** | Start/stop/restart/reload with PID file tracking |

---

## 🛠️ Tech Stack

| Layer | Technology |
|-------|------------|
| 🖥️ Backend | [Go](https://golang.org/) 1.26 + [Gin](https://gin-gonic.com/) v1.12 |
| ⚙️ Config | [Viper](https://github.com/spf13/viper) v1.21 + TOML |
| 🎨 Frontend | Vanilla HTML/CSS/JS + [jQuery](https://jquery.com/) 4.0 + [ECharts](https://echarts.apache.org/) 6.1 |
| 🌐 i18n | Custom implementation with `data-i18n` attributes |
| 🔌 API | GitHub Billing API (AI Credit usage) |

---

## 📁 Project Structure

```
CopilotLens/
├── main.go                          # 🚀 Entry point + GitHub client init
│
├── domain/dto/                      # 📦 Data Transfer Objects
│   ├── copilot.go                   #    CopilotRecord
│   ├── monthly_total.go             #    MonthlyTotalResponse
│   ├── monthly_user.go              #    UserUsage, UserModel, MonthlyUserResponse
│   └── monthly_model.go             #    ModelUsage, MonthlyModelResponse
│
├── internal/                        # 🔧 Internal packages
│   ├── client/                      #    Data loading utilities
│   │   └── client.go                #    LoadUsers, LoadUsernameMap, Round2
│   ├── config/                      #    Configuration (init() auto-load)
│   │   └── config.go                #    AppConfig + Config() getter
│   ├── github/                      #    GitHub API client
│   │   ├── client.go                #    Billing API fetchers
│   │   └── cache.go                 #    CacheManager (2h TTL, sync.RWMutex)
│   └── handler/                     #    HTTP handlers
│       ├── router.go                #    Route registration + IPWhitelist middleware
│       ├── total.go                 #    MonthlyTotal handler
│       ├── user.go                  #    MonthlyUser + DailyUser handlers
│       └── model.go                 #    MonthlyModel + DailyModel handlers
│
├── tasks/                           # ⏱️ Background tasks
│   ├── common.go                    #    ITask interface + Init/Stop
│   ├── cache_clean.go               #    Cache expiration cleanup (10min)
│   └── cache_warm.go                #    Hourly cache warm-up on the hour
│
├── web/                             # 🎨 Frontend
│   ├── index.html                   #    Homepage with hero banner
│   ├── monthly-total.html           #    Monthly total with hero card
│   ├── monthly-user.html            #    User usage charts + table
│   ├── monthly-model.html           #    Model usage charts
│   └── static/
│       ├── css/style.css            #    Global styles + loading animation
│       └── i18n/
│           ├── zh.js                #    Chinese translations
│           └── en.js                #    English translations
│
├── data/                            # 📊 Data files
│   └── username.csv                 #    Username → display name mapping
│
├── toml/config_template.toml        # ⚙️ Config template
├── build.sh / build.ps1             # 🔨 Build scripts
├── run.sh / run.ps1                 # ▶️ Process management
├── LICENSE                          # 📄 MIT License
└── .gitignore
```

---

## 🚀 Quick Start

### Prerequisites

- Go 1.26 or higher
- GitHub Personal Access Token with `copilot` scope

### 1️⃣ Clone & Build

```bash
# Clone the repository
git clone https://github.com/yuwenlong/CopilotLens.git
cd CopilotLens

# Build (choose your platform)
# Windows
.\build.ps1

# Linux/macOS
chmod +x build.sh
./build.sh
```

### 2️⃣ Configure

```bash
# Edit config file
vim bin/conf/config.toml
```

```toml
[server]
port = "8080"
whitelist = []  # Empty = allow all IPs

[github]
token = ""      # GitHub PAT (or set GITHUB_TOKEN env var)
org = "your-org"  # GitHub organization name
```

### 3️⃣ Run

```bash
# Windows
.\run.ps1 start

# Linux/macOS
chmod +x run.sh
./run.sh start
```

### 4️⃣ Access

Open your browser and visit: **http://localhost:8080**

---

## 📊 Data Source

### GitHub Billing API (Automatic)

CopilotLens automatically fetches AI Credits usage data from the GitHub Billing API:

1. **Authentication**: Uses GitHub Personal Access Token with `copilot` scope
2. **Data Source**: `GET /organizations/{org}/settings/billing/ai_credit/usage`
3. **Auto-sync**: Data is fetched in real-time on each request (fallback when cache misses)
4. **Caching**: Results cached for 2 hours to reduce API calls
5. **Warm-up**: A scheduled task pre-fetches the current month/day data every hour on the hour, so most requests hit the cache without calling GitHub API

### Configuration

Set your GitHub token via config file or environment variable:

```toml
[github]
token = "ghp_xxxxxxxxxxxx"
org = "your-org"
```

Or use environment variable:
```bash
export GITHUB_TOKEN="ghp_xxxxxxxxxxxx"
```

### Username Mapping

The `data/username.csv` file maps GitHub usernames to display names:

```csv
账号,姓名
xxxxxx,xx
yyyyyy,yy
```

**Purpose:**
- Maps usernames to human-readable display names
- Supports any language (Chinese, English, etc.)
- Used in charts and tables for better readability
- Format: `username,display_name` (first row is header)

---

## 🔌 API Documentation

### Base URL

```
http://localhost:8080
```

### Endpoints

#### 1️⃣ Monthly Total

```http
GET /api/monthly-total?month=YYYY-MM
```

**Response:**
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

#### 2️⃣ User Usage (Monthly)

```http
GET /api/monthly-user?month=YYYY-MM
```

**Response:**
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

#### 3️⃣ User Usage (Daily)

```http
GET /api/daily-user?date=YYYY-MM-DD
```

**Response:** Same structure as monthly, but for a single day.

#### 4️⃣ Model Usage (Monthly)

```http
GET /api/monthly-model?month=YYYY-MM
```

**Response:**
```json
{
  "month": "2026-06",
  "models": [
    { "model": "GPT-5.4", "total": 15234.56, "cost": 152.35 },
    { "model": "Claude Sonnet 4.6", "total": 12456.78, "cost": 124.57 }
  ]
}
```

#### 5️⃣ Model Usage (Daily)

```http
GET /api/daily-model?date=YYYY-MM-DD
```

**Response:** Same structure as monthly, but for a single day.

---

## 🌐 Internationalization

### Toggle Language

Click the **"中 | EN"** button in the navigation bar to switch languages.

### How It Works

- Default language: **Chinese (zh)**
- Toggle sets `localStorage('lang')` to `'en'` or `'zh'`
- Page reloads and loads the appropriate i18n file
- All `data-i18n` attributes are dynamically translated

### Adding Translations

Edit `web/static/i18n/zh.js` and `web/static/i18n/en.js`:

```javascript
window.i18n = {
    title: {
        index: 'Copilot AI Credits 用量分析系统',
        // ... more keys
    }
};
```

---

## 🔒 IP Whitelist

### Configuration

Edit `conf/config.toml`:

```toml
[server]
port = "8080"

# Empty array = allow all IPs
whitelist = []

# Or specify allowed IPs/subnets
whitelist = [
    "127.0.0.1",
    "::1",
    "192.168.1.0/24",
    "10.0.0.0/8"
]
```

### Supported Formats

| Format | Example | Description |
|--------|---------|-------------|
| Exact IP | `192.168.1.100` | Single IP address |
| IPv6 | `::1` | Loopback address |
| CIDR | `192.168.1.0/24` | Subnet range (256 addresses) |
| CIDR | `10.0.0.0/8` | Large subnet (16M addresses) |

### Behavior

- ✅ **Empty whitelist**: Allow all IPs
- ❌ **Non-empty whitelist**: Only listed IPs/subnets can access
- 🚫 **Blocked**: Returns `403 Forbidden` with `{"error": "ip not allowed"}`

---

## 📜 Scripts

### Build Scripts

| Script | Platform | Purpose |
|--------|----------|---------|
| `build.ps1` | Windows | Compile Go binary + copy assets to `bin/` |
| `build.sh` | Linux/macOS | Same as above |

### Run Scripts

| Command | Description |
|---------|-------------|
| `./run.ps1 start` | Start server in background |
| `./run.ps1 stop` | Stop running server |
| `./run.ps1 restart` | Restart server (no rebuild) |
| `./run.ps1 reload` | Rebuild + restart |

### PID Management

- Server PID is stored in `bin/.pid`
- `start` checks if server is already running
- `stop` gracefully terminates the process
- `reload` triggers full rebuild before restart

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<div align="center">

**Built with ❤️ using Go + Gin + ECharts + GitHub API**

</div>
