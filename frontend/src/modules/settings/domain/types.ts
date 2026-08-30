export interface RuntimeService {
  serviceName: string;
  displayName: string;
  activeVersion: string;
  availableVersions: string[];
  category: string;
  isRunning: boolean;
}

export interface RuntimeVersionSetting {
  serviceName: string;
  activeVersion: string;
  availableVersions: string[];
}
