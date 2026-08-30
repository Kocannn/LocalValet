import { useCallback, useEffect, useMemo, useState } from 'react';
import type { Project, ProjectViewMode } from '../domain/types';
import { WailsProjectRepository } from '../infrastructure/wails-project.repository';

export function useProjects() {
  const repository = useMemo(() => new WailsProjectRepository(), []);

  const [projects, setProjects] = useState<Project[]>([]);
  const [roots, setRoots] = useState<string[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isScanning, setIsScanning] = useState(false);
  const [viewMode, setViewMode] = useState<ProjectViewMode>('grid');
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedFramework, setSelectedFramework] = useState('All');
  const [rootsDialogOpen, setRootsDialogOpen] = useState(false);

  // Load initial data
  const loadData = useCallback(async () => {
    setIsLoading(true);
    try {
      const [projs, rootPaths] = await Promise.all([
        repository.getProjects(),
        repository.getProjectRoots(),
      ]);
      setProjects(projs);
      setRoots(rootPaths);
    } catch (err) {
      console.error('Failed to load projects data:', err);
    } finally {
      setIsLoading(false);
    }
  }, [repository]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  // Rescan projects
  const rescan = useCallback(async () => {
    setIsScanning(true);
    try {
      const scanned = await repository.scanProjects();
      setProjects(scanned);
    } catch (err) {
      console.error('Failed to scan projects:', err);
    } finally {
      setIsScanning(false);
    }
  }, [repository]);

  // Add Root Path
  const addRoot = useCallback(
    async (path: string) => {
      await repository.addProjectRoot(path);
      const [projs, rootPaths] = await Promise.all([
        repository.getProjects(),
        repository.getProjectRoots(),
      ]);
      setProjects(projs);
      setRoots(rootPaths);
    },
    [repository],
  );

  // Remove Root Path
  const removeRoot = useCallback(
    async (path: string) => {
      await repository.removeProjectRoot(path);
      const [projs, rootPaths] = await Promise.all([
        repository.getProjects(),
        repository.getProjectRoots(),
      ]);
      setProjects(projs);
      setRoots(rootPaths);
    },
    [repository],
  );

  // Toggle VHost
  const toggleVHost = useCallback(
    async (path: string, enable: boolean) => {
      // Optimistic update
      setProjects((prev) =>
        prev.map((p) => (p.path === path ? { ...p, vhostEnabled: enable } : p)),
      );
      try {
        await repository.toggleProjectVHost(path, enable);
      } catch (err) {
        console.error('Failed to toggle VHost:', err);
        // Rollback on failure
        setProjects((prev) =>
          prev.map((p) => (p.path === path ? { ...p, vhostEnabled: !enable } : p)),
        );
      }
    },
    [repository],
  );

  // Generate SSL
  const generateSSL = useCallback(
    async (path: string) => {
      try {
        await repository.generateProjectSSL(path);
        setProjects((prev) =>
          prev.map((p) => (p.path === path ? { ...p, sslEnabled: true } : p)),
        );
      } catch (err) {
        console.error('Failed to generate SSL:', err);
      }
    },
    [repository],
  );

  // Open Actions
  const openBrowser = useCallback(
    async (url: string) => {
      await repository.openInBrowser(url);
    },
    [repository],
  );

  const openTerminal = useCallback(
    async (path: string) => {
      await repository.openInTerminal(path);
    },
    [repository],
  );

  const openEditor = useCallback(
    async (path: string) => {
      await repository.openInEditor(path);
    },
    [repository],
  );

  // Framework filters and search
  const filteredProjects = useMemo(() => {
    return projects.filter((p) => {
      const matchesSearch =
        searchQuery === '' ||
        p.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        p.domain.toLowerCase().includes(searchQuery.toLowerCase()) ||
        p.path.toLowerCase().includes(searchQuery.toLowerCase());

      const matchesFramework =
        selectedFramework === 'All' ||
        p.framework.toLowerCase() === selectedFramework.toLowerCase();

      return matchesSearch && matchesFramework;
    });
  }, [projects, searchQuery, selectedFramework]);

  // Framework stats
  const frameworkCategories = useMemo(() => {
    const counts: Record<string, number> = { All: projects.length };
    for (const p of projects) {
      const f = p.framework.toLowerCase();
      counts[f] = (counts[f] || 0) + 1;
    }
    return counts;
  }, [projects]);

  return {
    projects,
    roots,
    filteredProjects,
    frameworkCategories,
    isLoading,
    isScanning,
    viewMode,
    setViewMode,
    searchQuery,
    setSearchQuery,
    selectedFramework,
    setSelectedFramework,
    rootsDialogOpen,
    setRootsDialogOpen,
    rescan,
    addRoot,
    removeRoot,
    toggleVHost,
    generateSSL,
    openBrowser,
    openTerminal,
    openEditor,
  };
}
