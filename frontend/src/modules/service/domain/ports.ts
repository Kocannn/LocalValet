import type { ServiceConfig, ServiceStatus } from './types';

export interface ServiceRepository {
  getAllServicesStatus(): Promise<ServiceStatus[]>;
  getAvailableServices(): Promise<ServiceConfig[]>;
  toggleService(serviceName: string, enable: boolean): Promise<void>;
}

export interface ServiceStatusEvents {
  subscribe(callback: (status: ServiceStatus) => void): () => void;
}
