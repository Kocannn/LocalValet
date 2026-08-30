import { useCallback, useEffect, useMemo, useState } from 'react';
import type { RuntimeService } from '../domain/types';
import { WailsRuntimeSettingsRepository } from '../infrastructure/wails-runtime-settings.repository';

interface UseRuntimeSettingsReturn {
  services: RuntimeService[];
  isLoading: boolean;
  savingMap: Record<string, boolean>;
  errorMap: Record<string, string>;
  changeVersion: (serviceName: string, version: string) => Promise<void>;
  refresh: () => Promise<void>;
}

export function useRuntimeSettings(): UseRuntimeSettingsReturn {
  const repository = useMemo(() => new WailsRuntimeSettingsRepository(), []);

  const [services, setServices] = useState<RuntimeService[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [savingMap, setSavingMap] = useState<Record<string, boolean>>({});
  const [errorMap, setErrorMap] = useState<Record<string, string>>({});

  const refresh = useCallback(async () => {
    setIsLoading(true);
    try {
      const data = await repository.getAllRuntimeServices();
      setServices(data);
    } catch (err) {
      console.error('Failed to load runtime services:', err);
    } finally {
      setIsLoading(false);
    }
  }, [repository]);

  useEffect(() => {
    refresh();
  }, [refresh]);

  const changeVersion = useCallback(
    async (serviceName: string, version: string) => {
      if (!serviceName || !version || savingMap[serviceName]) {
        return;
      }

      setSavingMap((prev) => ({ ...prev, [serviceName]: true }));
      setErrorMap((prev) => ({ ...prev, [serviceName]: '' }));

      try {
        await repository.setServiceVersion(serviceName, version);
        // Optimistically update
        setServices((prev) =>
          prev.map((svc) =>
            svc.serviceName === serviceName ? { ...svc, activeVersion: version } : svc,
          ),
        );
      } catch (err: any) {
        console.error(`Failed to change version for ${serviceName}:`, err);
        setErrorMap((prev) => ({
          ...prev,
          [serviceName]: err?.message || 'Failed to switch version',
        }));
      } finally {
        setSavingMap((prev) => ({ ...prev, [serviceName]: false }));
      }
    },
    [repository, savingMap],
  );

  return {
    services,
    isLoading,
    savingMap,
    errorMap,
    changeVersion,
    refresh,
  };
}
