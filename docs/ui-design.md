# LocalValet v2 - UI/UX Design Specification

## Layout Architecture

### Main Dashboard Layout

```
┌─────────────────────────────────────────────────────────────────┐
│ [≡] LocalValet                              [PHP 8.4 ▼] [🌙]  │
├──────────┬──────────────────────────────────────────────────────┤
│          │                                                      │
│ 🏠 Home  │  ┌─────────────────────────────────────────────┐   │
│          │  │  Server Modules                    [● Live]  │   │
│ ⚙️ Serv  │  ├─────────────────────────────────────────────┤   │
│          │  │  Apache   :8080  [● Active]    [━━━●━━━]    │   │
│ 📁 Proj  │  │  MySQL    :3306  [● Active]    [━━━●━━━]    │   │
│          │  │  Redis    :6379  [● Active]    [━━━●━━━]    │   │
│ 🔧 Sett  │  │  PHP-FPM  :9074  [● Active]    [━━━●━━━]    │   │
│          │  │  Nginx    :8080  [○ Inactive]  [━━━○━━━]    │   │
│          │  └─────────────────────────────────────────────┘   │
│          │                                                      │
│          │  ┌─────────────────────┐ ┌─────────────────────┐   │
│          │  │ Projects (12)       │ │ Server Logs         │   │
│          │  │ ┌───┐ ┌───┐ ┌───┐  │ │ [10:30:01] [INFO]   │   │
│          │  │ │🌐 │ │🛒 │ │📝 │  │ │ [10:30:02] [OK]     │   │
│          │  │ │Lrv│ │Wrd│ │Nxt│  │ │ [10:30:03] [WARN]   │   │
│          │  │ └───┘ └───┘ └───┘  │ │                       │   │
│          │  └─────────────────────┘ └─────────────────────┘   │
│          │                                                      │
└──────────┴──────────────────────────────────────────────────────┘
```

### Responsive Breakpoints

| Breakpoint | Layout Change |
|------------|---------------|
| < 768px    | Sidebar collapses to icons, single column |
| 768-1024px | Sidebar icons + labels, 2 column grid |
| > 1024px   | Full sidebar, 2-3 column grid |

### Header Components

```
┌─────────────────────────────────────────────────────────────────┐
│ [☰] LocalValet                    [PHP 8.4 ▼] [Node 20 ▼] [⚙️] │
└─────────────────────────────────────────────────────────────────┘
     │                                    │          │       │
     │                                    │          │       └─ Settings
     │                                    │          └─ Version selector
     │                                    └─ Version selector
     └─ Sidebar toggle
```

## Service Control Panel

### Service Card Design

```
┌─────────────────────────────────────────────────────────────┐
│  Apache                                    [●] [⋮]         │
│  localhost:8080                                              │
│  ─────────────────────────────────────────────────────────  │
│  Status: Running    Uptime: 2h 15m    Memory: 45MB         │
│  ─────────────────────────────────────────────────────────  │
│  [Restart] [Logs] [Config] [Terminal]                       │
└─────────────────────────────────────────────────────────────┘
```

### Service States

| State | Visual Indicator |
|-------|------------------|
| Running | Green badge, toggle ON |
| Stopped | Gray badge, toggle OFF |
| Starting | Yellow badge, spinner |
| Error | Red badge, error icon |
| Port Conflict | Orange badge, warning icon |

### Service Grouping

```
┌─ Web Servers ──────────────────────────────────────────────┐
│  [Apache] [Nginx]                                          │
└────────────────────────────────────────────────────────────┘

┌─ Databases ────────────────────────────────────────────────┐
│  [MySQL] [PostgreSQL] [Redis]                              │
└────────────────────────────────────────────────────────────┘

┌─ Runtimes ─────────────────────────────────────────────────┐
│  [PHP-FPM] [Node.js] [Python]                             │
└────────────────────────────────────────────────────────────┘
```

## Project Discovery View

### Grid View

