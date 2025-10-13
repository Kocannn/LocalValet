/**
 * Logs Feature - Public API
 * 
 * This is the ONLY file that should be imported from outside this feature.
 * Provides a clean public interface while hiding internal implementation.
 * 
 * Import pattern: import { LogViewer, useLogs } from '@/features/logs'
 */

export { LogViewer, LogRow } from './components';
export { useLogs } from './hooks';
export type { LogEntry, LogLevel } from './types';
