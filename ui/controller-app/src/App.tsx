import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { FluentProvider, webLightTheme } from '@fluentui/react-components';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { LoginPage } from './pages/LoginPage';
import { TripSelectionPage } from './pages/TripSelectionPage';
import { ScanPage } from './pages/ScanPage';
import { ProtectedRoute } from './components/ProtectedRoute';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
    },
  },
});

const App: React.FC = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <FluentProvider theme={webLightTheme}>
        <BrowserRouter>
          <Routes>
            <Route path="/login" element={<LoginPage />} />
            <Route
              path="/trips"
              element={
                <ProtectedRoute requiredRole="controller">
                  <TripSelectionPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="/scan"
              element={
                <ProtectedRoute requiredRole="controller">
                  <ScanPage />
                </ProtectedRoute>
              }
            />
            <Route path="*" element={<Navigate to="/trips" replace />} />
          </Routes>
        </BrowserRouter>
      </FluentProvider>
    </QueryClientProvider>
  );
};

export default App;
