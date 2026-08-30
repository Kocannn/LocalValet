export type FrameworkType =
  | 'laravel'
  | 'wordpress'
  | 'nextjs'
  | 'nuxt'
  | 'react'
  | 'vue'
  | 'php'
  | 'static'
  | 'unknown';

export interface Project {
  id: string;
  name: string;
  path: string;
  framework: FrameworkType | string;
  webRoot: string;
  domain: string;
  vhostEnabled: boolean;
  sslEnabled: boolean;
  targetPort?: number;
  phpVersion?: string;
  nodeVersion?: string;
  createdAt?: string;
  updatedAt?: string;
}

export type ProjectViewMode = 'grid' | 'list';
