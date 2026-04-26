import { memo } from 'react';
import type { LogEntry, LogLevel } from '../domain/types';

interface LogRowProps {
  log: LogEntry;
}

function getLogColor(level: LogLevel): string {
  switch (level) {
    case 'info':
      return 'text-sky-600';
    case 'success':
      return 'text-emerald-600';
    case 'warning':
      return 'text-amber-600';
    case 'error':
      return 'text-rose-600';
    default:
      return 'text-slate-500';
  }
}

export const LogRow = memo(({ log }: LogRowProps) => (
  <div
    className="flex flex-wrap gap-2 leading-tight py-0.5 items-start"
    style={{
      contain: 'layout style paint',
      contentVisibility: 'auto',
    }}
  >
    <span className="text-slate-500 flex-shrink-0 text-xs">[{log.timestamp}]</span>
    <span className={`font-semibold flex-shrink-0 text-xs ${getLogColor(log.level)}`}>
      [{log.level.toUpperCase()}]
    </span>
    <span className="text-slate-800 flex-1 min-w-0 text-xs break-words whitespace-pre-line">{log.message}</span>
  </div>
));

LogRow.displayName = 'LogRow';
