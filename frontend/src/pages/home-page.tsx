/**
 * Home Page
 * 
 * Dashboard overview page showing service status and recent logs.
 * Provides a high-level view of the application state.
 * 
 * Design Pattern: Page Component (Container)
 * Purpose: Compose features into a complete page view
 */

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { ServiceTable } from '@/features/services';
import { LogViewer, useLogs, type LogEntry } from '@/features/logs';
import { useServices } from '@/features/services';

/**
 * Home dashboard page
 */
export function HomePage() {
  const { services, toggleServiceStatus, isLoading } = useServices();
  const { logs, addLog, clearLogs } = useLogs();

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

  return (
    <div className='flex flex-col gap-4 flex-1 min-h-0 overflow-hidden'>
      {/* Services Section */}
      <div className="w-full">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle>Server Modules</CardTitle>
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

      {/* Logs Section */}
      <div className="w-full flex flex-col flex-1 min-h-0">
        <LogViewer logs={logs} onClear={clearLogs} />
      </div>
    </div>
  );
}
