import * as WailsApp from "../../wailsjs/go/main/App.js";
import type { ServiceConfig, ServiceStatus } from "@/types";

export async function getAllServicesStatus(): Promise<ServiceStatus[]> {
  try {
    return await (WailsApp as any).GetAllServicesStatus();
  } catch (error) {
    console.error('Failed to get services status:', error);
    throw error;
  }
}

export async function getAvailableServices(): Promise<ServiceConfig[]> {
  try {
    const services = await (WailsApp as any).GetAvailableServices();
    return (services ?? []).map((service: any) => ({
      displayName: service.DisplayName,
      serviceName: service.ServiceName,
    }));
  } catch (error) {
    console.error('Failed to get available services:', error);
    throw error;
  }
}

export async function toggleService(serviceName: string, enable: boolean): Promise<void> {
  try {
    await (WailsApp as any).ToggleService(serviceName, enable);
  } catch (error) {
    console.error(`Failed to toggle service ${serviceName}:`, error);
    throw error;
  }
}
