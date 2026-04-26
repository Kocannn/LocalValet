import { useCallback, useEffect, useMemo, useState } from 'react';
import { WailsRuntimeSettingsRepository } from '../infrastructure/wails-runtime-settings.repository';

const PHP_FPM_SERVICE = 'php-fpm';

interface UsePhpRuntimeSettingsReturn {
  versions: string[];
  activeVersion: string;
  label: string;
  isLoading: boolean;
  isSaving: boolean;
  error: string;
  changeVersion: (version: string) => Promise<void>;
}

function mergeVersions(versions: string[], active: string): string[] {
  const merged = [...versions, active]
    .map((version) => version.trim())
    .filter((version) => version.length > 0);

  return Array.from(new Set(merged));
}

export function usePhpRuntimeSettings(): UsePhpRuntimeSettingsReturn {
  const repository = useMemo(() => new WailsRuntimeSettingsRepository(), []);

  const [versions, setVersions] = useState<string[]>([]);
  const [activeVersion, setActiveVersion] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    const loadPhpVersions = async () => {
      setIsLoading(true);
      setError('');

      try {
        const [availableVersions, active] = await Promise.all([
          repository.getServiceVersions(PHP_FPM_SERVICE),
          repository.getActiveServiceVersion(PHP_FPM_SERVICE),
        ]);

        const mergedVersions = mergeVersions(availableVersions ?? [], active ?? '');
        setVersions(mergedVersions);
        setActiveVersion(active || mergedVersions[0] || '');
      } catch (loadError) {
        console.error('Failed to load PHP runtime versions:', loadError);
        setError('Gagal mengambil versi dari backend.');
        setVersions([]);
        setActiveVersion('');
      } finally {
        setIsLoading(false);
      }
    };

    loadPhpVersions();
  }, [repository]);

  const changeVersion = useCallback(
    async (version: string) => {
      if (!version || version === activeVersion || isLoading || isSaving) {
        return;
      }

      try {
        setIsSaving(true);
        await repository.setServiceVersion(PHP_FPM_SERVICE, version);
        setActiveVersion(version);
      } catch (saveError) {
        console.error('Failed to update PHP version:', saveError);
      } finally {
        setIsSaving(false);
      }
    },
    [activeVersion, isLoading, isSaving, repository],
  );

  return {
    versions,
    activeVersion,
    label: activeVersion ? `PHP ${activeVersion}` : 'No version selected',
    isLoading,
    isSaving,
    error,
    changeVersion,
  };
}
