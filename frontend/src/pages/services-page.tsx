/**
 * Services Page
 * 
 * Dedicated page for detailed service management.
 * Shows comprehensive service controls and status.
 * 
 * Design Pattern: Page Component (Container)
 * Purpose: Feature-focused page for service management
 */

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { ServiceTable, useServices } from '@/features/services';
import { useLogs } from '@/features/logs';

/**
 * Services management page
 */
export function ServicesPage() {
  const { services, toggleServiceStatus, isLoading } = useServices();
  const { addLog } = useLogs();

  // Handle service toggle with logging
  const handleServiceToggle = async (service: any, checked: boolean) => {
    // Add info log
    addLog({
      timestamp: new Date().toLocaleTimeString(),
      level: 'info',
      message: `${checked ? 'Starting' : 'Stopping'} ${service.name}...`
    });

    try {
      await toggleServiceStatus(service, checked);
    } catch (error: any) {
      // Add error log
      addLog({
        timestamp: new Date().toLocaleTimeString(),
        level: 'error',
        message: `Failed to ${checked ? 'start' : 'stop'} ${service.name}: ${error.message || error}`
      });
    }
  };

  // Calculate statistics
  const stats = {
    total: services.length,
    active: services.filter(s => s.isRunning).length,
    inactive: services.filter(s => !s.isRunning).length,
  };

  return (
    <div className='flex flex-col gap-4'>
      {/* Statistics Cards */}
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Services</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.total}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Active</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">{stats.active}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Inactive</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-gray-600">{stats.inactive}</div>
          </CardContent>
        </Card>
      </div>

      {/* Services Table */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle>Service Management</CardTitle>
          <Badge variant="outline" className="flex items-center gap-1 bg-green-50 text-green-700 border-green-200">
            <span className="relative flex h-2 w-2">
              <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
              <span className="relative inline-flex rounded-full h-2 w-2 bg-green-500"></span>
            </span>
            Real-time
          </Badge>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="text-center py-8">Loading services...</div>
          ) : (
            <ServiceTable 
              services={services}
              onServiceToggle={handleServiceToggle}
            />
          )}
        </CardContent>
      </Card>
    </div>
  );
}
