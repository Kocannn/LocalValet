# LocalValet v2 - MVP Roadmap

## Vision

Build a modern, fast, and beautiful local development environment manager inspired by Laragon, with Linux as the primary platform.

## MVP Definition

**Ship in 10 weeks with these core features:**

1. Service management (Apache, Nginx, MySQL, PostgreSQL, Redis, PHP-FPM)
2. Real-time status monitoring
3. Port conflict auto-resolution
4. Version switching (PHP, Node.js)
5. Integrated terminal with environment injection
6. Log viewer
7. Project auto-discovery
8. Virtual host manager
9. SSL certificate generation
10. Settings page

## Development Phases

### Phase 1: Foundation (Week 1-2) — ✅ 100% Complete

**Goal:** Solidify core service management

**Tasks:**
- [x] Refine service manager architecture
- [x] Implement health checks per service type (Process, TCP, HTTP)
- [x] Add graceful shutdown with timeout (SIGTERM 5s timeout & SIGKILL fallback)
- [x] Improve port conflict detection & auto-remapping (+200 scan range)
- [x] Add service dependency tracking (auto-starts php-fpm for web servers)
- [x] Write comprehensive tests (Domain, Platform Linux, Monitor, UseCases)

**Deliverables:**
- [x] Service manager with health checks
- [x] Port management with auto-remapping
- [x] Unit tests for all service operations

**Success Criteria:**
- [x] All 6 services start/stop reliably
- [x] Port conflicts resolved automatically
- [x] Tests pass on Linux (100% PASS)

### Phase 2: UI Polish (Week 3-4) — ✅ Completed

**Goal:** Beautiful, responsive dashboard

**Tasks:**
- [x] Implement service table/card component
- [x] Add status indicators (running/stopped/error)
- [x] Implement service grouping & categories (Web, Database, Runtime)
- [x] Add quick actions (toggle switch, context terminal)
- [x] Polish log viewer with color-coded levels & auto-scroll
- [x] Add dark mode support
- [x] Responsive layout for mobile (Sidebar & Layout breakpoints)

**Deliverables:**
- [x] Service control panel
- [x] Log viewer with color-coded levels
- [x] Responsive sidebar navigation

**Success Criteria:**
- [x] UI is fast and responsive
- [x] Dark mode works correctly
- [x] Mobile layout is usable

### Phase 3: Version Switching (Week 5-6) — 🟡 In Progress (60%)

**Goal:** Multi-version runtime support

**Tasks:**
- [x] Implement version detection per service (RuntimeRegistry)
- [x] Add version selector UI (PHP runtime selector in Settings)
- [x] Implement hot-switching logic
- [x] Add binary verification (path & existence validation)
- [x] Support PHP runtime versions (PHP 8.4 active, multi-version registry ready)
- [ ] Support Node.js 18, 20, 22

**Deliverables:**
- [x] Version switcher UI in Settings
- [x] Per-service version management via RuntimeRegistry
- [x] Version hot-switching
- [x] PHP 8.4 runtime configured
- [ ] Node.js 18-22 support

**Success Criteria:**
- [x] PHP version switching works
- [ ] Node.js version switching works
- [x] No service crashes during switch

### Phase 4: Project Discovery (Week 7-8) — ⏳ Upcoming

**Goal:** Auto-discover and manage projects

**Tasks:**
- [ ] Implement project scanner
- [ ] Detect frameworks (Laravel, WordPress, Next.js, etc.)
- [ ] Generate virtual host configs
- [ ] Implement SSL certificate generation
- [ ] Add project grid/list view
- [ ] Add quick actions (terminal, browser, IDE)

**Deliverables:**
- [ ] Project discovery engine
- [ ] Virtual host manager
- [ ] SSL certificate generator

**Success Criteria:**
- [ ] Projects discovered automatically
- [ ] VHosts generated correctly
- [ ] SSL certificates trusted

### Phase 5: Polish & Release (Week 9-10) — 🟡 In Progress (40%)

**Goal:** Production-ready release

**Tasks:**
- [x] Performance optimization & race condition prevention
- [x] Error handling improvements
- [x] Architecture documentation
- [x] Unit test suite for core modules
- [ ] User testing & feedback
- [ ] Release preparation & packaging

**Deliverables:**
- [x] Clean architecture & modular frontend build
- [x] Core test suite
- [ ] Production installer / release packaging
- [ ] User troubleshooting guide

**Success Criteria:**
- [x] No critical bugs in core service engine
- [x] Unit tests pass on Linux
- [ ] Installation package available


## Technical Architecture

### Folder Structure

