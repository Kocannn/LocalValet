/**
 * Projects Page
 * 
 * Auto-discovered local development projects dashboard.
 * Provides virtual host routing, SSL management, and quick IDE/terminal launches.
 */

import {
  FolderGit2,
  FolderRoot,
  LayoutGrid,
  List,
  Plus,
  RefreshCw,
  Search,
} from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
  ProjectCard,
  ProjectListRow,
  ProjectRootsDialog,
  useProjects,
} from '@/modules/project';

const FRAMEWORK_PILLS = [
  'All',
  'laravel',
  'wordpress',
  'nextjs',
  'react',
  'vue',
  'php',
  'static',
];

export function ProjectsPage() {
  const {
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
  } = useProjects();

  return (
    <div className="flex flex-col gap-5">
      {/* Top Banner Card */}
      <Card className="border-border/70 bg-card/95">
        <CardHeader className="pb-4">
          <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <div className="space-y-1">
              <div className="flex items-center gap-2">
                <CardTitle className="text-base md:text-lg">Project Discovery & Virtual Hosts</CardTitle>
                <Badge variant="outline" className="border-primary/20 bg-primary/5 text-primary text-xs">
                  Phase 4
                </Badge>
                <Badge variant="secondary" className="text-xs font-mono">
                  {projects.length} {projects.length === 1 ? 'Project' : 'Projects'}
                </Badge>
              </div>
              <CardDescription className="text-xs">
                Discovered projects are served automatically at <code className="font-mono text-foreground/90 font-medium">*.test</code> with FastCGI / reverse-proxy and local SSL.
              </CardDescription>
            </div>

            <div className="flex items-center gap-2 self-start sm:self-auto">
              <Button
                variant="outline"
                size="sm"
                className="h-8 gap-1.5 text-xs"
                onClick={() => setRootsDialogOpen(true)}
              >
                <FolderRoot className="h-3.5 w-3.5" />
                Scan Folders ({roots.length})
              </Button>
              <Button
                variant="default"
                size="sm"
                className="h-8 gap-1.5 text-xs"
                disabled={isScanning || isLoading}
                onClick={rescan}
              >
                <RefreshCw className={`h-3.5 w-3.5 ${isScanning ? 'animate-spin' : ''}`} />
                {isScanning ? 'Scanning...' : 'Rescan'}
              </Button>
            </div>
          </div>
        </CardHeader>
      </Card>

      {/* Filter and View Controls Toolbar */}
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        {/* Search Input */}
        <div className="relative w-full sm:w-72">
          <Search className="absolute left-2.5 top-2.5 h-3.5 w-3.5 text-muted-foreground" />
          <Input
            placeholder="Search project or domain..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="h-8 pl-8 text-xs"
          />
        </div>

        {/* View Mode Toggle (Grid vs List) */}
        <div className="flex items-center gap-1 self-end sm:self-auto">
          <Button
            variant={viewMode === 'grid' ? 'secondary' : 'ghost'}
            size="sm"
            className="h-8 w-8 p-0"
            onClick={() => setViewMode('grid')}
            title="Grid View"
          >
            <LayoutGrid className="h-4 w-4" />
          </Button>
          <Button
            variant={viewMode === 'list' ? 'secondary' : 'ghost'}
            size="sm"
            className="h-8 w-8 p-0"
            onClick={() => setViewMode('list')}
            title="List View"
          >
            <List className="h-4 w-4" />
          </Button>
        </div>
      </div>

      {/* Framework Category Pills */}
      <div className="flex flex-wrap items-center gap-1.5">
        {FRAMEWORK_PILLS.map((pill) => {
          const count = frameworkCategories[pill.toLowerCase()] || 0;
          const isSelected = selectedFramework.toLowerCase() === pill.toLowerCase();

          if (pill !== 'All' && count === 0) return null;

          return (
            <Button
              key={pill}
              variant={isSelected ? 'default' : 'outline'}
              size="sm"
              className="h-7 text-xs px-2.5 rounded-full capitalize gap-1"
              onClick={() => setSelectedFramework(pill)}
            >
              <span>{pill}</span>
              <span className={`text-[10px] font-mono px-1 rounded-full ${isSelected ? 'bg-primary-foreground/20' : 'bg-muted'}`}>
                {count}
              </span>
            </Button>
          );
        })}
      </div>

      {/* Main Projects Display */}
      {viewMode === 'grid' ? (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {filteredProjects.map((project) => (
            <ProjectCard
              key={project.path}
              project={project}
              onOpenBrowser={openBrowser}
              onOpenTerminal={openTerminal}
              onOpenEditor={openEditor}
              onToggleVHost={toggleVHost}
              onGenerateSSL={generateSSL}
            />
          ))}
        </div>
      ) : (
        <div className="space-y-2.5">
          {filteredProjects.map((project) => (
            <ProjectListRow
              key={project.path}
              project={project}
              onOpenBrowser={openBrowser}
              onOpenTerminal={openTerminal}
              onOpenEditor={openEditor}
              onToggleVHost={toggleVHost}
              onGenerateSSL={generateSSL}
            />
          ))}
        </div>
      )}

      {/* Empty State */}
      {filteredProjects.length === 0 && !isLoading && (
        <div className="flex flex-col items-center justify-center py-16 text-center text-muted-foreground border border-dashed rounded-lg border-border/70 p-6">
          <FolderGit2 className="h-10 w-10 opacity-30 mb-3" />
          <p className="text-base font-medium text-foreground">No projects found</p>
          <p className="mt-1 text-xs max-w-sm">
            {searchQuery || selectedFramework !== 'All'
              ? 'No projects match your current filters. Try changing or clearing your search.'
              : 'Add your projects root folder (e.g. ~/Projects, ~/Coding) and click Rescan.'}
          </p>
          <div className="mt-4 flex gap-2">
            <Button
              variant="outline"
              size="sm"
              className="gap-1 text-xs"
              onClick={() => setRootsDialogOpen(true)}
            >
              <Plus className="h-3.5 w-3.5" /> Add Project Root
            </Button>
          </div>
        </div>
      )}

      {/* Scan Roots Dialog */}
      <ProjectRootsDialog
        open={rootsDialogOpen}
        roots={roots}
        onOpenChange={setRootsDialogOpen}
        onAddRoot={addRoot}
        onRemoveRoot={removeRoot}
      />
    </div>
  );
}
