import type { ServiceRepository } from '../domain/ports';
import type { ServiceConfig, ServiceStatus } from '../domain/types';
import {
  getAllServicesStatus as getAllStatuses,
  getAvailableServices as getServices,
  toggleService as toggle,
} from '@/services/wails-service-control.service';

export class WailsServiceRepository implements ServiceRepository {
  async getAllServicesStatus(): Promise<ServiceStatus[]> {
    return getAllStatuses();
  }

  async getAvailableServices(): Promise<ServiceConfig[]> {
    return getServices();
  }

  async toggleService(serviceName: string, enable: boolean): Promise<void> {
    return toggle(serviceName, enable);
  }
}
