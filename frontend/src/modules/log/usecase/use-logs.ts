import { useCallback, useEffect, useMemo, useState } from 'react';
import type { LogEntry } from '../domain/types';
import { WailsLogEvents } from '../infrastructure/wails-log.events';

const MAX_LOGS = 50;

interface UseLogsReturn {
  logs: LogEntry[];
  addLog: (log: LogEntry) => void;
  clearLogs: () => void;
}

export function useLogs(): UseLogsReturn {
  const events = useMemo(() => new WailsLogEvents(), []);
  const [logs, setLogs] = useState<LogEntry[]>([
    {
      timestamp: new Date().toLocaleTimeString(),
      level: 'info',
      message: 'LocalValet initialized successfully',
    },
  ]);

  const addLog = useCallback((log: LogEntry) => {
    setLogs((prev) => {
      const newLogs = [...prev, log];
      return newLogs.length > MAX_LOGS ? newLogs.slice(-MAX_LOGS) : newLogs;
    });
  }, []);

  const clearLogs = useCallback(() => {
    setLogs([
      {
        timestamp: new Date().toLocaleTimeString(),
        level: 'info',
        message: 'Logs cleared',
      },
    ]);
  }, []);

  useEffect(() => {
    return events.subscribe((log: LogEntry) => {
      addLog(log);
    });
  }, [addLog, events]);

  return {
    logs,
    addLog,
    clearLogs,
  };
}
