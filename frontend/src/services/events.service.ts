/**
 * Events Service
 * 
 * Abstraction layer for Wails runtime events.
 * Provides type-safe event subscriptions and cleanup.
 * 
 * Design Pattern: Observer Pattern
 * Purpose: Centralize event management and prevent memory leaks
 */

import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime';
import type { ServiceStatus, LogEntry } from '@/types';

export type ServiceStatusCallback = (status: ServiceStatus) => void;
export type LogCallback = (log: LogEntry) => void;

/**
 * Subscribe to service status changes
 * @param callback - Function to call when status changes
 * @returns Cleanup function
 */
export function subscribeToServiceStatus(callback: ServiceStatusCallback): () => void {
  EventsOn('service:status-changed', callback);
  
  return () => {
    EventsOff('service:status-changed');
  };
}

/**
 * Subscribe to log events
 * @param callback - Function to call when new log arrives
 * @returns Cleanup function
 */
export function subscribeToLogs(callback: LogCallback): () => void {
  EventsOn('service:log', callback);
  
  return () => {
    EventsOff('service:log');
  };
}
