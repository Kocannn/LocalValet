# LocalValet 🚀

> A modern desktop application for managing local development services, inspired by Laragon. Built with Wails v2, React, and Go.

![Platform](https://img.shields.io/badge/platform-Linux-blue)
![License](https://img.shields.io/badge/license-MIT-green)

## ✨ Features

- 🎛️ **Service Management** — Start/Stop Apache, Nginx, MySQL, PostgreSQL, Redis, PHP-FPM
- 🔄 **Version Switching** — Switch between PHP, Node.js versions instantly
- 🌐 **Virtual Host Manager** — Auto-generate nginx configs for projects
- 🔒 **SSL Certificates** — Local CA with one-click cert generation
- 📁 **Project Discovery** — Auto-detect Laravel, WordPress, Next.js, etc.
- 📊 **Real-time Monitoring** — Live status updates with adaptive polling
- 📝 **Log Viewer** — Filterable logs with color-coded levels
- 🖥️ **Integrated Terminal** — Context-aware terminal with injected PATH
- 🌙 **Dark Mode** — System-aware theme switching

## 🚀 Quick Start

### Prerequisites

- [Go](https://golang.org/) 1.23+
- [Node.js](https://nodejs.org/) 18+
- [Wails](https://wails.io/) v2

### Installation

```bash
# Clone the repository
git clone https://github.com/Kocannn/LocalValet.git
cd LocalValet

# Install dependencies
go mod download
cd frontend && npm install && cd ..

# Run development server
wails dev

# Or build for production
wails build
```

### First Run

1. Launch LocalValet
2. Services will auto-detect from `config/runtime.json`
3. Use the toggle switches to start/stop services
4. Check the log viewer for status updates

## 📦 Supported Services

| Service | Default Port | Version Management |
|---------|--------------|-------------------|
| Apache | 8080 | ✅ |
| Nginx | 8080 | ✅ |
| MySQL | 3306 | ✅ |
| PostgreSQL | 5432 | ✅ |
| Redis | 6379 | ✅ |
| PHP-FPM | 9074 | ✅ |

## 🏗️ Architecture

```
LocalValet/
├── app.go                    # Main Wails app with bindings
├── main.go                   # Entry point with OnBeforeClose
├── config/
│   ├── runtime.json          # Service versions config
│   ├── vhosts.json           # Virtual host config
│   └── ssl.json              # SSL certificates
├── internal/
│   ├── domain/               # Interfaces and types
│   ├── usecase/              # Business logic
│   ├── infrastructure/       # Monitor, terminal
│   └── platform/             # OS-specific implementations
├── runtime/
│   ├── linux/                # Service binaries
│   ├── pids/                 # PID files
│   ├── logs/                 # Service logs
│   └── certs/                # SSL certificates
└── frontend/
    └── src/
        ├── modules/          # Feature modules
        ├── pages/            # Route pages
        ├── components/       # UI components
        └── services/         # Wails bindings
```

## ⚙️ Configuration

### Service Versions (`config/runtime.json`)

```json
{
  "services": {
    "php-fpm": {
      "activeVersion": "8.4",
      "versions": {
        "8.4": {
          "binary": "runtime/linux/php/8.4/sbin/php-fpm",
          "args": ["--nodaemonize"],
          "workingDir": "runtime/linux/php/8.4"
        }
      }
    }
  }
}
```

### Project Discovery

Projects are auto-discovered from these directories:
- `~/Projects`
- `~/projects`
- `~/Sites`
- `~/sites`
- `~/Code`
- `~/code`
- `~/www`

Frameworks detected: Laravel, WordPress, Next.js, Nuxt, Django, PHP, Node.js

## 🔧 Development

### Build Commands

```bash
# Development with hot reload
wails dev

# Build production binary
wails build

# Regenerate Wails bindings
wails generate module

# Run Go tests
go test ./internal/... -v

# Build frontend only
cd frontend && npm run build
```

### Adding a New Service

1. Add service config to `internal/domain/service/config.go`
2. Add default port to `internal/platform/linux/linux.go`
3. Add runtime config to `config/runtime.json`
4. Place binaries in `runtime/linux/<service>/<version>/`

### Port Conflict Handling

When a default port is in use:
- LocalValet scans +200 ports for an available one
- Config files are patched automatically (nginx, php-fpm)
- CLI flags are injected for other services (mysql, redis)

## 🐛 Troubleshooting

### Services won't start

1. Check if binaries exist in `runtime/linux/`
2. Verify `config/runtime.json` paths
3. Check `runtime/logs/` for error messages
4. Run `GetDiagnostics()` from the app

### Port conflicts

LocalValet auto-remaps ports when conflicts are detected. Check the status message for the actual port being used.

### Permission issues

Some operations require elevated privileges:
- Binding to ports < 1024
- Modifying `/etc/hosts`
- Trusting SSL certificates

### Reset everything

```bash
rm -rf runtime/pids/*
rm -rf runtime/logs/*
rm -rf runtime/certs/*
rm config/vhosts.json config/ssl.json
```

## 📄 License

MIT License - See [LICENSE](LICENSE) for details

## 🙏 Acknowledgments

- [Wails](https://wails.io/) - Go + Web framework
- [shadcn/ui](https://ui.shadcn.com/) - UI components
- [Tailwind CSS](https://tailwindcss.com/) - CSS framework
- [Laragon](https://laragon.org/) - Inspiration

---

Made with ❤️ using [Wails](https://wails.io/)
