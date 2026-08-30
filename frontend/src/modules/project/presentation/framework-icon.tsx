import {
  Code2,
  FileCode,
  Globe,
  Layers,
  Sparkles,
  Zap,
} from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import type { FrameworkType } from '../domain/types';

interface FrameworkBadgeProps {
  framework: FrameworkType | string;
  className?: string;
}

export function FrameworkBadge({ framework, className = '' }: FrameworkBadgeProps) {
  const norm = framework.toLowerCase();

  switch (norm) {
    case 'laravel':
      return (
        <Badge
          variant="outline"
          className={`border-red-200 bg-red-500/10 text-red-600 dark:text-red-400 gap-1 text-xs font-medium ${className}`}
        >
          <Zap className="h-3 w-3 fill-red-500" />
          Laravel
        </Badge>
      );
    case 'wordpress':
      return (
        <Badge
          variant="outline"
          className={`border-blue-200 bg-blue-500/10 text-blue-600 dark:text-blue-400 gap-1 text-xs font-medium ${className}`}
        >
          <Globe className="h-3 w-3" />
          WordPress
        </Badge>
      );
    case 'nextjs':
      return (
        <Badge
          variant="outline"
          className={`border-zinc-300 bg-zinc-500/10 text-zinc-900 dark:text-zinc-100 gap-1 text-xs font-medium ${className}`}
        >
          <Sparkles className="h-3 w-3" />
          Next.js
        </Badge>
      );
    case 'nuxt':
      return (
        <Badge
          variant="outline"
          className={`border-emerald-200 bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 gap-1 text-xs font-medium ${className}`}
        >
          <Layers className="h-3 w-3" />
          Nuxt
        </Badge>
      );
    case 'react':
      return (
        <Badge
          variant="outline"
          className={`border-cyan-200 bg-cyan-500/10 text-cyan-600 dark:text-cyan-400 gap-1 text-xs font-medium ${className}`}
        >
          <Code2 className="h-3 w-3" />
          React / Vite
        </Badge>
      );
    case 'vue':
      return (
        <Badge
          variant="outline"
          className={`border-emerald-200 bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 gap-1 text-xs font-medium ${className}`}
        >
          <Code2 className="h-3 w-3" />
          Vue
        </Badge>
      );
    case 'php':
      return (
        <Badge
          variant="outline"
          className={`border-indigo-200 bg-indigo-500/10 text-indigo-600 dark:text-indigo-400 gap-1 text-xs font-medium ${className}`}
        >
          <FileCode className="h-3 w-3" />
          PHP
        </Badge>
      );
    case 'static':
      return (
        <Badge
          variant="outline"
          className={`border-amber-200 bg-amber-500/10 text-amber-600 dark:text-amber-400 gap-1 text-xs font-medium ${className}`}
        >
          <FileCode className="h-3 w-3" />
          HTML / Static
        </Badge>
      );
    default:
      return (
        <Badge variant="outline" className={`gap-1 text-xs font-normal ${className}`}>
          <Code2 className="h-3 w-3 text-muted-foreground" />
          {framework || 'Project'}
        </Badge>
      );
  }
}
