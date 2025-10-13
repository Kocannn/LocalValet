# LocalValet Setup Instructions

## Regenerate Wails Bindings

Setelah menambahkan fungsi baru di `app.go`, Anda perlu regenerate bindings Wails:

```bash
# Development mode (akan auto-generate bindings)
wails dev

# Atau manual generate
wails generate module
```

## Service Names Configuration

### Linux (Ubuntu/Debian)

- Apache: `apache2`
- MySQL: `mysql`

### macOS (dengan Homebrew)

- Apache: `httpd`
- MySQL: `mysql`

### Permissions

Untuk Linux, Anda mungkin perlu menambahkan user ke sudoers atau memberikan permission:

```bash
# Tambahkan user ke sudo group
sudo usermod -aG sudo $USER

# Atau edit sudoers untuk allow tanpa password (DEVELOPMENT ONLY!)
sudo visudo
# Tambahkan: your_username ALL=(ALL) NOPASSWD: /bin/systemctl
```

## Testing

1. Build aplikasi:

```bash
wails build
```

2. Jalankan aplikasi:

```bash
./build/bin/LocalValet-dev-linux-amd64
```

## Features Implemented

### Backend (Go)

- ✅ `GetServiceStatus(serviceName string)` - Check if service is running
- ✅ `StartService(serviceName string)` - Start a service
- ✅ `StopService(serviceName string)` - Stop a service
- ✅ `ToggleService(serviceName, shouldStart)` - Toggle service on/off
- ✅ Cross-platform support (Linux, macOS, Windows)

### Frontend (React + TypeScript)

- ✅ Service management UI with switches
- ✅ Real-time status checking (every 5 seconds)
- ✅ Loading states during operations
- ✅ Integrated logging system
- ✅ Color-coded log levels (info, success, warning, error)
- ✅ Auto-scroll logs
- ✅ Clear logs button

## Service Management

The app supports managing these services:

- **Apache** - Web server
- **MySQL** - Database server

You can easily add more services by updating the `modules` state in `App.tsx`.

## Troubleshooting

### Bindings not found

Run `wails dev` or `wails generate module` to regenerate TypeScript bindings.

### Permission denied

Make sure your user has sudo permissions for systemctl commands.

### Service not starting

Check logs in the UI for specific error messages.
Check service status manually:

```bash
# Linux
sudo systemctl status apache2
sudo systemctl status mysql

# macOS
brew services list
```
