# LocalValet v2 - Development Phases

## Phase 1: Foundation (Week 1-2)

### Week 1: Service Manager Core

**Day 1-2: Service Manager Interface**
- Define `Manager` interface in `internal/domain/service/manager.go`
- Implement `LinuxManager` in `internal/platform/linux/linux.go`
- Add health check methods per service type

**Day 3-4: Port Management**
- Implement port conflict detection
- Add auto-remapping logic
- Add config file patching for nginx, php-fpm

**Day 5: PID Management**
- Implement PID file read/write
- Add process alive check
- Add graceful shutdown with timeout

### Week 2: Testing & Polish

**Day 1-2: Unit Tests**
- Test service start/stop
- Test port conflict resolution
- Test config parsing

**Day 3-4: Integration Tests**
- Test real service start/stop
- Test health checks
- Test file system operations

**Day 5: Bug Fixes**
- Fix any issues found in testing
- Polish error messages
- Update documentation

### Deliverables

- [ ] Service manager with health checks
- [ ] Port management with auto-remapping
- [ ] Unit tests for all operations
- [ ] Integration tests for real services

### Success Criteria

- All 6 services start/stop reliably
- Port conflicts resolved automatically
- Tests pass on Linux

---

## Phase 2: UI Polish (Week 3-4)

### Week 3: Service Control Panel

**Day 1-2: Service Card Component**
- Create `ServiceCard` component
- Add status indicators (running/stopped/error)
- Add toggle switch
- Add quick actions menu

**Day 3-4: Service Grouping**
- Implement service grouping (Web, Database, Runtime)
- Add group headers
- Add group-level actions

**Day 5: Log Viewer**
- Implement `LogViewer` component
- Add color-coded log levels
- Add auto-scroll
- Add clear button

### Week 4: Layout & Polish

**Day 1-2: Responsive Layout**
- Implement sidebar navigation
- Add mobile responsive breakpoints
- Add dark mode support

**Day 3-4: Settings Page**
- Create settings page layout
- Add general settings section
- Add service settings section

**Day 5: UI Polish**
- Add animations
- Improve error states
- Add loading states

### Deliverables

- [ ] Service control panel
- [ ] Log viewer with filtering
- [ ] Responsive sidebar navigation
- [ ] Dark mode support

### Success Criteria

- UI is fast and responsive
- Dark mode works correctly
- Mobile layout is usable

---

## Phase 3: Version Switching (Week 5-6)

### Week 5: Version Management

**Day 1-2: Version Detection**
- Implement version detection per service
- Add version list from runtime directory
- Add active version tracking

**Day 3-4: Version Switching**
- Implement hot-switching logic
- Add service restart on version change
- Add binary verification

**Day 5: Version UI**
- Add version selector in header
- Add per-service version dropdown
- Add active version indicator

### Week 6: Multi-Version Support

**Day 1-2: PHP Versions**
- Support PHP 8.1, 8.2, 8.3, 8.4
- Add PHP-specific config handling
- Test version switching

**Day 3-4: Node.js Versions**
- Support Node.js 18, 20, 22
- Add Node.js-specific config handling
- Test version switching

**Day 5: Testing**
- Test all version combinations
- Fix any compatibility issues
- Update documentation

### Deliverables

- [ ] Version switcher in header
- [ ] Per-service version management
- [ ] Version hot-switching
- [ ] PHP 8.1-8.4 support
- [ ] Node.js 18-22 support

### Success Criteria

- PHP version switching works
- Node.js version switching works
- No service crashes during switch

---

## Phase 4: Project Discovery (Week 7-8)

### Week 7: Project Scanner

**Day 1-2: Project Detection**
- Implement project scanner
- Detect frameworks (Laravel, WordPress, Next.js, etc.)
- Add project metadata extraction

**Day 3-4: Project UI**
- Create project grid/list view
- Add project cards with framework icons
- Add quick actions (terminal, browser, IDE)

**Day 5: Project Storage**
- Implement project persistence
- Add project search/filter
- Add project favorites

### Week 8: Virtual Hosts & SSL

**Day 1-2: Virtual Host Manager**
- Implement VHost config generation
- Add nginx/apache templates
- Add DNS resolution setup

**Day 3-4: SSL Certificate Manager**
- Implement local CA generation
- Add per-domain cert generation
- Add cert trust management

**Day 5: Integration**
- Integrate VHost with project discovery
- Test full workflow
- Fix any issues

