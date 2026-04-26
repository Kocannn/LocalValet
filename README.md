# LocalValet 🚀

> A modern desktop application for managing local development services (Apache, MySQL, etc.) built with Wails, React, TypeScript, and shadcn/ui.

![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20macOS%20%7C%20Windows-blue)
![License](https://img.shields.io/badge/license-MIT-green)

## ✨ Features

- 🎛️ **Service Management** - Start/Stop local services with toggle switches
- 📊 **Real-time Monitoring** - Auto-refresh service status every 5 seconds
- 📝 **Integrated Logging** - Color-coded logs with timestamps and log levels
- 🎨 **Modern UI** - Built with shadcn/ui components and Tailwind CSS
- 🧱 **Isolated Runtime** - Services run from app-managed `runtime/` binaries and configs
- ⚡ **Fast & Lightweight** - Native performance with Go backend and React frontend

## 🚀 Quick Start

### Prerequisites

- [Go](https://golang.org/) (1.18+)
- [Node.js](https://nodejs.org/) (16+)
- [Wails](https://wails.io/) v2

### Installation & Run

```bash
# Install dependencies
go mod download
cd frontend && npm install && cd ..

# Run development server
wails dev

# Or use the helper script
chmod +x dev.sh test-services.sh
./dev.sh
```

### Runtime Isolation (Linux)

LocalValet now uses app-isolated runtimes defined by `config/runtime.json`.

```bash
config/runtime.json
runtime/linux/php/8.2/
runtime/linux/php/8.3/
runtime/pids/
runtime/logs/
```

To switch PHP version, update `activeVersion` for `php-fpm` in `config/runtime.json`
or use Settings page in the app.

### Test Services

```bash
# Run service test script
./test-services.sh
```

## 📚 Documentation

- **[QUICKSTART.md](QUICKSTART.md)** - Get started quickly
- **[SETUP.md](SETUP.md)** - Detailed setup instructions
- **[IMPLEMENTATION_GUIDE.md](IMPLEMENTATION_GUIDE.md)** - Complete implementation guide

## 🎯 Supported Services

| Service    | Canonical Key |
| ---------- | ------------- |
| Apache     | apache        |
| MySQL      | mysql         |
| PostgreSQL | postgresql    |
| Redis      | redis         |
| Nginx      | nginx         |
| PHP-FPM    | php-fpm       |

## 🏗️ Architecture

### Backend (Go)

```
app.go              # Main application logic
├── GetServiceStatus()   # Check service status
├── StartService()       # Start a service
├── StopService()        # Stop a service
└── ToggleService()      # Toggle service on/off

service_utils.go    # Helper functions
service_config.go   # Service configurations
```

### Frontend (React + TypeScript)

```
src/
├── App.tsx                    # Main app with service management
├── components/
│   ├── app-layout.tsx         # Layout wrapper with sidebar
│   ├── app-sidebar.tsx        # Navigation sidebar
│   └── ui/                    # shadcn/ui components
```

## 🛠️ Development

### Adding New Services

Update `App.tsx`:

```typescript
const [modules, setModules] = useState<ServiceModule[]>([
  {
    name: "Apache",
    serviceName: "apache2",
    isRunning: false,
    isLoading: false,
  },
  { name: "MySQL", serviceName: "mysql", isRunning: false, isLoading: false },
  // Add new service here
]);
```

### Build for Production

```bash
wails build
```

## 🐛 Troubleshooting

### Wails Bindings Not Found

```bash
wails generate module
# Or restart: wails dev
```

### Service Not Starting

1. Ensure binary paths in `config/runtime.json` exist
2. Check logs in `runtime/logs/`
3. Verify active version and runtime config for the service

## 📧 Contact

- GitHub: [@Kocannn](https://github.com/Kocannn)
- Repository: [LocalValet](https://github.com/Kocannn/LocalValet)

## 🙏 Acknowledgments

- [Wails](https://wails.io/) - Go + Web framework
- [shadcn/ui](https://ui.shadcn.com/) - UI components
- [Tailwind CSS](https://tailwindcss.com/) - CSS framework

---

Made with ❤️ using [Wails](https://wails.io/)
