# LocalValet v2 - System Architecture

## Architecture Overview

LocalValet uses a **layered clean architecture** with Wails v2 as the IPC bridge between Go backend and React frontend.

```
┌─────────────────────────────────────────────────────────────┐
│                    Frontend (React/TS)                       │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐          │
│  │ Service  │ │  VHost  │ │  SSL    │ │ Project │          │
│  │ Manager  │ │ Manager │ │ Manager │ │ Finder  │          │
│  └────┬─────┘ └────┬────┘ └────┬────┘ └────┬────┘          │
│       └────────────┼──────────┼────────────┘               │
│                    │ Wails Bindings │                        │
├────────────────────┼───────────────┼────────────────────────┤
│                    │ Go Backend    │                        │
│  ┌─────────────────┴───────────────┴─────────────────┐     │
│  │              Use Case Layer                        │     │
│  │  ServiceUC │ VHostUC │ SSLUC │ ProjectUC │ TermUC  │     │
│  └─────────────────────┬───────────────────────────┘     │
│                        │                                   │
│  ┌─────────────────────┴───────────────────────────┐     │
│  │              Domain Layer                         │     │
│  │  Interfaces │ Types │ Ports │ Events              │     │
│  └─────────────────────┬───────────────────────────┘     │
│                        │                                   │
│  ┌─────────────────────┴───────────────────────────┐     │
│  │              Infrastructure Layer                 │     │
│  │  Monitor │ Terminal │ Config │ Events │ Plugin     │     │
│  └─────────────────────┬───────────────────────────┘     │
│                        │                                   │
│  ┌─────────────────────┴───────────────────────────┐     │
│  │              Platform Layer (OS-specific)          │     │
│  │  Linux │ Windows │ macOS │ Docker                   │     │
│  └───────────────────────────────────────────────────┘     │
└─────────────────────────────────────────────────────────────┘
```

## Layer Responsibilities

### Domain Layer (`internal/domain/`)
Pure interfaces and types. No dependencies on other layers.

```go
// service/manager.go
type Manager interface {
    StartService(name string) error
    StopService(name string) error
    GetServiceStatus(name string) (bool, string)
    SetServiceVersion(name, version string) error
    GetActiveServiceVersion(name string) (string, error)
    GetAvailableVersions(name string) ([]string, error)
}

// vhost/manager.go
type VHostManager interface {
    DiscoverProjects(root string) ([]Project, error)
    GenerateVHost(project Project) error
    RemoveVHost(domain string) error
    ListVHosts() ([]VHost, error)
}

// ssl/manager.go
type SSLManager interface {
    GenerateCert(domain string) (CertPair, error)
    TrustCert(certPath string) error
    IsTrusted(domain string) bool
}

// terminal/manager.go
type TerminalManager interface {
    Launch(opts LaunchOptions) (LaunchResult, error)
}
```

### Use Case Layer (`internal/usecase/`)
Business logic orchestration. Depends only on domain interfaces.

```go
// service/usecase.go
type ServiceUseCase struct {
    manager domain.ServiceManager
    configs []domain.ServiceConfig
}

func (uc *ServiceUseCase) ToggleService(name string, start bool) LogMessage {
    if start {
        return uc.StartService(name)
    }
    return uc.StopService(name)
}
```

### Infrastructure Layer (`internal/infrastructure/`)
External integrations and cross-cutting concerns.

- **Monitor**: Adaptive polling (300ms fast / 5s slow), debounced events
- **Terminal**: Shell detection, init script generation, env injection
- **Config**: JSON file read/write with atomic operations
- **Events**: Wails event emitter wrapper
- **Plugin**: Plugin loader and registry

### Platform Layer (`internal/platform/`)
OS-specific implementations behind factory pattern.

```go
// factory.go
func NewServiceManager() domain.ServiceManager {
    switch runtime.GOOS {
    case "linux":
        return linux.NewLinuxManager()
    case "darwin":
        return darwin.NewDarwinManager()
    case "windows":
        return windows.NewWindowsManager()
    }
}
```