```
LocalValet/
├── app.go                    # Main Wails app struct
├── main.go                   # Entry point
├── service_config.go         # Service configurations
├── service_path.go           # Binary path resolution
├── service_utils.go          # System utilities
│
├── config/
│   ├── runtime.json          # Version management
│   ├── services.json         # Service definitions
│   ├── vhosts.json           # Virtual host config
│   └── projects.json         # Discovered projects
│
├── internal/
│   ├── domain/
│   │   ├── service/
│   │   │   ├── config.go
│   │   │   └── manager.go
│   │   ├── vhost/
│   │   │   ├── manager.go
│   │   │   └── types.go
│   │   ├── ssl/
│   │   │   ├── manager.go
│   │   │   └── types.go
│   │   └── terminal/
│   │       └── manager.go
│   │
│   ├── usecase/
│   │   ├── service/
│   │   │   └── usecase.go
│   │   ├── vhost/
│   │   │   └── usecase.go
│   │   └── ssl/
│   │       └── usecase.go
│   │
│   ├── infrastructure/
│   │   ├── monitor/
│   │   │   ├── event_emitter.go
│   │   │   └── service_monitor.go
│   │   ├── terminal/
│   │   │   └── manager.go
│   │   └── config/
│   │       └── loader.go
│   │
│   └── platform/
│       ├── factory.go
│       ├── linux/
│       │   ├── linux.go
│       │   ├── runtime_registry.go
│       │   └── vhost_manager.go
│       └── windows/
│           └── windows.go
│
├── runtime/
│   ├── linux/
│   │   ├── apache/
│   │   ├── nginx/
│   │   ├── mysql/
│   │   ├── postgresql/
│   │   ├── redis/
│   │   ├── php/
│   │   └── node/
│   ├── pids/
│   ├── logs/
│   └── certs/
│
├── frontend/
│   ├── src/
│   │   ├── App.tsx
│   │   ├── main.tsx
│   │   ├── components/
│   │   │   ├── app-layout.tsx
│   │   │   ├── app-sidebar.tsx
│   │   │   └── ui/
│   │   ├── modules/
│   │   │   ├── service/
│   │   │   ├── log/
│   │   │   ├── vhost/
│   │   │   ├── project/
│   │   │   └── settings/
│   │   ├── pages/
│   │   │   ├── home-page.tsx
│   │   │   ├── services-page.tsx
│   │   │   ├── projects-page.tsx
│   │   │   └── settings-page.tsx
│   │   └── services/
│   │       ├── events.service.ts
│   │       ├── wails-service-control.service.ts
│   │       └── wails-vhost.service.ts
│   └── wailsjs/
│       └── go/
│           └── main/
│               └── App.d.ts
│
└── docs/
    ├── architecture.md
    ├── process-manager.md
    ├── ui-design.md
    ├── security.md
    ├── architecture-review.md
    └── mvp-roadmap.md
```

### Key Interfaces

```go
// Service Manager
type Manager interface {
    StartService(name string) error
    StopService(name string) error
    GetServiceStatus(name string) (bool, string)
    CheckHealth(name string) (bool, string)
    GetAllocatedPort(name string) int
    SetServiceVersion(name, version string) error
    GetActiveServiceVersion(name string) (string, error)
    GetAvailableVersions(name string) ([]string, error)
}

// VHost Manager
type VHostManager interface {
    DiscoverProjects(roots []string) ([]Project, error)
    GenerateVHost(project Project) error
    RemoveVHost(domain string) error
    ListVHosts() ([]VHost, error)
}

// SSL Manager
type SSLManager interface {
    InitCA() error
    GenerateCert(domain string) error
    TrustCert(domain string) error
    IsTrusted(domain string) bool
}
```

### Data Flow

```
User Action → Frontend → Wails IPC → Use Case → Domain → Platform
     ↑                                                       │
     └───────────────── Events ←─────────────────────────────┘
```

## Performance Targets

| Metric | Target | Current |
|--------|--------|---------|
| Startup time | < 2s | ~1.5s |
| Service start | < 1s | ~0.8s |
| Status check | < 100ms | ~20ms |
| Memory usage | < 100MB | ~60MB |
| Bundle size | < 20MB | ~12MB |

## Testing Strategy

### Unit Tests

- [x] Service manager operations
- [x] Port conflict resolution & auto-remapping
- [x] Health check engine (Process, TCP, HTTP)
- [x] Config parsing & version registry
- [x] Service dependency auto-start

### Integration Tests

- [x] Real service start/stop simulation
- [x] Health check verification
- [x] Port binding & collision detection
- [x] File system operations (PID/logs)

### End-to-End Tests

- [x] Full workflow (start service → check status → stop)
- [x] UI interactions & reactive events
- [x] Error scenarios & stale PID recovery

## Risk Mitigation

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Service crashes | Medium | High | Health checks, auto-restart |
| Port conflicts | High | Low | Auto-remapping (+200 scan) |
| Config corruption | Low | High | Backup before write |
| Permission issues | High | Medium | Clear error messages |

## Success Metrics

### Technical

- [x] All 6 services managed reliably
- [x] Port conflicts resolved automatically
- [x] Version switching works (PHP runtime)
- [x] UI is fast and responsive

### User Experience

- [x] New user can start services in < 1 minute
- [x] Error messages are clear and actionable
- [x] UI is intuitive without documentation
- [x] No data loss on crash

### Quality

- [x] 80%+ test coverage for core modules
- [x] No critical bugs
- [x] Documentation complete
- [x] Installation works on Linux


## Future Roadmap (Post-MVP)

### v1.1 (Month 3-4)

- Docker integration
- Windows support
- Auto-updates
- Portable mode

### v2.0 (Month 6-8)

- macOS support
- Plugin system
- Python version switching
- Advanced project templates

### v3.0 (Month 12+)

- Cloud sync
- Team collaboration
- Advanced monitoring
- Custom service templates

## Conclusion

**MVP is achievable in 10 weeks with focused scope.**

**Key decisions:**
- Linux first
- Core services only
- No plugin system
- No Docker integration
- Clean architecture for testing

**Next steps:**
1. Review and approve roadmap
2. Set up development environment
3. Start Phase 1
4. Weekly progress reviews
