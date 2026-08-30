import { useState } from 'react';
import { FolderPlus, FolderRoot, Trash2 } from 'lucide-react';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

interface ProjectRootsDialogProps {
  open: boolean;
  roots: string[];
  onOpenChange: (open: boolean) => void;
  onAddRoot: (path: string) => Promise<void>;
  onRemoveRoot: (path: string) => Promise<void>;
}

export function ProjectRootsDialog({
  open,
  roots,
  onOpenChange,
  onAddRoot,
  onRemoveRoot,
}: ProjectRootsDialogProps) {
  const [newPath, setNewPath] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleAdd = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newPath.trim()) return;

    try {
      setIsSubmitting(true);
      await onAddRoot(newPath.trim());
      setNewPath('');
    } catch (err) {
      console.error(err);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2 text-base">
            <FolderRoot className="h-4 w-4 text-primary" />
            Project Scan Directories
          </DialogTitle>
          <DialogDescription className="text-xs">
            LocalValet scans these root folders for Laravel, WordPress, Next.js, and other web projects.
          </DialogDescription>
        </DialogHeader>

        {/* Existing Roots List */}
        <div className="space-y-2 py-2">
          <Label className="text-xs text-muted-foreground">Configured Scan Paths</Label>
          <div className="max-h-48 space-y-1.5 overflow-y-auto pr-1">
            {roots.map((root) => (
              <div
                key={root}
                className="flex items-center justify-between gap-2 rounded-md border border-border/70 bg-muted/40 px-3 py-1.5 text-xs font-mono"
              >
                <span className="truncate" title={root}>{root}</span>
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-6 w-6 p-0 text-muted-foreground hover:text-destructive shrink-0"
                  onClick={() => onRemoveRoot(root)}
                >
                  <Trash2 className="h-3.5 w-3.5" />
                </Button>
              </div>
            ))}
            {roots.length === 0 && (
              <p className="py-3 text-center text-xs text-muted-foreground">
                No custom scan directories configured. Default fallback is used.
              </p>
            )}
          </div>
        </div>

        {/* Add New Root Form */}
        <form onSubmit={handleAdd} className="space-y-2 pt-2 border-t border-border/50">
          <Label htmlFor="root-path" className="text-xs font-medium">Add Directory Path</Label>
          <div className="flex gap-2">
            <Input
              id="root-path"
              placeholder="/home/username/Projects"
              value={newPath}
              onChange={(e) => setNewPath(e.target.value)}
              className="text-xs font-mono h-8"
              disabled={isSubmitting}
            />
            <Button type="submit" size="sm" className="h-8 gap-1 text-xs" disabled={isSubmitting || !newPath.trim()}>
              <FolderPlus className="h-3.5 w-3.5" />
              Add
            </Button>
          </div>
        </form>

        <DialogFooter className="pt-2">
          <Button variant="outline" size="sm" onClick={() => onOpenChange(false)}>
            Close
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
