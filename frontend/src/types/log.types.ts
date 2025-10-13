/**
 * Log Types
 * 
 * Defines all type definitions related to logging system.
 */

export type LogLevel = 'info' | 'warning' | 'error' | 'success';

export interface LogEntry {
  /** Timestamp of the log entry */
  timestamp: string;
  
  /** Severity level of the log */
  level: LogLevel;
  
  /** Log message content */
  message: string;
}
