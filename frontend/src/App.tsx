/**
 * App Component
 * 
 * Root application component with routing configuration.
 * Sets up React Router with nested routes and layouts.
 * 
 * Architecture: Feature-based with routing
 * Pattern: Route-based code splitting
 */

import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { RootLayout } from '@/layouts';
import { HomePage, ServicesPage, SettingsPage } from '@/pages';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<RootLayout />}>
          <Route path="/" element={<HomePage />} />
          <Route path="/services" element={<ServicesPage />} />
          <Route path="/settings" element={<SettingsPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
