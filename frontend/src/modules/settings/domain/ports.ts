export interface RuntimeSettingsRepository {
  getServiceVersions(serviceName: string): Promise<string[]>;
  getActiveServiceVersion(serviceName: string): Promise<string>;
  setServiceVersion(serviceName: string, version: string): Promise<void>;
}
