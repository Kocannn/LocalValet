/**
 * Services Feature - Public API
 * 
 * This is the ONLY file that should be imported from outside this feature.
 * Provides a clean public interface while hiding internal implementation.
 * 
 * Import pattern: import { ServiceTable, useServices } from '@/features/services'
 */

export { ServiceTable } from './components';
export { useServices } from './hooks';
export type { ServiceModule, ServiceStatus } from './types';
