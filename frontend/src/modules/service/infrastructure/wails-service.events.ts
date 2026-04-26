import type { ServiceStatusEvents } from '../domain/ports';
import type { ServiceStatus } from '../domain/types';
import { subscribeToServiceStatus } from '@/services/events.service';

export class WailsServiceStatusEvents implements ServiceStatusEvents {
  subscribe(callback: (status: ServiceStatus) => void): () => void {
    return subscribeToServiceStatus(callback);
  }
}
