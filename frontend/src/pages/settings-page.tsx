/**
 * Settings Page
 * 
 * Application settings, runtime version management, and system-level tools (SSL CA, DNS).
 */

import { useEffect, useState } from 'react';
import {
  Check,
  CheckCircle2,
  Globe,
  Layers,
  Lock,
  RefreshCw,
  ShieldAlert,
  ShieldCheck,
  Sparkles,
} from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { ServiceRuntimeCard, useRuntimeSettings } from '@/modules/settings';
import * as WailsApp from '../../wailsjs/go/main/App.js';

export function SettingsPage() {
  const { services, isLoading, savingMap, errorMap, changeVersion, refresh } = useRuntimeSettings();
  const [selectedCategory, setSelectedCategory] = useState<string>('All');

  // SSL and DNS state
  const [isCATrusted, setIsCATrusted] = useState<boolean>(false);
  const [isTrustingCA, setIsTrustingCA] = useState<boolean>(false);
  const [isSyncingHosts, setIsSyncingHosts] = useState<boolean>(false);
  const [caFeedback, setCaFeedback] = useState<string>('');
  const [dnsFeedback, setDnsFeedback] = useState<string>('');

  const checkCAStatus = async () => {
    try {
      const trusted = await (WailsApp as any).IsRootCATrusted();
      setIsCATrusted(!!trusted);
    } catch {
      setIsCATrusted(false);
    }
  };

  useEffect(() => {
    checkCAStatus();
  }, []);

  const handleTrustCA = async () => {
    setIsTrustingCA(true);
    setCaFeedback('');
    try {
      const res = await (WailsApp as any).TrustRootCA();
      setCaFeedback(res?.message || 'Local Root CA processed.');
      await checkCAStatus();
    } catch (err: any) {
      setCaFeedback(err?.message || 'Failed to trust CA.');
    } finally {
      setIsTrustingCA(false);
    }
  };

  const handleSyncHosts = async () => {
    setIsSyncingHosts(true);
    setDnsFeedback('');
    try {
      const res = await (WailsApp as any).SyncHostsDomains();
      setDnsFeedback(res?.message || 'Domains synchronized.');
    } catch (err: any) {
      setDnsFeedback(err?.message || 'Failed to sync hosts.');
    } finally {
      setIsSyncingHosts(false);
    }
  };

  const categories = ['All', 'Runtime', 'Database', 'Web'];

  const filteredServices = selectedCategory === 'All'
    ? services
    : services.filter((s) => s.category.toLowerCase() === selectedCategory.toLowerCase());

  return (
    <div className="flex flex-col gap-6">
      {/* Header Banner */}
      <Card className="border-border/70 bg-card/95">
        <CardHeader>
          <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <div className="flex items-center gap-2">
                <CardTitle className="text-base md:text-lg">Settings & System Tools</CardTitle>
                <Badge variant="outline" className="border-primary/20 bg-primary/5 text-primary text-xs">
                  Production Ready
                </Badge>
              </div>
              <CardDescription className="mt-1">
                Configure runtime versions, manage system Root CA trust, and sync local development domains.
              </CardDescription>
            </div>
            <Button
              variant="outline"
              size="sm"
              className="self-start gap-1.5 sm:self-auto"
              disabled={isLoading}
              onClick={refresh}
            >
              <RefreshCw className={`h-3.5 w-3.5 ${isLoading ? 'animate-spin' : ''}`} />
              Refresh
            </Button>
          </div>
        </CardHeader>
      </Card>

      {/* System Integration Section: SSL & DNS */}
      <div className="grid gap-4 md:grid-cols-2">
        {/* Root CA Trust Card */}
        <Card className="border-border/70 bg-card/95">
          <CardHeader className="pb-3">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <CardTitle className="text-base font-semibold">Local Root CA (HTTPS)</CardTitle>
                {isCATrusted ? (
                  <Badge variant="outline" className="border-emerald-200 bg-emerald-50 text-emerald-700 dark:border-emerald-800 dark:bg-emerald-950/40 dark:text-emerald-400 gap-1 text-xs">
                    <ShieldCheck className="h-3 w-3" />
                    Trusted
                  </Badge>
                ) : (
                  <Badge variant="outline" className="border-amber-200 bg-amber-50 text-amber-700 dark:border-amber-800 dark:bg-amber-950/40 dark:text-amber-400 gap-1 text-xs">
                    <ShieldAlert className="h-3 w-3" />
                    Untrusted
                  </Badge>
                )}
              </div>
            </div>
            <CardDescription className="text-xs">
              Installs LocalValet Root CA into <code className="font-mono text-foreground/80">/usr/local/share/ca-certificates/</code> so <code className="font-mono text-foreground/80">https://*.test</code> is trusted across your system.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-3">
            <div className="flex items-center gap-2">
              <Button
                variant={isCATrusted ? 'outline' : 'default'}
                size="sm"
                className="gap-1.5 text-xs"
                disabled={isTrustingCA}
                onClick={handleTrustCA}
              >
                <Lock className="h-3.5 w-3.5" />
                {isTrustingCA ? 'Installing...' : isCATrusted ? 'Re-install Root CA' : 'Trust Root CA in System'}
              </Button>
            </div>
            {caFeedback && (
              <p className="text-xs font-medium text-primary">{caFeedback}</p>
            )}
            <p className="text-[11px] text-muted-foreground">
              CA certificate file location: <code className="font-mono">runtime/certs/ca.crt</code>
            </p>
          </CardContent>
        </Card>

        {/* DNS /etc/hosts Sync Card */}
        <Card className="border-border/70 bg-card/95">
          <CardHeader className="pb-3">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <CardTitle className="text-base font-semibold">Local Domain Resolution</CardTitle>
                <Badge variant="outline" className="text-xs font-normal">/etc/hosts</Badge>
              </div>
            </div>
            <CardDescription className="text-xs">
              Synchronizes discovered <code className="font-mono text-foreground/80">*.test</code> project domains to <code className="font-mono text-foreground/80">/etc/hosts</code> inside a managed block.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-3">
            <div className="flex items-center gap-2">
              <Button
                variant="outline"
                size="sm"
                className="gap-1.5 text-xs"
                disabled={isSyncingHosts}
                onClick={handleSyncHosts}
              >
                <Globe className="h-3.5 w-3.5" />
                {isSyncingHosts ? 'Syncing...' : 'Sync Domains to /etc/hosts'}
              </Button>
            </div>
            {dnsFeedback && (
              <p className="text-xs font-medium text-primary">{dnsFeedback}</p>
            )}
            <p className="text-[11px] text-muted-foreground">
              Safe operation: updates only the <code className="font-mono"># LocalValet Managed</code> section.
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Category Filter Pills */}
      <div className="flex flex-wrap items-center gap-2 pt-2">
        <span className="text-xs font-medium text-muted-foreground mr-1">Filter Runtimes:</span>
        {categories.map((category) => (
          <Button
            key={category}
            variant={selectedCategory === category ? 'default' : 'outline'}
            size="sm"
            className="h-7 text-xs px-3 rounded-full"
            onClick={() => setSelectedCategory(category)}
          >
            {category}
          </Button>
        ))}
      </div>

      {/* Runtime Services Grid */}
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {filteredServices.map((service) => (
          <ServiceRuntimeCard
            key={service.serviceName}
            service={service}
            isSaving={savingMap[service.serviceName] || false}
            error={errorMap[service.serviceName]}
            onVersionChange={changeVersion}
          />
        ))}
      </div>

      {filteredServices.length === 0 && !isLoading && (
        <div className="py-12 text-center text-muted-foreground">
          <Layers className="mx-auto h-8 w-8 opacity-40 mb-2" />
          <p className="text-sm">No services found for the selected category.</p>
        </div>
      )}

      {/* Additional Settings Info */}
      <div className="grid gap-4 md:grid-cols-2 pt-2">
        <Card className="border-border/70 bg-card/95">
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle className="text-base">Terminal & PATH Environment</CardTitle>
              <Badge variant="outline" className="text-xs font-normal">Active</Badge>
            </div>
            <CardDescription className="text-xs">
              Context terminal automatically injects active PHP, Node.js, Composer, and MySQL binaries into PATH.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <p className="text-xs text-muted-foreground">
              Hot-restarting a service updates its environment instantly in newly launched terminal sessions.
            </p>
          </CardContent>
        </Card>

        <Card className="border-border/70 bg-card/95">
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle className="text-base">Custom Binary Locations</CardTitle>
              <Badge variant="secondary" className="text-xs">Isolated (runtime/)</Badge>
            </div>
            <CardDescription className="text-xs">
              Binaries are isolated in <code className="font-mono text-xs">runtime/linux/&lt;service&gt;/&lt;version&gt;</code>.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <p className="text-xs text-muted-foreground">
              You can drop new version folders into <code className="font-mono text-xs">runtime/linux/</code> and LocalValet will detect them automatically.
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