```
┌─────────────────────────────────────────────────────────────┐
│ Projects                                    [Grid] [List]  │
├─────────────────────────────────────────────────────────────┤
│ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│ │   🌐        │ │   🛒        │ │   📝        │           │
│ │  Laravel    │ │  WordPress  │ │  Next.js    │           │
│ │ myproject   │ │  blog       │ │  frontend   │           │
│ │ PHP 8.4     │ │ PHP 8.2     │ │ Node 20     │           │
│ │             │ │             │ │             │           │
│ │ [🖥️][🌐][⚙️]│ │ [🖥️][🌐][⚙️]│ │ [🖥️][🌐][⚙️]│           │
│ └─────────────┘ └─────────────┘ └─────────────┘           │
│                                                             │
│ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐           │
│ │   🐍        │ │   🐳        │ │   📦        │           │
│ │  Django     │ │  Docker     │ │  Symfony    │           │
│ │  api        │ │  stack      │ │  app        │           │
│ │ Python 3.11 │ │  Compose    │ │ PHP 8.3     │           │
│ │             │ │             │ │             │           │
│ │ [🖥️][🌐][⚙️]│ │ [🖥️][🌐][⚙️]│ │ [🖥️][🌐][⚙️]│           │
│ └─────────────┘ └─────────────┘ └─────────────┘           │
└─────────────────────────────────────────────────────────────┘

Legend: [🖥️] Terminal  [🌐] Browser  [⚙️] Settings
```

### List View

```
┌─────────────────────────────────────────────────────────────┐
│ Projects                                    [Grid] [List]  │
├─────────────────────────────────────────────────────────────┤
│ 🌐 myproject        Laravel   PHP 8.4  /home/user/proj... │
│ 🛒 blog             WordPress PHP 8.2  /home/user/blog... │
│ 📝 frontend         Next.js   Node 20  /home/user/fron... │
│ 🐍 api              Django    Py 3.11  /home/user/api...  │
│ 🐳 stack            Docker    Compose  /home/user/stac... │
└─────────────────────────────────────────────────────────────┘
```

### Project Card Actions

| Action | Icon | Description |
|--------|------|-------------|
| Open Terminal | `>_` | Launch terminal in project directory |
| Open Browser | 🌐 | Open http://project.test |
| Open in IDE | ✏️ | Open in VS Code / default editor |
| Settings | ⚙️ | Project-specific configuration |

## Version Switcher UI

### Global Version Selector (Header)

```
┌─────────────────────────────────────────────────────────────┐
│ [PHP 8.4 ▼] [Node 20 ▼] [Python 3.11 ▼]                   │
└─────────────────────────────────────────────────────────────┘
         │
         ▼
┌──────────────┐
│ PHP 8.4  ✓  │
│ PHP 8.3     │
│ PHP 8.2     │
│ PHP 8.1     │
│ ─────────── │
│ Manage...   │
└──────────────┘
```

### Per-Service Version (Settings)

```
┌─────────────────────────────────────────────────────────────┐
│ PHP Runtime                                                 │
│ ─────────────────────────────────────────────────────────  │
│ Active Version: [PHP 8.4 ▼]                                │
│                                                             │
│ Available Versions:                                         │
│ ✓ PHP 8.4 (active)                                          │
│   PHP 8.3                                                   │
│   PHP 8.2                                                   │
│   PHP 8.1                                                   │
│                                                             │
│ [Download New Version]                                      │
└─────────────────────────────────────────────────────────────┘
```

## Log Viewer

### Log Panel Design

```
┌─────────────────────────────────────────────────────────────┐
│ Server Logs                            [Filter ▼] [🗑️ Clear]│
├─────────────────────────────────────────────────────────────┤
│ [10:30:01] [INFO]    LocalValet started on linux            │
│ [10:30:02] [SUCCESS]  Apache started on port 8080           │
│ [10:30:02] [SUCCESS]  MySQL started on port 3306            │
│ [10:30:03] [WARNING]  Port 8080 conflict, using 8081        │
│ [10:30:04] [ERROR]    Redis failed to start: permission     │
│ [10:30:05] [INFO]    PHP-FPM started on port 9074           │
│                                                            │
│                                                            │
│                                                            │
│ ─────────────────────────────────────────────────────────── │
│ [Auto-scroll □]  Showing 6 of 142 entries                   │
└─────────────────────────────────────────────────────────────┘
```

### Log Level Colors

| Level | Color | Hex |
|-------|-------|-----|
| INFO | Sky Blue | `text-sky-600` |
| SUCCESS | Emerald | `text-emerald-600` |
| WARNING | Amber | `text-amber-600` |
| ERROR | Rose | `text-rose-600` |

### Log Filtering

```
┌─────────────────────────────────────────────────────────────┐
│ [All ▼] [Apache] [MySQL] [Redis] [PHP-FPM] [System]       │
└─────────────────────────────────────────────────────────────┘
```

