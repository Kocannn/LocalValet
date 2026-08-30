import type { ProjectRepository } from '../domain/ports';
import type { Project } from '../domain/types';
import * as WailsApp from '../../../../wailsjs/go/main/App.js';

export class WailsProjectRepository implements ProjectRepository {
  async getProjects(): Promise<Project[]> {
    try {
      const data = await (WailsApp as any).GetProjects();
      return (data ?? []).map(this.mapProject);
    } catch (error) {
      console.error('Failed to get projects:', error);
      return [];
    }
  }

  async scanProjects(): Promise<Project[]> {
    try {
      const data = await (WailsApp as any).ScanProjects();
      return (data ?? []).map(this.mapProject);
    } catch (error) {
      console.error('Failed to scan projects:', error);
      return [];
    }
  }

  async getProjectRoots(): Promise<string[]> {
    try {
      return (await (WailsApp as any).GetProjectRoots()) ?? [];
    } catch (error) {
      console.error('Failed to get project roots:', error);
      return [];
    }
  }

  async addProjectRoot(path: string): Promise<void> {
    try {
      await (WailsApp as any).AddProjectRoot(path);
    } catch (error) {
      console.error(`Failed to add root ${path}:`, error);
      throw error;
    }
  }

  async removeProjectRoot(path: string): Promise<void> {
    try {
      await (WailsApp as any).RemoveProjectRoot(path);
    } catch (error) {
      console.error(`Failed to remove root ${path}:`, error);
      throw error;
    }
  }

  async toggleProjectVHost(projectPath: string, enable: boolean): Promise<void> {
    try {
      await (WailsApp as any).ToggleProjectVHost(projectPath, enable);
    } catch (error) {
      console.error(`Failed to toggle VHost for ${projectPath}:`, error);
      throw error;
    }
  }

  async generateProjectSSL(projectPath: string): Promise<void> {
    try {
      await (WailsApp as any).GenerateProjectSSL(projectPath);
    } catch (error) {
      console.error(`Failed to generate SSL for ${projectPath}:`, error);
      throw error;
    }
  }

  async openInEditor(projectPath: string, editor = ''): Promise<void> {
    try {
      await (WailsApp as any).OpenProjectInEditor(projectPath, editor);
    } catch (error) {
      console.error(`Failed to open ${projectPath} in editor:`, error);
      throw error;
    }
  }

  async openInBrowser(url: string): Promise<void> {
    try {
      await (WailsApp as any).OpenProjectInBrowser(url);
    } catch (error) {
      console.error(`Failed to open ${url} in browser:`, error);
      throw error;
    }
  }

  async openInTerminal(projectPath: string): Promise<void> {
    try {
      await (WailsApp as any).OpenContextTerminal(projectPath);
    } catch (error) {
      console.error(`Failed to open terminal in ${projectPath}:`, error);
      throw error;
    }
  }

  private mapProject(p: any): Project {
    return {
      id: p.id ?? p.ID ?? '',
      name: p.name ?? p.Name ?? '',
      path: p.path ?? p.Path ?? '',
      framework: p.framework ?? p.Framework ?? 'unknown',
      webRoot: p.webRoot ?? p.WebRoot ?? '.',
      domain: p.domain ?? p.Domain ?? '',
      vhostEnabled: p.vhostEnabled ?? p.VHostEnabled ?? false,
      sslEnabled: p.sslEnabled ?? p.SSLEnabled ?? false,
      targetPort: p.targetPort ?? p.TargetPort,
      phpVersion: p.phpVersion ?? p.PHPVersion,
      nodeVersion: p.nodeVersion ?? p.NodeVersion,
      createdAt: p.createdAt ?? p.CreatedAt,
      updatedAt: p.updatedAt ?? p.UpdatedAt,
    };
  }
}
