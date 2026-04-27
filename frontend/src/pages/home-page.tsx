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
import { Button } from '@/components/ui/button';
import { ServiceTable, useServices } from '@/modules/service';
import { LogViewer, useLogs } from '@/modules/log';
import { openContextTerminal } from '@/services';

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

  const handleOpenContextTerminal = async () => {
    try {
      await openContextTerminal('');
    } catch (error) {
      console.error('Failed to open context terminal:', error);
    }
  };

  return (
    <div className='flex flex-1 min-h-0 overflow-hidden'>
      <div className="grid w-full gap-4 min-h-0 lg:grid-cols-[minmax(0,1.25fr)_minmax(320px,0.75fr)]">
        <div className="flex flex-col gap-4 min-h-0">
          {/* Services Section */}
          <div className="w-full">
            <Card className="border-border/70 bg-card/95">
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <div>
                  <CardTitle className="text-base md:text-lg">Server Modules</CardTitle>
                  <p className="mt-1 text-sm text-muted-foreground">Monitor and control your local services.</p>
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

          {/* Context Terminal Section */}
          <div className="w-full">
            <Card className="border-border/70 bg-card/95">
              <CardHeader className="pb-2">
                <CardTitle className="text-base md:text-lg">Context Terminal</CardTitle>
                <p className="mt-1 text-sm text-muted-foreground">
                  Open a terminal with the LocalValet runtime environment.
                </p>
              </CardHeader>
              <CardContent>
                <Button className="w-full sm:w-auto" onClick={handleOpenContextTerminal}>
                  Open Context Terminal
                </Button>
              </CardContent>
            </Card>
          </div>
        </div>

        {/* Logs Section */}
        <div className="w-full min-h-0 lg:sticky lg:top-0 lg:self-start">
          <LogViewer logs={logs} onClear={clearLogs} />
        </div>
      </div>
    </div>
  );
}
