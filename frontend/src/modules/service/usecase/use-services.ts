import { useCallback, useEffect, useMemo, useState } from 'react';
import type { ServiceModule, ServiceStatus } from '../domain/types';
import { WailsServiceRepository } from '../infrastructure/wails-service.repository';
import { WailsServiceStatusEvents } from '../infrastructure/wails-service.events';

const DEFAULT_SERVICES: ServiceModule[] = [];

interface UseServicesReturn {
  services: ServiceModule[];
  toggleServiceStatus: (service: ServiceModule, checked: boolean) => Promise<void>;
  isLoading: boolean;
}

export function useServices(): UseServicesReturn {
  const repository = useMemo(() => new WailsServiceRepository(), []);
  const events = useMemo(() => new WailsServiceStatusEvents(), []);

  const [services, setServices] = useState<ServiceModule[]>(DEFAULT_SERVICES);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const loadInitialStatus = async () => {
      try {
        setIsLoading(true);
        const [configs, statuses] = await Promise.all([
          repository.getAvailableServices(),
          repository.getAllServicesStatus(),
        ]);

        const statusMap = new Map<string, ServiceStatus>(
          statuses.map((status) => [status.name, status]),
        );

        setServices(
          configs.map((config) => ({
            name: config.displayName,
            serviceName: config.serviceName,
            isRunning: statusMap.get(config.serviceName)?.isRunning ?? false,
            isLoading: false,
          })),
        );
      } catch (error) {
        console.error('Failed to load initial service status:', error);
      } finally {
        setIsLoading(false);
      }
    };

    loadInitialStatus();
  }, [repository]);

  useEffect(() => {
    return events.subscribe((status: ServiceStatus) => {
      setServices((prev) =>
        prev.map((service) =>
          service.serviceName === status.name
            ? { ...service, isRunning: status.isRunning, isLoading: false }
            : service,
        ),
      );
    });
  }, [events]);

  const toggleServiceStatus = useCallback(
    async (service: ServiceModule, checked: boolean) => {
      setServices((prev) =>
        prev.map((s) =>
          s.serviceName === service.serviceName
            ? { ...s, isLoading: true }
            : s,
        ),
      );

      try {
        await repository.toggleService(service.serviceName, checked);
      } catch (error: any) {
        console.error(`Failed to toggle ${service.name}:`, error);
        setServices((prev) =>
          prev.map((s) =>
            s.serviceName === service.serviceName
              ? { ...s, isLoading: false }
              : s,
          ),
        );
        throw error;
      }
    },
    [repository],
  );

  return {
    services,
    toggleServiceStatus,
    isLoading,
  };
}
