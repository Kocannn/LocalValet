/**
 * ServiceTable Component
 * 
 * Displays a table of services with their status and controls.
 * Pure presentational component with callback props.
 * 
 * Design Pattern: Presentational Component
 * Purpose: UI-only, no business logic
 */

import { memo } from 'react';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Switch } from '@/components/ui/switch';
import { Loader2 } from 'lucide-react';
import type { ServiceModule } from '../types';

interface ServiceTableProps {
  /** List of services to display */
  services: ServiceModule[];
  
  /** Callback when service switch is toggled */
  onServiceToggle: (service: ServiceModule, checked: boolean) => void;
}

/**
 * Service management table component
 */
export const ServiceTable = memo(({ services, onServiceToggle }: ServiceTableProps) => {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className="text-center">Module</TableHead>
          <TableHead className='text-center'>Status</TableHead>
          <TableHead className='text-center'>Action</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {services.map((service) => (
          <TableRow key={service.serviceName}>
            <TableCell className="font-medium">{service.name}</TableCell>
            <TableCell>
              <div className="flex items-center justify-center">
                {service.isLoading ? (
                  <Badge variant="secondary" className="flex items-center gap-1">
                    <Loader2 className="h-3 w-3 animate-spin" />
                    Loading...
                  </Badge>
                ) : (
                  <Badge variant={service.isRunning ? "default" : "secondary"}>
                    {service.isRunning ? "Active" : "Inactive"}
                  </Badge>
                )}
              </div>
            </TableCell>
            <TableCell>
              <div className='flex space-x-2 justify-center'>
                <Switch
                  id={`service-${service.serviceName}`}
                  checked={service.isRunning}
                  disabled={service.isLoading}
                  onCheckedChange={(checked) => onServiceToggle(service, checked)}
                />
              </div>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
});

ServiceTable.displayName = 'ServiceTable';
