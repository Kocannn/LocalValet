import { memo } from 'react';
import { Loader2 } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Switch } from '@/components/ui/switch';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import type { ServiceModule } from '../domain/types';

interface ServiceTableProps {
  services: ServiceModule[];
  onServiceToggle: (service: ServiceModule, checked: boolean) => void;
}

export const ServiceTable = memo(({ services, onServiceToggle }: ServiceTableProps) => {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Module</TableHead>
          <TableHead className='text-center'>Status</TableHead>
          <TableHead className='text-right'>Action</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {services.map((service) => (
          <TableRow key={service.serviceName}>
            <TableCell className="font-medium tracking-tight">{service.name}</TableCell>
            <TableCell>
              <div className="flex items-center justify-center">
                {service.isLoading ? (
                  <Badge variant="secondary" className="flex items-center gap-1 rounded-md">
                    <Loader2 className="h-3 w-3 animate-spin" />
                    Loading...
                  </Badge>
                ) : (
                  <Badge variant={service.isRunning ? "default" : "secondary"} className="rounded-md">
                    {service.isRunning ? "Active" : "Inactive"}
                  </Badge>
                )}
              </div>
            </TableCell>
            <TableCell>
              <div className='flex justify-end'>
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
