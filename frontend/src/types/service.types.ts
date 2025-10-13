/**
 * Service Types
 * 
 * Defines all type definitions related to service management.
 * Used across the application for type safety.
 */

export interface ServiceModule {
  /** Display name of the service */
  name: string;
  
  /** Actual service name used by systemctl/brew commands */
  serviceName: string;
  
  /** Current running status of the service */
  isRunning: boolean;
  
  /** Loading state during service operations */
  isLoading: boolean;
}

export interface ServiceStatus {
  /** Service name identifier */
  name: string;
  
  /** Running status */
  isRunning: boolean;
}
