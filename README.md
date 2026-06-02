<div align="center">

# 🚀 CopilotLens

**Enterprise AI Usage Monitoring & Analytics Platform**

*Analyze GitHub Copilot AI Credits consumption with beautiful visualizations*

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

---

## 📁 Project Structure

```
CopilotLens/
├── main.go                          # 🚀 Entry point (24 lines)
│
├── domain/dto/                      # 📦 Data Transfer Objects
│   ├── copilot.go                   #    CopilotRecord
│   ├── monthly_total.go             #    MonthlyTotalResponse
│   ├── monthly_user.go              #    UserUsage, UserModel, MonthlyUserResponse
│   └── monthly_model.go             #    ModelUsage, MonthlyModelResponse
│
├── internal/                        # 🔧 Internal packages
│   ├── client/                      #    Data loading
│   │   └── client.go                #    LoadUsernameMap, LoadCopilotCSV, Round2
│   ├── conf/                        #    Configuration
│   │   └── config.go                #    Global config + IP whitelist
│   └── handler/                     #    HTTP handlers
│       ├── router.go                #    Route registration + IPWhitelist middleware
│       ├── total.go                 #    MonthlyTotal handler
│       ├── user.go                  #    MonthlyUser handler
│       └── model.go                 #    MonthlyModel handler
│
├── web/                             # 🎨 Frontend
│   ├── index.html                   #    Homepage with hero banner
│   ├── monthly-total.html           #    Monthly total with hero card
│   ├── monthly-user.html            #    User usage charts + table
│   ├── monthly-model.html           #    Model usage charts
│   └── static/
│       ├── css/style.css            #    Global styles (423 lines)
│       └── i18n/
│           ├── zh.js                #    Chinese translations
│           └── en.js                #    English translations
│
├── data/                            # 📊 Data files
│   ├── copilot/                     #    Copilot usage CSVs (from GitHub)
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
- Git (optional)

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
# whitelist = ["127.0.0.1", "::1", "192.168.1.0/24"]
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

### Copilot Usage Data

Export from GitHub:

1. Go to **GitHub Organization** → **Settings** → **Billing and licensing** → **AI usage**
2. Or visit directly: `/organizations/{org}/settings/billing/ai_usage`
3. Click **"Get usage report"** to export current or previous month's data
4. Save as `YYYY-MM.csv` (e.g., `2026-06.csv`)
5. Place in `data/copilot/` directory

### Username Mapping

The `data/username.csv` file maps GitHub usernames to display names:

```csv
账号,姓名
xxxxxx,xx
yyyyyyy,yy
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
  "total": 36754.18
}
```

#### 2️⃣ User Usage

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

#### 3️⃣ Model Usage

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
| `./run.sh start` | Start server in background |
| `./run.sh stop` | Stop running server |
| `./run.sh restart` | Restart server |
| `./run.sh reload` | Rebuild + restart |

**Windows:**
```powershell
.\run.ps1 start
.\run.ps1 stop
.\run.ps1 restart
.\run.ps1 reload
```

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

**Built with ❤️ using Go + Gin + ECharts**

</div>
