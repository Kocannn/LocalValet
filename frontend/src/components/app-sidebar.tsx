import * as React from "react"
import { NavLink } from "react-router-dom"
import { Code2, FolderKanban, Home, Server, Settings, Terminal } from "lucide-react"

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
} from "@/components/ui/sidebar"

// Navigation data with routing
const data = {
  navMain: [
    {
      title: "Main",
      items: [
        {
          title: "Home",
          url: "/",
          icon: Home,
        },
        {
          title: "Services",
          url: "/services",
          icon: Server,
        },
        {
          title: "Projects",
          url: "/projects",
          icon: FolderKanban,
        },
        {
          title: "Settings",
          url: "/settings",
          icon: Settings,
        },
      ],
    },
  ],
}


export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  return (
    <Sidebar {...props}>
      <SidebarHeader className="border-b border-sidebar-border/70 px-3 py-4">
        <div className="px-2">
          <p className="text-sm font-semibold tracking-tight">LocalValet</p>
          <p className="text-xs text-muted-foreground">Service control center</p>
        </div>
      </SidebarHeader>
      <SidebarContent>
        {data.navMain.map((group) => (
          <SidebarGroup key={group.title} className="pt-3">
            <SidebarGroupLabel>{group.title}</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {group.items.map((item) => (
                  <SidebarMenuItem key={item.title}>
                    <SidebarMenuButton asChild className="h-9 rounded-md">
                      <NavLink 
                        to={item.url}
                        className={({ isActive }) =>
                          isActive
                            ? 'bg-accent text-accent-foreground font-medium'
                            : 'text-sidebar-foreground/85 hover:text-sidebar-foreground'
                        }
                      >
                        {item.icon && <item.icon className="mr-2 h-4 w-4" />}
                        {item.title}
                      </NavLink>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        ))}
      </SidebarContent>
      <SidebarFooter className="border-t border-sidebar-border/70 p-3">
        <div className="flex flex-col gap-1.5 rounded-lg bg-sidebar-accent/50 p-2 text-xs">
          <div className="flex items-center justify-between text-muted-foreground font-medium">
            <span className="flex items-center gap-1">
              <Code2 className="h-3 w-3" /> Runtime
            </span>
            <NavLink to="/settings" className="hover:underline text-[11px] text-primary">
              Switch
            </NavLink>
          </div>
          <div className="flex flex-wrap gap-1 text-[11px]">
            <span className="rounded bg-background/80 px-1.5 py-0.5 font-mono text-muted-foreground">
              PHP 8.4
            </span>
            <span className="rounded bg-background/80 px-1.5 py-0.5 font-mono text-muted-foreground">
              Node 22
            </span>
          </div>
        </div>
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}

