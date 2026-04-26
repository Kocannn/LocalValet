import type { LogEvents } from '../domain/ports';
import type { LogEntry } from '../domain/types';
import { subscribeToLogs } from '@/services/events.service';

export class WailsLogEvents implements LogEvents {
  subscribe(callback: (log: LogEntry) => void): () => void {
    return subscribeToLogs(callback);
  }
}
