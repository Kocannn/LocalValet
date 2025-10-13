/**
 * Root Layout
 * 
 * Main layout wrapper for the application.
 * Provides consistent sidebar navigation and content area.
 * 
 * Design Pattern: Layout Component
 * Purpose: Consistent UI structure across all pages
 */

import { Outlet } from 'react-router-dom';
import { AppLayout } from '@/components/app-layout';

/**
 * Root layout with sidebar and main content area
 */
export function RootLayout() {
  return (
    <AppLayout title="LocalValet">
      <Outlet />
    </AppLayout>
  );
}
