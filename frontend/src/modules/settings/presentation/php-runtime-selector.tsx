import { ChevronDown } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';

interface PhpRuntimeSelectorProps {
  versions: string[];
  label: string;
  isLoading: boolean;
  isSaving: boolean;
  error: string;
  onChange: (version: string) => Promise<void>;
}

export function PhpRuntimeSelector({
  versions,
  label,
  isLoading,
  isSaving,
  error,
  onChange,
}: PhpRuntimeSelectorProps) {
  return (
    <div className="space-y-2">
      <Label className="text-sm">PHP Runtime Version</Label>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button
            variant="outline"
            className="w-full justify-between"
            disabled={isSaving || isLoading}
          >
            {label}
            <ChevronDown className="h-4 w-4" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent className="w-[var(--radix-dropdown-menu-trigger-width)]" align="start">
          {versions.length === 0 ? (
            <DropdownMenuItem disabled>Tidak ada versi tersedia</DropdownMenuItem>
          ) : (
            versions.map((version) => (
              <DropdownMenuItem
                key={version}
                onSelect={() => onChange(version)}
              >
                PHP {version}
              </DropdownMenuItem>
            ))
          )}
        </DropdownMenuContent>
      </DropdownMenu>
      {isLoading && (
        <p className="text-xs text-muted-foreground">Memuat daftar versi PHP...</p>
      )}
      {error && (
        <p className="text-xs text-amber-600">{error}</p>
      )}
      <p className="text-xs text-muted-foreground">
        Mengubah versi aktif akan dipakai saat service PHP-FPM berikutnya dijalankan.
      </p>
    </div>
  );
}
