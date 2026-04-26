import { memo, useEffect, useRef } from 'react';
import { Trash2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import type { LogEntry } from '../domain/types';
import { LogRow } from './log-row';

interface LogViewerProps {
  logs: LogEntry[];
  onClear: () => void;
}

export const LogViewer = memo(function LogViewer({ logs, onClear }: LogViewerProps) {
  const logsEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const timeoutId = setTimeout(() => {
      if (logsEndRef.current) {
        requestAnimationFrame(() => {
          if (logsEndRef.current) {
            const container = logsEndRef.current.parentElement;
            if (container) {
              container.scrollTo({
                top: container.scrollHeight,
                behavior: 'smooth',
              });
            }
          }
        });
      }
    }, 100);

    return () => clearTimeout(timeoutId);
  }, [logs.length]);

  return (
    <Card className='flex-1 flex flex-col min-h-0 border-border/70 bg-card/95'>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-base md:text-lg">Server Logs</CardTitle>
        <Button
          variant="ghost"
          size="sm"
          onClick={onClear}
          className="h-8 w-8 p-0 text-muted-foreground hover:text-foreground"
        >
          <Trash2 className="h-4 w-4" />
        </Button>
      </CardHeader>
      <CardContent className="flex flex-col flex-1 min-h-0">
        <div className="rounded-lg border border-border/70 bg-muted/35 p-4 flex-1 min-h-0 overflow-y-auto font-mono text-sm [scroll-behavior:smooth]">
          {logs.length === 0 ? (
            <div className="py-8 text-center text-sm text-muted-foreground">
              No logs available. Waiting for server activity...
            </div>
          ) : (
            <div className="space-y-1">
              {logs.map((log, index) => (
                <LogRow
                  key={`log-${index}-${log.timestamp}`}
                  log={log}
                />
              ))}
              <div ref={logsEndRef} />
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
});

LogViewer.displayName = 'LogViewer';
