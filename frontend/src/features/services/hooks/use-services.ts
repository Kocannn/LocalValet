/**
 * useServices Hook
 * 
 * Custom hook for managing service state and operations.
 * Encapsulates all service-related business logic.
 * 
 * Design Pattern: Custom Hook Pattern
 * Purpose: Separation of concerns - business logic vs UI
 */

import { useState, useEffect, useCallback } from 'react';
import { getAllServicesStatus, toggleService } from '@/services';
import { subscribeToServiceStatus } from '@/services/events.service';
import type { ServiceModule, ServiceStatus } from '../types';

// Default services configuration
const DEFAULT_SERVICES: ServiceModule[] = [
  { name: 'Apache', serviceName: 'apache2', isRunning: false, isLoading: false },
  { name: 'MySQL', serviceName: 'mysql', isRunning: false, isLoading: false },
  { name: 'PostgreSQL', serviceName: 'postgresql', isRunning: false, isLoading: false },
  { name: 'Redis', serviceName: 'redis-server', isRunning: false, isLoading: false },
  { name: 'Nginx', serviceName: 'nginx', isRunning: false, isLoading: false },
  { name: 'PHP-FPM', serviceName: 'php8.1-fpm', isRunning: false, isLoading: false },
];

interface UseServicesReturn {
  services: ServiceModule[];
  toggleServiceStatus: (service: ServiceModule, checked: boolean) => Promise<void>;
  isLoading: boolean;
}

/**
 * Hook for service management
 * @returns Service state and control functions
 */
export function useServices(): UseServicesReturn {
  const [services, setServices] = useState<ServiceModule[]>(DEFAULT_SERVICES);
  const [isLoading, setIsLoading] = useState(true);

  // Load initial service status
  useEffect(() => {
    const loadInitialStatus = async () => {
      try {
        setIsLoading(true);
        const statuses = await getAllServicesStatus();
        
        setServices(prev => prev.map(service => {
          const status = statuses.find((s: ServiceStatus) => s.name === service.serviceName);
          return status ? { ...service, isRunning: status.isRunning } : service;
        }));
      } catch (error) {
        console.error('Failed to load initial service status:', error);
      } finally {
        setIsLoading(false);
      }
    };

    loadInitialStatus();
  }, []);

  // Subscribe to real-time status updates
  useEffect(() => {
    const unsubscribe = subscribeToServiceStatus((status: ServiceStatus) => {
      setServices(prev => prev.map(service =>
        service.serviceName === status.name
          ? { ...service, isRunning: status.isRunning, isLoading: false }
          : service
      ));
    });

    return unsubscribe;
  }, []);

  // Toggle service function
  const toggleServiceStatus = useCallback(async (service: ServiceModule, checked: boolean) => {
    // Set loading state
    setServices(prev => prev.map(s =>
      s.serviceName === service.serviceName
        ? { ...s, isLoading: true }
        : s
    ));

    try {
      await toggleService(service.serviceName, checked);
      // Status update will come via event
    } catch (error: any) {
      console.error(`Failed to toggle ${service.name}:`, error);
      
      // Reset loading state on error
      setServices(prev => prev.map(s =>
        s.serviceName === service.serviceName
          ? { ...s, isLoading: false }
          : s
      ));
      
      throw error;
    }
  }, []);

  return {
    services,
    toggleServiceStatus,
    isLoading,
  };
}
