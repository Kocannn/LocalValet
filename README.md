# LocalValet

<div align="center">

![LocalValet Banner](https://img.shields.io/badge/LocalValet-v2.0.0--MVP-6366f1?style=for-the-badge&logo=electron&logoColor=white)

**The Lightweight, Modern Local Development Environment Orchestrator for Linux**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![React](https://img.shields.io/badge/React-18-61DAFB?style=flat-square&logo=react)](https://reactjs.org)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.0+-3178C6?style=flat-square&logo=typescript)](https://typescriptlang.org)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind-v4-38B2AC?style=flat-square&logo=tailwind-css)](https://tailwindcss.com)
[![Wails v2](https://img.shields.io/badge/Wails-v2-DF0000?style=flat-square&logo=wails)](https://wails.io)
[![License: MIT](https://img.shields.io/badge/License-MIT-emerald.svg?style=flat-square)](LICENSE)
[![Status: MVP Complete](https://img.shields.io/badge/MVP-100%25%20Complete-10b981.svg?style=flat-square)](docs/mvp-roadmap.md)

</div>

---

## 🌟 Overview

**LocalValet** is an all-in-one local development orchestrator inspired by Laragon and Laravel Valet, designed natively for Linux. Built with a clean architecture combining **Go** on the backend and **React + TypeScript + Tailwind CSS** on the frontend (powered by **Wails v2**).

LocalValet gives developers an instant, native alternative to heavy Docker setups for daily web development, offering zero-delay process startup, multi-version runtime hot-switching, automatic framework discovery, virtual hosts generation, and a pure-Go local SSL Certificate Authority.

---

## ✨ Key Features

### 🚀 1. Native Service Lifecycle Manager
- **Supported Modules**: Apache, Nginx, MariaDB / MySQL, PostgreSQL, Redis, and PHP-FPM.
- **Port Conflict Auto-Remapping**: Automatically detects when default ports are occupied and remaps across a `+200` range with dynamic CLI argument injection (`--port`, `-p`, `-D PORT`).
- **Health Check Engine**: Real-time process PID polling, TCP socket connection probes, and HTTP response validation.
- **Graceful Shutdown**: Strict process group isolation (`syscall.Setpgid`) with 5-second `SIGTERM` polling and `SIGKILL` timeout fallback.
- **Dependency Auto-Start**: Starting web servers (`apache` or `nginx`) automatically boots `php-fpm`.

### ⚡ 2. Multi-Version Runtime Switching & Hot-Restart
- **PHP Multi-Version**: Switch between PHP 8.4, 8.3, 8.2, and 8.1 seamlessly.
- **Node.js Multi-Version**: Switch between Node.js 22, 20, and 18.
- **Seamless Hot-Restart**: If a service is currently running when its version is changed, LocalValet gracefully restarts it using the new binary without crashing.
- **Dynamic Discovery**: Drop new version folders into `runtime/linux/<service>/<version>` and LocalValet will auto-detect them dynamically.

### 🔍 3. Automatic Project Discovery & Virtual Hosts
- **Framework Auto-Detection**:
  - **Laravel**: Detected via `artisan` / `composer.json` (document root: `public/`).
  - **WordPress**: Detected via `wp-config.php` / `wp-load.php`.
  - **Next.js / Nuxt**: Detected via `next.config.*` / `nuxt.config.*` (reverse-proxy to `:3000`).
  - **React / Vue / Vite**: Detected via `vite.config.*` (document root: `dist/`).
  - **Generic PHP / Static HTML**: Automatically served.
- **Nginx Virtual Hosts**: Automatically generates `.conf` virtual host templates in `runtime/linux/nginx/vhosts/<domain>.conf` with FastCGI proxying, SPA fallback routing (`try_files $uri $uri/ /index.php?$query_string`), and WebSocket-enabled reverse proxying.
- **Local `.test` Domains**: Standardized local domain mapping (e.g. `my-project.test`).

### 🔒 4. Pure-Go Local SSL Certificate Authority
- **Zero CLI Dependencies**: Built exclusively with Go standard library `crypto/x509`, `crypto/rsa`, and `encoding/pem` (no `mkcert` or `openssl` CLI needed).
- **Local Root CA**: Generates a persistent Root CA (`runtime/certs/ca.crt` & `ca.key`) with 10-year validity.
- **Per-Project Certificates**: Auto-issues signed leaf certificates with Subject Alternative Names (`DNS:project.test`, `DNS:*.project.test`, `IP:127.0.0.1`).
- **System Trust Helper**: 1-click installer to trust the Root CA in `/usr/local/share/ca-certificates/` so browsers display a trusted padlock (`https://`).

### 🖥️ 5. Context Terminal & IDE Launcher
- **Injected PATH Environment**: Launches your favorite Linux terminal (`kitty`, `alacritty`, `gnome-terminal`, `konsole`, `xfce4-terminal`, etc.) with active PHP, Node.js, Composer, and MySQL binaries placed at the front of `PATH`.
- **Direct IDE Opening**: Open discovered projects in VS Code (`code`), Cursor, PhpStorm, or system `$EDITOR` in one click.

---

## 🏛️ Architecture & Clean Design

```
┌─────────────────────────────────────────────────────────────┐
│                    Frontend (React 18 / TS / Vite)          │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌────────┐ │
│  │ Service │ │ Project │ │  VHost  │ │ Version │ │  Logs  │ │
│  │ Control │ │ Browser │ │ Manager │ │ Switch  │ │ Viewer │ │
│  └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘ └───┬────┘ │
│       └───────────┴───────────┼───────────┴──────────┘      │
│                               │ Wails v2 IPC Bridge         │
├───────────────────────────────┼─────────────────────────────┤
│                               │ Go Backend                  │
│  ┌────────────────────────────┴──────────────────────────┐  │
│  │                     Use Case Layer                    │  │
│  │  ServiceUC  │  ProjectUC  │  VHostUC  │  SSLUC        │  │
│  └────────────────────────────┬──────────────────────────┘  │
│                               │                             │
│  ┌────────────────────────────┴──────────────────────────┐  │
│  │                    Domain & Platform Layer            │  │
│  │  LinuxManager │ Scanner │ NginxGen │ CAManager │ DNS │  │
│  └────────────────────────────┬──────────────────────────┘  │
│                               │                             │
│  ┌────────────────────────────┴──────────────────────────┐  │
│  │              Runtime Filesystem (Isolated)            │  │
│  │  runtime/linux/*  │  runtime/certs/*  │  runtime/logs │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

---

## 📦 Getting Started

### Prerequisites
- **Linux** (Ubuntu, Debian, Fedora, Arch, Manjaro, Pop!_OS, etc.)
- **Go 1.21+**
- **Node.js 18+** & `npm`
- **Wails v2 CLI** (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

### Quick Setup

```bash
# 1. Clone the repository
git clone https://github.com/Kocannn/LocalValet.git
cd LocalValet

# 2. Setup runtime environment structure
./scripts/setup-runtime.sh

# 3. Install frontend dependencies
cd frontend && npm install && cd ..

# 4. Run in development mode
wails dev
```

### Production Build

```bash
# Build standalone optimized binary & frontend bundle
./scripts/build-linux.sh
```
The compiled executable is created at `build/bin/localvalet`.

---

## 🌐 Local Domain & SSL Setup

### 1. Synchronizing `/etc/hosts`
To route `*.test` domains to `127.0.0.1`, click **Sync Domains to /etc/hosts** in the **Settings** page or let LocalValet sync automatically. It updates only the managed block:
```hosts
# BEGIN LocalValet Managed Domains
127.0.0.1    my-laravel-app.test
127.0.0.1    my-wp-blog.test
127.0.0.1    nextjs-app.test
# END LocalValet Managed Domains
```

### 2. Trusting the Root CA
In **Settings**, click **Trust Root CA in System Store**. This copies `runtime/certs/ca.crt` to `/usr/local/share/ca-certificates/` and runs `update-ca-certificates`.

- **Chrome / Chromium / Brave / Edge**: Automatically respects system CA certificates.
- **Firefox**: Ensure `"Security & Privacy" -> "Certificates" -> "View Certificates" -> "Authorities"` has imported `runtime/certs/ca.crt` or set `security.enterprise_roots.enabled = true` in `about:config`.

---

## 📁 Project & Runtime Directory Layout

```
LocalValet/
├── config/
│   ├── runtime.json        # Service versions & binary definitions
│   └── projects.json       # Configured scan roots & project metadata
├── runtime/
│   ├── certs/              # LocalValet Root CA & per-domain SSL certs
│   ├── logs/               # Service log outputs
│   ├── pids/               # Running service process IDs
│   └── linux/              # Portable service runtime binaries
│       ├── php/            # Multi-version PHP (8.4, 8.3, 8.2)
│       ├── node/           # Multi-version Node.js (22, 20, 18)
│       ├── mysql/          # MariaDB / MySQL
│       ├── nginx/          # Nginx binary & vhosts/
│       ├── apache/         # Apache httpd
│       ├── redis/          # Redis server
│       └── postgresql/     # PostgreSQL
├── scripts/
│   ├── build-linux.sh      # Production build automation script
│   └── setup-runtime.sh    # Runtime folder initialization script
└── build/bin/localvalet    # Standalone application binary
```

---

## 🛠️ Testing

LocalValet includes an extensive unit and integration test suite covering domain interfaces, platform implementations, framework scanners, cryptographic certificate generators, and use case orchestration:

```bash
# Run all Go unit and integration tests
go test -v ./internal/...

# Run frontend type-check & Vite build verification
cd frontend && npm run build
```

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
