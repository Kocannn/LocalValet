import type { Project } from './types';

export interface ProjectRepository {
  getProjects(): Promise<Project[]>;
  scanProjects(): Promise<Project[]>;
  getProjectRoots(): Promise<string[]>;
  addProjectRoot(path: string): Promise<void>;
  removeProjectRoot(path: string): Promise<void>;
  toggleProjectVHost(projectPath: string, enable: boolean): Promise<void>;
  generateProjectSSL(projectPath: string): Promise<void>;
  openInEditor(projectPath: string, editor?: string): Promise<void>;
  openInBrowser(url: string): Promise<void>;
  openInTerminal(projectPath: string): Promise<void>;
}
