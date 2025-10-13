/**
 * useLogs Hook
 * 
 * Custom hook for managing log state and operations.
 * Handles log collection, limiting, and real-time updates.
 * 
 * Design Pattern: Custom Hook Pattern
 * Purpose: Encapsulate log management logic
 */

import { useState, useEffect, useCallback } from 'react';
import { subscribeToLogs } from '@/services/events.service';
import type { LogEntry } from '../types';

const MAX_LOGS = 50; // Maximum logs to keep in memory

interface UseLogsReturn {
  logs: LogEntry[];
  addLog: (log: LogEntry) => void;
  clearLogs: () => void;
}

/**
 * Hook for log management
 * @returns Log state and control functions
 */
export function useLogs(): UseLogsReturn {
  const [logs, setLogs] = useState<LogEntry[]>([
    {
      timestamp: new Date().toLocaleTimeString(),
      level: 'info',
      message: 'LocalValet initialized successfully'
    },
  ]);

  // Subscribe to real-time log events
  useEffect(() => {
    const unsubscribe = subscribeToLogs((log: LogEntry) => {
      addLog(log);
    });

    return unsubscribe;
  }, []);

  // Add log entry with automatic limiting
  const addLog = useCallback((log: LogEntry) => {
    setLogs(prev => {
      const newLogs = [...prev, log];
      // Keep only last MAX_LOGS entries to prevent memory bloat
      return newLogs.length > MAX_LOGS ? newLogs.slice(-MAX_LOGS) : newLogs;
    });
  }, []);

  // Clear all logs
  const clearLogs = useCallback(() => {
    setLogs([{
      timestamp: new Date().toLocaleTimeString(),
      level: 'info',
      message: 'Logs cleared'
    }]);
  }, []);

  return {
    logs,
    addLog,
    clearLogs,
  };
}
