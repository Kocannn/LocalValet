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
import { ServiceTable, useServices } from '@/modules/service';
import { useLogs } from '@/modules/log';

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
        <Card className="border-border/70 bg-card/95">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Services</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-semibold tracking-tight">{stats.total}</div>
          </CardContent>
        </Card>
        <Card className="border-border/70 bg-card/95">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Active</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-semibold tracking-tight text-emerald-600">{stats.active}</div>
          </CardContent>
        </Card>
        <Card className="border-border/70 bg-card/95">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Inactive</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-semibold tracking-tight text-slate-600">{stats.inactive}</div>
          </CardContent>
        </Card>
      </div>

      {/* Services Table */}
      <Card className="border-border/70 bg-card/95">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <div>
            <CardTitle className="text-base md:text-lg">Service Management</CardTitle>
            <p className="mt-1 text-sm text-muted-foreground">Start, stop, and review status instantly.</p>
          </div>
          <Badge variant="outline" className="border-emerald-200 bg-emerald-50 text-emerald-700">
            Live
          </Badge>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="py-8 text-center text-sm text-muted-foreground">Loading services...</div>
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
