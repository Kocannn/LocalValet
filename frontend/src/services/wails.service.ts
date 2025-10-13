/**
 * Wails Service
 * 
 * Abstraction layer for Wails Go functions.
 * Provides type-safe access to backend functionality.
 * 
 * Design Pattern: Service Layer Pattern
 * Purpose: Decouple business logic from UI components
 */

import * as WailsApp from "../../wailsjs/go/main/App.js";
import type { ServiceStatus } from "@/types";

/**
 * Get status of all services
 * @returns Promise with array of service statuses
 */
export async function getAllServicesStatus(): Promise<ServiceStatus[]> {
  try {
    return await (WailsApp as any).GetAllServicesStatus();
  } catch (error) {
    console.error('Failed to get services status:', error);
    throw error;
  }
}

/**
 * Toggle service on/off
 * @param serviceName - Name of the service to toggle
 * @param enable - True to start, false to stop
 */
export async function toggleService(serviceName: string, enable: boolean): Promise<void> {
  try {
    await (WailsApp as any).ToggleService(serviceName, enable);
  } catch (error) {
    console.error(`Failed to toggle service ${serviceName}:`, error);
    throw error;
  }
}
