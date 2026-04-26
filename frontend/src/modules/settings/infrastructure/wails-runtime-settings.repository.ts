import type { RuntimeSettingsRepository } from '../domain/ports';
import {
  getActiveServiceVersion,
  getServiceVersions,
  setServiceVersion,
} from '@/services/wails-runtime-settings.service';

export class WailsRuntimeSettingsRepository implements RuntimeSettingsRepository {
  getServiceVersions(serviceName: string): Promise<string[]> {
    return getServiceVersions(serviceName);
  }

  getActiveServiceVersion(serviceName: string): Promise<string> {
    return getActiveServiceVersion(serviceName);
  }

  setServiceVersion(serviceName: string, version: string): Promise<void> {
    return setServiceVersion(serviceName, version);
  }
}
