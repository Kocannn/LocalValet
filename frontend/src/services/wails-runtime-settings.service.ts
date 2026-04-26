import * as WailsApp from "../../wailsjs/go/main/App.js";

export async function getServiceVersions(serviceName: string): Promise<string[]> {
  try {
    return await (WailsApp as any).GetServiceVersions(serviceName);
  } catch (error) {
    console.error(`Failed to get versions for ${serviceName}:`, error);
    throw error;
  }
}

export async function getActiveServiceVersion(serviceName: string): Promise<string> {
  try {
    return await (WailsApp as any).GetActiveServiceVersion(serviceName);
  } catch (error) {
    console.error(`Failed to get active version for ${serviceName}:`, error);
    throw error;
  }
}

export async function setServiceVersion(serviceName: string, version: string): Promise<void> {
  try {
    await (WailsApp as any).SetServiceVersion(serviceName, version);
  } catch (error) {
    console.error(`Failed to set version for ${serviceName}:`, error);
    throw error;
  }
}
