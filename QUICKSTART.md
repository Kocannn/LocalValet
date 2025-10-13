# LocalValet - Quick Start Guide

## ✅ Yang Sudah Dibuat

### 1. Backend (Go) - 3 Files

**📄 app.go**

- `GetServiceStatus(serviceName)` - Cek status service
- `StartService(serviceName)` - Start service
- `StopService(serviceName)` - Stop service
- `ToggleService(serviceName, shouldStart)` - Toggle on/off
- Support: Linux (systemctl), macOS (brew), Windows (net)

**📄 service_utils.go**

- `CheckSudoAccess()` - Cek sudo permission
- `GetSystemInfo()` - Info sistem OS
- Helper functions untuk command execution

**📄 main.go**

- Sudah ada, tidak perlu diubah

### 2. Frontend (React + TypeScript)

**📄 App.tsx** - Fully Integrated!

- ✅ Service management table (Apache & MySQL)
- ✅ Real-time status monitoring (every 5s)
- ✅ Toggle switches untuk start/stop
- ✅ Loading states dengan spinner
- ✅ Integrated logging system
- ✅ Auto-scroll logs
- ✅ Clear logs button
- ✅ Color-coded log levels

**📄 app-layout.tsx**

- ✅ Sidebar layout component
- ✅ Header dengan title
- ✅ Responsive design

**📄 app-sidebar.tsx**

- ✅ Navigation sidebar
- ✅ Version switcher
- ✅ Search form

### 3. Documentation

**📄 SETUP.md** - Setup instructions
**📄 IMPLEMENTATION_GUIDE.md** - Complete implementation guide
**📄 dev.sh** - Development helper script

## 🚀 Cara Menjalankan

### Step 1: Generate Wails Bindings

```bash
cd /home/kocan/Coding/LocalValet

# Cara tercepat (recommended)
wails dev

# Atau manual
wails generate module
```

### Step 2: Setup Permissions (Linux Only)

```bash
# Edit sudoers untuk allow systemctl tanpa password
sudo visudo

# Tambahkan (ganti 'kocan' dengan username Anda):
kocan ALL=(ALL) NOPASSWD: /bin/systemctl start *
kocan ALL=(ALL) NOPASSWD: /bin/systemctl stop *
kocan ALL=(ALL) NOPASSWD: /bin/systemctl status *
kocan ALL=(ALL) NOPASSWD: /bin/systemctl is-active *
```

### Step 3: Run Development

```bash
# Menggunakan script helper
chmod +x dev.sh
./dev.sh

# Atau langsung
wails dev
```

## 🎯 Cara Menggunakan Aplikasi

1. **Lihat Status Service**

   - Status akan auto-update setiap 5 detik
   - Badge hijau = Active, abu-abu = Inactive

2. **Start/Stop Service**

   - Toggle switch untuk on/off
   - Lihat loading spinner saat proses
   - Cek logs untuk hasil operasi

3. **Monitor Logs**
   - Semua operasi tercatat di logs
   - Color-coded: Info (biru), Success (hijau), Warning (kuning), Error (merah)
   - Auto-scroll ke log terbaru
   - Click trash icon untuk clear logs

## 🔧 Service Names per Platform

**Linux (Ubuntu/Debian):**

```typescript
serviceName: "apache2"; // untuk Apache
serviceName: "mysql"; // untuk MySQL
```

**macOS (Homebrew):**

```typescript
serviceName: "httpd"; // untuk Apache
serviceName: "mysql"; // untuk MySQL
```

**Windows:**

```typescript
serviceName: "Apache2.4"; // untuk Apache
serviceName: "MySQL80"; // untuk MySQL
```

## 📁 File Structure

```
LocalValet/
├── app.go                          # Backend logic
├── service_utils.go                # Helper functions
├── main.go                         # Entry point
├── dev.sh                          # Development script
├── SETUP.md                        # Setup guide
├── IMPLEMENTATION_GUIDE.md         # Full guide
├── frontend/
│   └── src/
│       ├── App.tsx                 # Main app with service management
│       └── components/
│           ├── app-layout.tsx      # Layout component
│           ├── app-sidebar.tsx     # Sidebar navigation
│           └── ui/                 # shadcn components
└── wailsjs/
    └── go/
        └── main/
            └── App.js              # Generated bindings (after wails dev)
```

## 🎨 Features Overview

### Service Management

- [x] Start/Stop Apache
- [x] Start/Stop MySQL
- [x] Real-time status monitoring
- [x] Loading states
- [x] Error handling

### Logging System

- [x] Timestamp setiap log
- [x] Log levels (info, success, warning, error)
- [x] Color-coded display
- [x] Auto-scroll
- [x] Clear logs

### UI/UX

- [x] Sidebar navigation
- [x] Responsive layout
- [x] Loading spinners
- [x] Status badges
- [x] Clean design dengan shadcn/ui

## 🐛 Common Issues

### 1. Bindings Not Found

```bash
# Solution:
wails generate module
# Atau restart dev server
```

### 2. Permission Denied (Linux)

```bash
# Test sudo access:
sudo -n systemctl status apache2

# Setup sudoers (lihat Step 2)
```

### 3. Service Not Starting

- Cek service name sesuai OS
- Cek logs di UI untuk error detail
- Test manual: `sudo systemctl start apache2`

### 4. Frontend Errors

- Cek browser console
- Pastikan wails dev running
- Restart dev server

## 📝 Next Development Ideas

1. **Add More Services**

   - PostgreSQL
   - Redis
   - Nginx
   - PHP-FPM

2. **Enhanced Features**

   - Restart service button
   - Service configuration viewer
   - Port management
   - Service dependencies
   - Auto-start on boot

3. **Advanced Logging**

   - Export logs to file
   - Filter logs by level
   - Search in logs
   - Tail service logs (journalctl)

4. **System Info**
   - Resource usage (CPU, Memory)
   - Disk space
   - Network stats
   - Service uptime

## 🎓 Learning Resources

- **Wails**: https://wails.io/docs/introduction
- **shadcn/ui**: https://ui.shadcn.com
- **React**: https://react.dev
- **TypeScript**: https://www.typescriptlang.org
- **systemctl**: https://www.freedesktop.org/software/systemd/man/systemctl.html

## 💡 Pro Tips

1. **Development**: Gunakan `wails dev` untuk hot-reload
2. **Debugging**: Check browser console dan logs UI
3. **Testing**: Test service commands di terminal dulu
4. **Security**: Untuk production, gunakan proper service management
5. **Performance**: Status check 5s sudah optimal

## 📞 Need Help?

1. Baca `IMPLEMENTATION_GUIDE.md` untuk detail teknis
2. Baca `SETUP.md` untuk setup instructions
3. Check Wails documentation
4. Test commands manually di terminal

---

**Status**: ✅ Ready to Use!
**Platform**: Linux (primary), macOS, Windows
**Framework**: Wails v2 + React + TypeScript + shadcn/ui
