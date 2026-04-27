/**
 * Settings Page
 * 
 * Application settings and configuration page.
 * Placeholder for future settings functionality.
 * 
 * Design Pattern: Page Component
 * Purpose: Settings and configuration management
 */

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { PhpRuntimeSelector, usePhpRuntimeSettings } from '@/modules/settings';

/**
 * Settings page
 */
export function SettingsPage() {
  const {
    versions,
    label,
    isLoading,
    isSaving,
    error,
    changeVersion,
  } = usePhpRuntimeSettings();

  return (
    <div className='flex flex-col gap-4'>
      <Card className="border-border/70 bg-card/95">
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle className="text-base md:text-lg">Settings</CardTitle>
              <CardDescription>
                Configure your LocalValet application
              </CardDescription>
            </div>
            <Badge variant="secondary">Coming Soon</Badge>
          </div>
        </CardHeader>
        <CardContent>
          <div className="py-10 text-center text-muted-foreground">
            <p className="text-lg font-medium tracking-tight">Settings module is in progress.</p>
            <p className="mt-2 text-sm">Configuration options will be available in a future update.</p>
          </div>
        </CardContent>
      </Card>

      {/* Placeholder for future settings sections */}
      <div className="grid gap-4 md:grid-cols-2">
        <Card className="border-border/70 bg-card/95">
          <CardHeader>
            <CardTitle className="text-base">General</CardTitle>
            <CardDescription>
              Application preferences and behavior
            </CardDescription>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground">Coming soon...</p>
          </CardContent>
        </Card>

        <Card className="border-border/70 bg-card/95">
          <CardHeader>
            <CardTitle className="text-base">Services</CardTitle>
            <CardDescription>
              Isolated runtime and service versions
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-3">
            <PhpRuntimeSelector
              versions={versions}
              label={label}
              isLoading={isLoading}
              isSaving={isSaving}
              error={error}
              onChange={changeVersion}
            />
          </CardContent>
        </Card>

        <Card className="border-border/70 bg-card/95">
          <CardHeader>
            <CardTitle className="text-base">Notifications</CardTitle>
            <CardDescription>
              Alert and notification preferences
            </CardDescription>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground">Coming soon...</p>
          </CardContent>
        </Card>

        <Card className="border-border/70 bg-card/95">
          <CardHeader>
            <CardTitle className="text-base">Advanced</CardTitle>
            <CardDescription>
              Advanced configuration options
            </CardDescription>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground">Coming soon...</p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
