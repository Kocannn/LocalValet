import { memo } from 'react';
import {
  ExternalLink,
  Folder,
  Globe,
  Lock,
  Terminal,
  Code,
  ShieldCheck,
  SwitchCamera,
} from 'lucide-react';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Switch } from '@/components/ui/switch';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip';
import type { Project } from '../domain/types';
import { FrameworkBadge } from './framework-icon';

interface ProjectCardProps {
  project: Project;
  onOpenBrowser: (url: string) => void;
  onOpenTerminal: (path: string) => void;
  onOpenEditor: (path: string) => void;
  onToggleVHost: (path: string, enable: boolean) => void;
  onGenerateSSL: (path: string) => void;
}

export const ProjectCard = memo(
  ({
    project,
    onOpenBrowser,
    onOpenTerminal,
    onOpenEditor,
    onToggleVHost,
    onGenerateSSL,
  }: ProjectCardProps) => {
    const url = `https://${project.domain}`;
    const httpUrl = `http://${project.domain}`;

    return (
      <TooltipProvider>
        <Card className="flex flex-col justify-between border-border/70 bg-card/95 transition-all hover:border-border hover:shadow-sm">
          <CardHeader className="pb-3">
            <div className="flex items-start justify-between gap-2">
              <div className="space-y-1 overflow-hidden">
                <CardTitle className="text-base font-semibold truncate" title={project.name}>
                  {project.name}
                </CardTitle>
                <div className="flex flex-wrap items-center gap-1.5 pt-0.5">
                  <FrameworkBadge framework={project.framework} />
                  {project.sslEnabled ? (
                    <Badge
                      variant="outline"
                      className="border-emerald-200 bg-emerald-50 text-emerald-700 dark:border-emerald-800 dark:bg-emerald-950/40 dark:text-emerald-400 gap-1 text-[11px] font-normal cursor-pointer"
                      onClick={() => onGenerateSSL(project.path)}
                      title="SSL Active (click to regenerate)"
                    >
                      <ShieldCheck className="h-3 w-3" />
                      SSL
                    </Badge>
                  ) : (
                    <Badge
                      variant="outline"
                      className="text-[11px] font-normal text-muted-foreground cursor-pointer"
                      onClick={() => onGenerateSSL(project.path)}
                    >
                      No SSL
                    </Badge>
                  )}
                  {project.targetPort && (
                    <Badge variant="secondary" className="text-[11px] font-mono font-normal">
                      :{project.targetPort}
                    </Badge>
                  )}
                </div>
              </div>

              {/* VHost Active Switch */}
              <div className="flex flex-col items-end gap-1">
                <Tooltip>
                  <TooltipTrigger asChild>
                    <div className="flex items-center gap-1.5">
                      <Switch
                        checked={project.vhostEnabled}
                        onCheckedChange={(checked) => onToggleVHost(project.path, checked)}
                      />
                    </div>
                  </TooltipTrigger>
                  <TooltipContent>
                    <p className="text-xs">{project.vhostEnabled ? 'VHost Active' : 'VHost Disabled'}</p>
                  </TooltipContent>
                </Tooltip>
              </div>
            </div>

            <CardDescription className="pt-2">
              <button
                type="button"
                onClick={() => onOpenBrowser(project.sslEnabled ? url : httpUrl)}
                className="group flex items-center gap-1 font-mono text-xs text-primary hover:underline"
              >
                <Globe className="h-3.5 w-3.5 text-muted-foreground group-hover:text-primary" />
                <span>{project.domain}</span>
                <ExternalLink className="h-3 w-3 opacity-60 group-hover:opacity-100" />
              </button>
            </CardDescription>
          </CardHeader>

          <CardContent className="py-0">
            <div className="flex items-center gap-1 text-[11px] text-muted-foreground truncate" title={project.path}>
              <Folder className="h-3 w-3 shrink-0" />
              <span className="truncate">{project.path}</span>
            </div>
            {project.webRoot && project.webRoot !== '.' && (
              <div className="text-[10px] text-muted-foreground/70 font-mono mt-0.5">
                root: /{project.webRoot}
              </div>
            )}
          </CardContent>

          <CardFooter className="pt-4 pb-3 flex items-center justify-between border-t border-border/40 mt-3">
            <div className="flex items-center gap-1">
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant="outline"
                    size="sm"
                    className="h-7 px-2 text-xs gap-1"
                    onClick={() => onOpenBrowser(project.sslEnabled ? url : httpUrl)}
                  >
                    <Globe className="h-3.5 w-3.5" />
                    Open
                  </Button>
                </TooltipTrigger>
                <TooltipContent>
                  <p className="text-xs">Open in Browser ({project.domain})</p>
                </TooltipContent>
              </Tooltip>

              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant="outline"
                    size="sm"
                    className="h-7 px-2 text-xs gap-1"
                    onClick={() => onOpenTerminal(project.path)}
                  >
                    <Terminal className="h-3.5 w-3.5" />
                    Terminal
                  </Button>
                </TooltipTrigger>
                <TooltipContent>
                  <p className="text-xs">Open Context Terminal with Injected PATH</p>
                </TooltipContent>
              </Tooltip>

              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant="outline"
                    size="sm"
                    className="h-7 px-2 text-xs gap-1"
                    onClick={() => onOpenEditor(project.path)}
                  >
                    <Code className="h-3.5 w-3.5" />
                    IDE
                  </Button>
                </TooltipTrigger>
                <TooltipContent>
                  <p className="text-xs">Open in Code Editor</p>
                </TooltipContent>
              </Tooltip>
            </div>
          </CardFooter>
        </Card>
      </TooltipProvider>
    );
  },
);

ProjectCard.displayName = 'ProjectCard';