### Deliverables

- [ ] Project discovery engine
- [ ] Virtual host manager
- [ ] SSL certificate generator
- [ ] Project grid/list view

### Success Criteria

- Projects discovered automatically
- VHosts generated correctly
- SSL certificates trusted

---

## Phase 5: Polish & Release (Week 9-10)

### Week 9: Performance & Testing

**Day 1-2: Performance Optimization**
- Optimize startup time
- Reduce memory usage
- Improve UI responsiveness

**Day 3-4: Error Handling**
- Improve error messages
- Add error recovery
- Add user-friendly error dialogs

**Day 5: User Testing**
- Conduct user testing
- Gather feedback
- Prioritize fixes

### Week 10: Release Preparation

**Day 1-2: Bug Fixes**
- Fix critical bugs
- Fix UI issues
- Fix edge cases

**Day 3-4: Documentation**
- Write user documentation
- Write installation guide
- Write troubleshooting guide

**Day 5: Release**
- Create release build
- Test installation
- Publish release

### Deliverables

- [ ] Production build
- [ ] User documentation
- [ ] Installation guide
- [ ] Troubleshooting guide

### Success Criteria

- No critical bugs
- Documentation complete
- Installation works on Linux

---

## Parallel Workstreams

### Backend Development

```
Week 1-2: Service Manager
Week 3-4: VHost Manager
Week 5-6: Version Manager
Week 7-8: SSL Manager
Week 9-10: Integration & Polish
```

### Frontend Development

```
Week 1-2: Basic UI Shell
Week 3-4: Service Control Panel
Week 5-6: Version Switcher UI
Week 7-8: Project Discovery UI
Week 9-10: UI Polish
```

### Testing

```
Week 1-2: Unit Tests (Service Manager)
Week 3-4: Integration Tests (UI)
Week 5-6: Version Switching Tests
Week 7-8: Project Discovery Tests
Week 9-10: End-to-End Tests
```

---

## Dependencies

```
Phase 1 (Foundation)
    ↓
Phase 2 (UI Polish) ←── Can start in parallel
    ↓
Phase 3 (Version Switching) ←── Depends on Phase 1
    ↓
Phase 4 (Project Discovery) ←── Depends on Phase 1
    ↓
Phase 5 (Polish & Release) ←── Depends on all phases
```

---

## Resource Allocation

### Solo Developer (Recommended for MVP)

- 100% on backend (Week 1-2)
- 50% backend, 50% frontend (Week 3-8)
- 100% on polish & testing (Week 9-10)

### Small Team (2-3 developers)

- 2 backend, 1 frontend (Week 1-2)
- 1 backend, 2 frontend (Week 3-8)
- All on polish & testing (Week 9-10)

---

## Risk Management

### Technical Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Service crashes | Medium | High | Health checks, auto-restart |
| Port conflicts | High | Low | Auto-remapping |
| Config corruption | Low | High | Backup before write |
| Permission issues | High | Medium | Clear error messages |

### Schedule Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Scope creep | High | High | Strict MVP scope |
| Technical debt | Medium | Medium | Refactor in Phase 5 |
| Testing time | Medium | Medium | Test as you go |
| Integration issues | Medium | High | Early integration |

---

## Quality Gates

### Phase Completion Criteria

**Phase 1 Complete:**
- [ ] All services start/stop reliably
- [ ] Port conflicts resolved
- [ ] Tests pass

**Phase 2 Complete:**
- [ ] UI is responsive
- [ ] Dark mode works
- [ ] Mobile layout usable

**Phase 3 Complete:**
- [ ] Version switching works
- [ ] No crashes during switch
- [ ] PHP and Node supported

**Phase 4 Complete:**
- [ ] Projects discovered
- [ ] VHosts generated
- [ ] SSL certificates trusted

**Phase 5 Complete:**
- [ ] No critical bugs
- [ ] Documentation complete
- [ ] Installation works

---

## Communication

### Weekly Sync

- Progress review
- Blocker discussion
- Priority adjustment
- Demo of completed work

### Documentation

- Architecture decisions
- API documentation
- User guides
- Troubleshooting guides

---

## Conclusion

**MVP is achievable in 10 weeks with focused scope.**

**Key success factors:**
1. Strict MVP scope
2. Test-driven development
3. Regular user feedback
4. Early integration
5. Performance focus

**Next steps:**
1. Approve roadmap
2. Set up development environment
3. Start Phase 1
4. Weekly progress reviews
