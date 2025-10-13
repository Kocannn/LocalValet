import { useState, useEffect, useRef } from 'react';
import * as WailsApp from "../wailsjs/go/main/App.js";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card.js"
import { Button } from './components/ui/button.js';
import { AppLayout } from "./components/app-layout";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from './components/ui/table.js';
import { Trash2, Loader2 } from 'lucide-react';
import { Badge } from './components/ui/badge.js';
import { Switch } from './components/ui/switch.js';

interface LogEntry {
  timestamp: string;
  level: 'info' | 'warning' | 'error' | 'success';
  message: string;
}

interface ServiceModule {
  name: string;
  serviceName: string; // actual service name for systemctl/brew
  isRunning: boolean;
  isLoading: boolean;
}

function App() {
  const [modules, setModules] = useState<ServiceModule[]>([
    { name: 'Apache', serviceName: 'apache2', isRunning: false, isLoading: false },
    { name: 'MySQL', serviceName: 'mysql', isRunning: false, isLoading: false },
    { name: 'PostgreSQL', serviceName: 'postgresql', isRunning: false, isLoading: false },
    { name: 'Redis', serviceName: 'redis-server', isRunning: false, isLoading: false },
    { name: 'Nginx', serviceName: 'nginx', isRunning: false, isLoading: false },
    { name: 'PHP-FPM', serviceName: 'php8.1-fpm', isRunning: false, isLoading: false },
  ]);
  const [logs, setLogs] = useState<LogEntry[]>([
    {
      timestamp: new Date().toLocaleTimeString(),
      level: 'info',
      message: 'LocalValet initialized successfully'
    },
  ]);

  const logsEndRef = useRef<HTMLDivElement>(null);


  // Add log entry
  const addLog = (log: LogEntry) => {
    setLogs(prev => [...prev, log]);
  };

  // Check service status on mount and periodically
  useEffect(() => {
    const checkAllServices = async () => {
      for (const module of modules) {
        try {
          const status = await (WailsApp as any).GetServiceStatus(module.serviceName);
          setModules(prev => prev.map(m =>
            m.serviceName === module.serviceName
              ? { ...m, isRunning: status.isRunning }
              : m
          ));
        } catch (error) {
          console.error(`Failed to check ${module.serviceName} status:`, error);
        }
      }
    };

    checkAllServices();
    const interval = setInterval(checkAllServices, 5000); // Check every 5 seconds

    return () => clearInterval(interval);
  }, []);

  // Handle service toggle
  const handleServiceToggle = async (module: ServiceModule, checked: boolean) => {
    // Set loading state
    setModules(prev => prev.map(m =>
      m.serviceName === module.serviceName
        ? { ...m, isLoading: true }
        : m
    ));

    // Add info log
    addLog({
      timestamp: new Date().toLocaleTimeString(),
      level: 'info',
      message: `${checked ? 'Starting' : 'Stopping'} ${module.name}...`
    });

    try {
      const result = await (WailsApp as any).ToggleService(module.serviceName, checked);

      // Add result log
      addLog({
        timestamp: result.timestamp || new Date().toLocaleTimeString(),
        level: result.level,
        message: result.message
      });

      // Update module status
      setModules(prev => prev.map(m =>
        m.serviceName === module.serviceName
          ? { ...m, isRunning: checked, isLoading: false }
          : m
      ));
    } catch (error: any) {
      // Add error log
      addLog({
        timestamp: new Date().toLocaleTimeString(),
        level: 'error',
        message: `Failed to ${checked ? 'start' : 'stop'} ${module.name}: ${error.message || error}`
      });

      // Reset loading state
      setModules(prev => prev.map(m =>
        m.serviceName === module.serviceName
          ? { ...m, isLoading: false }
          : m
      ));
    }
  };

  // Auto-scroll to bottom when new log is added
  useEffect(() => {
    logsEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [logs]);



  const clearLogs = () => {
    setLogs([]);
  };

  const getLogColor = (level: string) => {
    switch (level) {
      case 'info':
        return 'text-blue-400';
      case 'success':
        return 'text-green-400';
      case 'warning':
        return 'text-yellow-400';
      case 'error':
        return 'text-red-400';
      default:
        return 'text-gray-400';
    }
  };

  return (
    <AppLayout title="LocalValet">
      <div className="w-full">
        <Card>
          <CardHeader>
            <CardTitle>Server Modules</CardTitle>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="text-center">Module</TableHead>
                  <TableHead className='text-center'>Status</TableHead>
                  <TableHead className='text-center'>Action</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {modules.map((module) => (
                  <TableRow key={module.serviceName}>
                    <TableCell className="font-medium">{module.name}</TableCell>
                    <TableCell>
                      <div className="flex items-center justify-center">
                        {module.isLoading ? (
                          <Badge variant="secondary" className="flex items-center gap-1">
                            <Loader2 className="h-3 w-3 animate-spin" />
                            Loading...
                          </Badge>
                        ) : (
                          <Badge variant={module.isRunning ? "default" : "secondary"}>
                            {module.isRunning ? "Active" : "Inactive"}
                          </Badge>
                        )}
                      </div>
                    </TableCell>
                    <TableCell>
                      <div className='flex space-x-2 justify-center'>
                        <Switch
                          id={`module-${module.serviceName}`}
                          checked={module.isRunning}
                          disabled={module.isLoading}
                          onCheckedChange={(checked) => handleServiceToggle(module, checked)}
                        />
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </div>

      {/* Server Logs */}
      <div className="w-full">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle>Server Logs</CardTitle>
            <Button
              variant="ghost"
              size="sm"
              onClick={clearLogs}
              className="h-8 w-8 p-0"
            >
              <Trash2 className="h-4 w-4" />
            </Button>
          </CardHeader>
          <CardContent>
            <div className="relative">
              <div className="bg-[#FAFAFA] rounded-lg p-4 h-[400px] overflow-y-auto font-mono text-sm">
                {logs.length === 0 ? (
                  <div className="text-black text-center py-8">
                    No logs available. Waiting for server activity...
                  </div>
                ) : (
                  <div className="space-y-1">
                    {logs.map((log, index) => (
                      <div key={index} className="flex gap-2">
                        <span className="text-black">[{log.timestamp}]</span>
                        <span className={`font-semibold ${getLogColor(log.level)}`}>
                          [{log.level.toUpperCase()}]
                        </span>
                        <span className="text-black">{log.message}</span>
                      </div>
                    ))}
                    <div ref={logsEndRef} />
                  </div>
                )}
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </AppLayout>
  )
} export default App