## IPC Communication Layer

### Wails v2 Bindings

Go methods on the `App` struct are auto-generated to TypeScript:

```go
// app.go - Go side
type App struct { ... }

func (a *App) GetServiceStatus(name string) domain.Status { ... }
func (a *App) ToggleService(name string, start bool) LogMessage { ... }
```

```typescript
// wailsjs/go/main/App.d.ts - Auto-generated
export function GetServiceStatus(name: string): Promise<Status>;
export function ToggleService(name: string, start: boolean): Promise<LogMessage>;
```

### Event System (Real-time Updates)

```
Go Backend                          Frontend
    │                                   │
    ├─ EventsEmit("service:status") ──►├─ EventsOn("service:status")
    ├─ EventsEmit("service:log") ─────►├─ EventsOn("service:log")
    ├─ EventsEmit("vhost:discovered")─►├─ EventsOn("vhost:discovered")
    └─ EventsEmit("project:found") ──►└─ EventsOn("project:found")
```

### Error Handling Pattern

```go
// Go side - return structured errors
type ServiceError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Service string `json:"service"`
}

// Frontend - handle by error code
try {
    await ToggleService("mysql", true);
} catch (err) {
    if (err.code === "PORT_CONFLICT") {
        showPortConflictDialog(err);
    }
}
```

## Plugin System

### Plugin Interface

```go
// plugin/plugin.go
type Plugin interface {
    Name() string
    Version() string
    Init(app AppContext) error
    Services() []ServiceDefinition
    Destroy() error
}

type AppContext interface {
    RegisterService(def ServiceDefinition)
    EmitEvent(event string, data interface{})
    GetConfig() *Config
}
```

### Plugin Lifecycle

```
Load Plugin → Init(ctx) → Register Services → [Running] → Destroy()
```

### Service Definition

```go
type ServiceDefinition struct {
    Name         string
    DisplayName  string
    Binary       string
    DefaultPort  int
    HealthCheck  func() bool
    ConfigFiles  []string
}
```

## Storage Schema

### config/services.json
```json
{
  "services": {
    "apache": {
      "displayName": "Apache",
      "defaultPort": 8080,
      "binary": "runtime/linux/apache/2.4/bin/httpd",
      "configFiles": ["conf/httpd.conf"],
      "healthCheck": "http://localhost:8080"
    }
  }
}
```

### config/runtime.json
```json
{
  "services": {
    "php-fpm": {
      "activeVersion": "8.4",
      "versions": {
        "8.4": {
          "binary": "runtime/linux/php/8.4/sbin/php-fpm",
          "args": ["--nodaemonize", "--fpm-config", "runtime/linux/php/8.4/etc/php-fpm.conf"],
          "workingDir": "runtime/linux/php/8.4"
        }
      }
    }
  }
}
```

### config/vhosts.json
```json
{
  "vhosts": {
    "myproject.test": {
      "root": "/home/user/projects/myproject/public",
      "framework": "laravel",
      "ssl": true,
      "port": 443
    }
  }
}
```

### config/projects.json
```json
{
  "projects": [
    {
      "name": "myproject",
      "path": "/home/user/projects/myproject",
      "framework": "laravel",
      "phpVersion": "8.4",
      "hasDocker": true
    }
  ]
}
```

## Module Boundaries

```
┌─────────────────────────────────────────────────────────────┐
│                     App (app.go)                             │
├─────────────────────────────────────────────────────────────┤
│ ServiceManager │ VHostManager │ SSLManager │ ProjectFinder  │
├────────────────┼──────────────┼────────────┼────────────────┤
│ VersionSwitcher│ Terminal     │ Docker     │ LogViewer      │
└────────────────┴──────────────┴────────────┴────────────────┘
```

Each module:
- Has its own domain interfaces
- Has its own use case layer
- Can be tested independently
- Can be disabled via feature flags
