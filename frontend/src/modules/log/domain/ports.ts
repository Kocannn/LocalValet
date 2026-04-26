import type { LogEntry } from './types';

export interface LogEvents {
  subscribe(callback: (log: LogEntry) => void): () => void;
}
