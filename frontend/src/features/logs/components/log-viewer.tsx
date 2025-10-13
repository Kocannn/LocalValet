/**
 * LogViewer Component
 * 
 * Displays a scrollable list of log entries with auto-scroll.
 * Handles rendering optimization for large log lists.
 * 
 * Design Pattern: Presentational Component
 * Purpose: Display logs with performance optimizations
 */

import { useRef, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Trash2 } from 'lucide-react';
import { LogRow } from './log-row';
import type { LogEntry } from '../types';

interface LogViewerProps {
  /** Array of log entries to display */
  logs: LogEntry[];
  
  /** Callback when clear button is clicked */
  onClear: () => void;
}

/**
 * Log viewer component with auto-scroll
 */
export function LogViewer({ logs, onClear }: LogViewerProps) {
  const logsEndRef = useRef<HTMLDivElement>(null);

  // Auto-scroll to bottom when new log is added (debounced)
  useEffect(() => {
    const timeoutId = setTimeout(() => {
      if (logsEndRef.current) {
        requestAnimationFrame(() => {
          if (logsEndRef.current) {
            const container = logsEndRef.current.parentElement;
            if (container) {
              container.scrollTo({
                top: container.scrollHeight,
                behavior: 'smooth'
              });
            }
          }
        });
      }
    }, 100);

    return () => clearTimeout(timeoutId);
  }, [logs.length]);

  return (
    <Card className='flex-1 flex flex-col min-h-0'>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle>Server Logs</CardTitle>
        <Button
          variant="ghost"
          size="sm"
          onClick={onClear}
          className="h-8 w-8 p-0"
        >
          <Trash2 className="h-4 w-4" />
        </Button>
      </CardHeader>
      <CardContent className="flex flex-col flex-1 min-h-0">
        <div 
          className="bg-[#FAFAFA] rounded-lg p-4 flex-1 min-h-0 overflow-y-auto font-mono text-sm"
          style={{
            willChange: 'scroll-position',
            WebkitOverflowScrolling: 'touch',
            contain: 'strict',
            scrollBehavior: 'smooth',
            transform: 'translateZ(0)',
            backfaceVisibility: 'hidden',
            perspective: 1000
          }}
        >
          {logs.length === 0 ? (
            <div className="text-black text-center py-8">
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
}
