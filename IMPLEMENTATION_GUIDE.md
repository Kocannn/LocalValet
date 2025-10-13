# LocalValet - Service Management Implementation Guide

## 📋 Overview

LocalValet sekarang memiliki sistem lengkap untuk mengelola service Apache dan MySQL dengan fitur:

- ✅ Start/Stop services dengan toggle switch
- ✅ Real-time status monitoring
- ✅ Integrated logging system
- ✅ Cross-platform support (Linux, macOS, Windows)

## 🏗️ Arsitektur

### Backend (Go)

**File: `app.go`**

- `GetServiceStatus(serviceName string) ServiceStatus` - Cek status service
- `StartService(serviceName string) LogMessage` - Start service
- `StopService(serviceName string) LogMessage` - Stop service
- `ToggleService(serviceName, shouldStart bool) LogMessage` - Toggle service

**File: `service_utils.go`**

- `CheckSudoAccess() bool` - Cek sudo access
- `GetSystemInfo() map[string]string` - Informasi sistem
- `getServiceCommand(action, serviceName) *exec.Cmd` - Helper untuk command

### Frontend (React + TypeScript)

**File: `frontend/src/App.tsx`**

**Interfaces:**

```typescript
interface LogEntry {
  timestamp: string;
  level: "info" | "warning" | "error" | "success";
  message: string;
}

interface ServiceModule {
  name: string; // Display name
  serviceName: string; // System service name
  isRunning: boolean; // Current status
  isLoading: boolean; // Loading state
}
```

**Features:**

- Real-time service status (check every 5 seconds)
- Loading states during operations
- Auto-scroll logs
- Color-coded log levels
- Clear logs functionality

## 🚀 Cara Menggunakan

### 1. Generate Wails Bindings

Setelah membuat/mengubah fungsi Go, generate bindings:

```bash
# Cara 1: Menggunakan script helper
./dev.sh

# Cara 2: Manual
wails generate module
wails dev
```

### 2. Service Names per Platform

Update `modules` state di `App.tsx` sesuai platform:

**Linux:**

```typescript
const [modules, setModules] = useState<ServiceModule[]>([
  {
    name: "Apache",
    serviceName: "apache2",
    isRunning: false,
    isLoading: false,
  },
  { name: "MySQL", serviceName: "mysql", isRunning: false, isLoading: false },
]);
```

**macOS:**

```typescript
const [modules, setModules] = useState<ServiceModule[]>([
  { name: "Apache", serviceName: "httpd", isRunning: false, isLoading: false },
  { name: "MySQL", serviceName: "mysql", isRunning: false, isLoading: false },
]);
```

### 3. Permissions Setup (Linux)

Untuk menjalankan systemctl tanpa password prompt:

```bash
# Edit sudoers
sudo visudo

# Tambahkan baris ini (ganti username dengan user Anda):
username ALL=(ALL) NOPASSWD: /bin/systemctl start *
username ALL=(ALL) NOPASSWD: /bin/systemctl stop *
username ALL=(ALL) NOPASSWD: /bin/systemctl restart *
username ALL=(ALL) NOPASSWD: /bin/systemctl status *
username ALL=(ALL) NOPASSWD: /bin/systemctl is-active *
```

⚠️ **CATATAN:** Ini untuk development. Untuk production, gunakan polkit atau service management yang lebih secure.

### 4. Build & Run

```bash
# Development mode
wails dev

# Build for production
wails build

# Run build
./build/bin/LocalValet-dev-linux-amd64
```

## 🎨 UI Components

### Service Management Table

- **Module Column**: Display name (Apache, MySQL)
- **Status Column**: Badge dengan warna (Active=hijau, Inactive=abu)
- **Action Column**: Toggle switch untuk on/off

### Server Logs Panel

- **Background**: Light (#FAFAFA) untuk visibility
- **Font**: Monospace untuk console feel
- **Format**: `[HH:MM:SS] [LEVEL] Message`
- **Colors**:
  - 🔵 INFO: Blue
  - 🟢 SUCCESS: Green
  - 🟡 WARNING: Yellow
  - 🔴 ERROR: Red
- **Features**: Auto-scroll, Clear button

## 🔧 Cara Menambah Service Baru

1. **Update Frontend (`App.tsx`):**

```typescript
const [modules, setModules] = useState<ServiceModule[]>([
  {
    name: "Apache",
    serviceName: "apache2",
    isRunning: false,
    isLoading: false,
  },
  { name: "MySQL", serviceName: "mysql", isRunning: false, isLoading: false },
  // Tambahkan service baru
  {
    name: "PostgreSQL",
    serviceName: "postgresql",
    isRunning: false,
    isLoading: false,
  },
  { name: "Redis", serviceName: "redis", isRunning: false, isLoading: false },
]);
```

2. **Backend sudah generic**, tidak perlu diubah!

## 🐛 Troubleshooting

### Error: "Module has no exported member"

```bash
# Solution: Regenerate Wails bindings
wails generate module
# Atau restart dev server
wails dev
```

### Error: "Permission denied"

```bash
# Check sudo access
sudo -n true

# Setup sudoers (lihat section Permissions Setup)
```

### Service tidak start/stop

1. Cek service name sesuai dengan sistem OS
2. Cek permission user
3. Lihat logs di UI untuk detail error
4. Test manual di terminal:

```bash
# Linux
sudo systemctl status apache2
sudo systemctl start apache2

# macOS
brew services list
brew services start httpd
```

### Logs tidak muncul

- Check browser console untuk errors
- Pastikan fungsi `addLog()` dipanggil
- Pastikan state `logs` terupdate

## 📝 Flow Diagram

```
User Toggle Switch
       ↓
handleServiceToggle()
       ↓
Set Loading State → Add Info Log
       ↓
Call WailsApp.ToggleService()
       ↓
Backend: StartService/StopService
       ↓
Execute systemctl/brew command
       ↓
Return LogMessage
       ↓
Frontend: Add Result Log
       ↓
Update Module Status → Remove Loading
       ↓
Auto Status Check (5s interval)
```

## 🎯 Next Steps

1. **Auto-detect platform** dan set service names otomatis
2. **Add restart functionality** untuk services
3. **Service logs streaming** dari systemctl journal
4. **Configuration management** untuk service settings
5. **Health checks** untuk services
6. **Notification system** untuk errors
7. **Service dependencies** handling

## 📚 Resources

- [Wails Documentation](https://wails.io/docs/introduction)
- [shadcn/ui Components](https://ui.shadcn.com)
- [systemctl Documentation](https://www.freedesktop.org/software/systemd/man/systemctl.html)
- [Homebrew Services](https://docs.brew.sh/Manpage#services-subcommand)

## 💡 Tips

1. **Testing**: Test di terminal dulu sebelum implement di app
2. **Logging**: Selalu log error untuk debugging
3. **State Management**: Gunakan loading states untuk UX yang baik
4. **Error Handling**: Handle semua possible errors
5. **Cross-platform**: Test di semua platform yang didukung