## Terminal Integration

### Embedded Terminal

```
┌─────────────────────────────────────────────────────────────┐
│ Terminal                                    [+] [×]         │
├─────────────────────────────────────────────────────────────┤
│ LocalValet Active | PHP 8.4 | MySQL 8.0 | Nginx 1.26      │
│                                                            │
│ user@localhost:~/projects/myproject$ php artisan serve      │
│ Starting development server: http://127.0.0.1:8000         │
│ Listening on http://127.0.0.1:8000                         │
│ Press Ctrl+C to stop the server                            │
│                                                            │
│ user@localhost:~/projects/myproject$ █                     │
└─────────────────────────────────────────────────────────────┘
```

### Terminal Features

- **Context-aware**: Auto PATH, env vars injected
- **Multiple tabs**: Each with independent shell
- **Themes**: Match system dark/light mode
- **Project-aware**: Auto-cd to project directory

## Settings Page

### Settings Layout

```
┌─────────────────────────────────────────────────────────────┐
│ Settings                                                    │
├─────────────────────────────────────────────────────────────┤
│ ┌─ General ────────────────────────────────────────────────┐│
│ │ Start on boot: [□]                                        ││
│ │ Minimize to tray: [☑]                                     ││
│ │ Dark mode: [☑]                                            ││
│ │ Default project root: [/home/user/projects] [Browse]      ││
│ └──────────────────────────────────────────────────────────┘│
│                                                             │
│ ┌─ Services ───────────────────────────────────────────────┐│
│ │ PHP Version: [8.4 ▼]                                      ││
│ │ Node Version: [20 ▼]                                      ││
│ │ MySQL Port: [3306]                                         ││
│ │ Apache Port: [8080]                                        ││
│ └──────────────────────────────────────────────────────────┘│
│                                                             │
│ ┌─ Virtual Hosts ──────────────────────────────────────────┐│
│ │ Auto-generate: [☑]                                        ││
│ │ Domain suffix: [.test]                                     ││
│ │ Auto SSL: [☑]                                             ││
│ └──────────────────────────────────────────────────────────┘│
│                                                             │
│ ┌─ Docker ─────────────────────────────────────────────────┐│
│ │ Docker path: [/usr/bin/docker]                             ││
│ │ Auto-detect compose: [☑]                                  ││
│ └──────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────┘
```

## Design System

### Color Palette

**Light Mode:**
```css
--background: oklch(0.985 0.006 255);      /* Near white */
--foreground: oklch(0.22 0.02 255);         /* Dark slate */
--primary: oklch(0.5 0.15 250);             /* Blue */
--success: oklch(0.6 0.15 150);             /* Emerald */
--warning: oklch(0.7 0.15 80);              /* Amber */
--destructive: oklch(0.577 0.245 27.325);   /* Rose */
```

**Dark Mode:**
```css
--background: oklch(0.145 0 0);             /* Near black */
--foreground: oklch(0.985 0 0);             /* Near white */
--primary: oklch(0.922 0 0);                /* Light gray */
```

### Typography Scale

| Element | Size | Weight |
|---------|------|--------|
| H1 | 1.5rem | 600 |
| H2 | 1.25rem | 600 |
| H3 | 1rem | 600 |
| Body | 0.875rem | 400 |
| Small | 0.75rem | 400 |
| Badge | 0.75rem | 500 |

### Spacing System

| Token | Value |
|-------|-------|
| xs | 0.25rem |
| sm | 0.5rem |
| md | 1rem |
| lg | 1.5rem |
| xl | 2rem |

### Component Library (shadcn/ui)

| Component | Usage |
|-----------|-------|
| Card | Service cards, project cards, settings sections |
| Badge | Status indicators, port numbers, versions |
| Switch | Service toggle |
| Button | Actions, navigation |
| Table | Service list, project list |
| Dropdown | Version selectors, filters |
| Sheet | Mobile sidebar, dialogs |
| Tooltip | Hover information |
| Separator | Visual dividers |

### Icons (lucide-react)

| Icon | Usage |
|------|-------|
| Home | Home page |
| Server | Services page |
| FolderOpen | Projects page |
| Settings | Settings page |
| Terminal | Terminal actions |
| Play/Stop | Service control |
| RefreshCw | Restart service |
| Trash2 | Clear logs |
| ExternalLink | Open in browser |
| Code | Open in IDE |
