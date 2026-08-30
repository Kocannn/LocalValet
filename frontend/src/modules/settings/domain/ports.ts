import type { RuntimeService } from './types';

export interface RuntimeSettingsRepository {
  getAllRuntimeServices(): Promise<RuntimeService[]>;
  getServiceVersions(serviceName: string): Promise<string[]>;
  getActiveServiceVersion(serviceName: string): Promise<string>;
  setServiceVersion(serviceName: string, version: string): Promise<void>;
}
