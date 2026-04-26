export interface ServiceModule {
  name: string;
  serviceName: string;
  isRunning: boolean;
  isLoading: boolean;
}

export interface ServiceStatus {
  name: string;
  isRunning: boolean;
}

export interface ServiceConfig {
  displayName: string;
  serviceName: string;
}
