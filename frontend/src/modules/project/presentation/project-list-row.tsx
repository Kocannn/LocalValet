import { memo } from 'react';
import {
  Code,
  ExternalLink,
  Folder,
  Globe,
  ShieldCheck,
  Terminal,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Switch } from '@/components/ui/switch';
import { Badge } from '@/components/ui/badge';
import type { Project } from '../domain/types';
import { FrameworkBadge } from './framework-icon';

interface ProjectListRowProps {
  project: Project;
  onOpenBrowser: (url: string) => void;
  onOpenTerminal: (path: string) => void;
  onOpenEditor: (path: string) => void;
  onToggleVHost: (path: string, enable: boolean) => void;
  onGenerateSSL: (path: string) => void;
}

export const ProjectListRow = memo(
  ({
    project,
    onOpenBrowser,
    onOpenTerminal,
    onOpenEditor,
    onToggleVHost,
    onGenerateSSL,
  }: ProjectListRowProps) => {
    const url = `https://${project.domain}`;
    const httpUrl = `http://${project.domain}`;

    return (
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 rounded-lg border border-border/70 bg-card/95 p-3 transition-all hover:border-border">
        {/* Project Info */}
        <div className="flex items-center gap-3 min-w-0">
          <FrameworkBadge framework={project.framework} />
          <div className="min-w-0 space-y-0.5">
            <div className="flex items-center gap-2">
              <span className="font-semibold text-sm truncate" title={project.name}>
                {project.name}
              </span>
              {project.sslEnabled ? (
                <Badge
                  variant="outline"
                  className="border-emerald-200 bg-emerald-50 text-emerald-700 dark:border-emerald-800 dark:bg-emerald-950/40 dark:text-emerald-400 gap-0.5 text-[10px] h-4 px-1"
                  onClick={() => onGenerateSSL(project.path)}
                >
                  <ShieldCheck className="h-2.5 w-2.5" />
                  SSL
                </Badge>
              ) : null}
            </div>
            <div className="flex items-center gap-2 text-xs text-muted-foreground">
              <button
                type="button"
                onClick={() => onOpenBrowser(project.sslEnabled ? url : httpUrl)}
                className="font-mono hover:underline text-primary flex items-center gap-1"
              >
                {project.domain}
                <ExternalLink className="h-2.5 w-2.5" />
              </button>
              <span>•</span>
              <span className="truncate max-w-[200px] text-[11px]" title={project.path}>
                {project.path}
              </span>
            </div>
          </div>
        </div>

        {/* Actions & VHost Toggle */}
        <div className="flex items-center justify-between sm:justify-end gap-2 shrink-0 border-t sm:border-t-0 pt-2 sm:pt-0 border-border/40">
          <div className="flex items-center gap-1">
            <Button
              variant="outline"
              size="sm"
              className="h-7 px-2 text-xs gap-1"
              onClick={() => onOpenBrowser(project.sslEnabled ? url : httpUrl)}
            >
              <Globe className="h-3 w-3" />
              Open
            </Button>
            <Button
              variant="outline"
              size="sm"
              className="h-7 px-2 text-xs gap-1"
              onClick={() => onOpenTerminal(project.path)}
            >
              <Terminal className="h-3 w-3" />
              Terminal
            </Button>
            <Button
              variant="outline"
              size="sm"
              className="h-7 px-2 text-xs gap-1"
              onClick={() => onOpenEditor(project.path)}
            >
              <Code className="h-3 w-3" />
              IDE
            </Button>
          </div>

          <div className="flex items-center gap-1.5 ml-2 border-l border-border/60 pl-2">
            <Switch
              checked={project.vhostEnabled}
              onCheckedChange={(checked) => onToggleVHost(project.path, checked)}
            />
          </div>
        </div>
      </div>
    );
  },
);

ProjectListRow.displayName = 'ProjectListRow';
