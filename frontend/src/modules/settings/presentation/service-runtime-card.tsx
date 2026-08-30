import { memo } from 'react';
import { Check, ChevronDown, Loader2, Sparkles } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import type { RuntimeService } from '../domain/types';

interface ServiceRuntimeCardProps {
  service: RuntimeService;
  isSaving: boolean;
  error?: string;
  onVersionChange: (serviceName: string, version: string) => Promise<void>;
}

export const ServiceRuntimeCard = memo(
  ({ service, isSaving, error, onVersionChange }: ServiceRuntimeCardProps) => {
    const { displayName, serviceName, activeVersion, availableVersions, category, isRunning } = service;

    // Merge versions to ensure activeVersion is always visible
    const allVersions = Array.from(
      new Set([activeVersion, ...availableVersions].filter(Boolean)),
    );

    return (
      <Card className="border-border/70 bg-card/95 transition-all hover:border-border">
        <CardHeader className="pb-3">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <CardTitle className="text-base font-semibold">{displayName}</CardTitle>
              <Badge variant="outline" className="text-xs font-normal">
                {category}
              </Badge>
            </div>
            {isRunning && (
              <Badge variant="outline" className="border-emerald-200 bg-emerald-50 text-emerald-700 text-xs">
                Running
              </Badge>
            )}
          </div>
          <CardDescription className="text-xs text-muted-foreground">
            Canonical name: <code className="font-mono text-xs text-foreground/80">{serviceName}</code>
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-3">
          <div className="space-y-1.5">
            <div className="flex items-center justify-between text-xs text-muted-foreground">
              <span>Active Runtime</span>
              <span className="font-mono text-xs font-medium text-foreground">
                v{activeVersion || 'default'}
              </span>
            </div>

            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button
                  variant="outline"
                  className="w-full justify-between"
                  disabled={isSaving || allVersions.length <= 1}
                >
                  <span className="flex items-center gap-1.5 truncate">
                    {isSaving ? (
                      <>
                        <Loader2 className="h-3.5 w-3.5 animate-spin text-primary" />
                        Switching version...
                      </>
                    ) : (
                      <>
                        <Sparkles className="h-3.5 w-3.5 text-muted-foreground" />
                        Version {activeVersion || 'default'}
                      </>
                    )}
                  </span>
                  <ChevronDown className="h-4 w-4 opacity-50" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent className="w-[var(--radix-dropdown-menu-trigger-width)]" align="start">
                {allVersions.map((version) => (
                  <DropdownMenuItem
                    key={version}
                    className="flex items-center justify-between"
                    onSelect={() => onVersionChange(serviceName, version)}
                  >
                    <span>Version {version}</span>
                    {version === activeVersion && <Check className="h-4 w-4 text-emerald-600" />}
                  </DropdownMenuItem>
                ))}
              </DropdownMenuContent>
            </DropdownMenu>
          </div>

          {error && <p className="text-xs font-medium text-destructive">{error}</p>}

          <p className="text-[11px] text-muted-foreground">
            {isRunning
              ? 'Hot-switch active: service will automatically restart on change.'
              : 'Applies to the next service startup.'}
          </p>
        </CardContent>
      </Card>
    );
  },
);

ServiceRuntimeCard.displayName = 'ServiceRuntimeCard';
