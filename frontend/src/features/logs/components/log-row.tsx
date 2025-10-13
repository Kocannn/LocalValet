/**
 * LogRow Component
 * 
 * Single log entry display component.
 * Memoized to prevent unnecessary re-renders.
 * 
 * Design Pattern: Presentational Component (Memoized)
 * Purpose: Optimize rendering performance
 */

import { memo } from 'react';
import type { LogEntry, LogLevel } from '../types';

interface LogRowProps {
  log: LogEntry;
}

/**
 * Get color class for log level
 */
function getLogColor(level: LogLevel): string {
  switch (level) {
    case 'info':
      return 'text-blue-400';
    case 'success':
      return 'text-green-400';
    case 'warning':
      return 'text-yellow-400';
    case 'error':
      return 'text-red-400';
    default:
      return 'text-gray-400';
  }
}

/**
 * Individual log entry row
 */
export const LogRow = memo(({ log }: LogRowProps) => (
  <div 
    className="flex flex-wrap gap-2 leading-tight py-0.5 items-start"
    style={{ 
      contain: 'layout style paint',
      contentVisibility: 'auto'
    }}
  >
    <span className="text-black flex-shrink-0 text-xs">[{log.timestamp}]</span>
    <span className={`font-semibold flex-shrink-0 text-xs ${getLogColor(log.level)}`}> 
      [{log.level.toUpperCase()}]
    </span>
    <span className="text-black flex-1 min-w-0 text-xs break-words whitespace-pre-line">{log.message}</span>
  </div>
));

LogRow.displayName = 'LogRow';
