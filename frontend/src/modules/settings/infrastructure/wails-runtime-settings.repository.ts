import type { RuntimeSettingsRepository } from '../domain/ports';
import type { RuntimeService } from '../domain/types';
import * as WailsApp from '../../../../wailsjs/go/main/App.js';

export class WailsRuntimeSettingsRepository implements RuntimeSettingsRepository {
  async getAllRuntimeServices(): Promise<RuntimeService[]> {
    try {
      const services = await (WailsApp as any).GetAllRuntimeServices();
      return (services ?? []).map((s: any) => ({
        serviceName: s.serviceName ?? s.ServiceName ?? '',
        displayName: s.displayName ?? s.DisplayName ?? s.serviceName ?? '',
        activeVersion: s.activeVersion ?? s.ActiveVersion ?? '',
        availableVersions: s.availableVersions ?? s.AvailableVersions ?? [],
        category: s.category ?? s.Category ?? 'Runtime',
        isRunning: s.isRunning ?? s.IsRunning ?? false,
      }));
    } catch (error) {
      console.error('Failed to get all runtime services:', error);
      return [];
    }
  }

  async getServiceVersions(serviceName: string): Promise<string[]> {
    try {
      return await (WailsApp as any).GetServiceVersions(serviceName);
    } catch (error) {
      console.error(`Failed to get versions for ${serviceName}:`, error);
      return [];
    }
  }

  async getActiveServiceVersion(serviceName: string): Promise<string> {
    try {
      return await (WailsApp as any).GetActiveServiceVersion(serviceName);
    } catch (error) {
      console.error(`Failed to get active version for ${serviceName}:`, error);
      return '';
    }
  }

  async setServiceVersion(serviceName: string, version: string): Promise<void> {
    try {
      await (WailsApp as any).SetServiceVersion(serviceName, version);
    } catch (error) {
      console.error(`Failed to set active version for ${serviceName}:`, error);
      throw error;
    }
  }
}
