import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { FluentProvider, webLightTheme, Toaster } from '@fluentui/react-components';
import { ProtectedRoute } from '@/components/ProtectedRoute';
import { LoginPage } from '@/pages/LoginPage';
import { DashboardPage } from '@/pages/DashboardPage';
import { SchedulesPage } from '@/pages/SchedulesPage';
import { UsersPage } from '@/pages/UsersPage';
import { StationsPage } from '@/pages/StationsPage';
import { RoutesPage } from '@/pages/RoutesPage';
import { TripsPage } from '@/pages/TripsPage';
import { AuditPage } from '@/pages/AuditPage';
import { ReportsPage } from '@/pages/ReportsPage';
import { MonitoringPage } from '@/pages/MonitoringPage';
import { useAuthStore } from '@/stores/authStore';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
});

function App() {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);

  return (
    <QueryClientProvider client={queryClient}>
      <FluentProvider theme={webLightTheme}>
        <Toaster />
        <BrowserRouter>
          <Routes>
            <Route 
              path="/login" 
              element={
                isAuthenticated ? <Navigate to="/" replace /> : <LoginPage />
              } 
            />
            
            <Route
              path="/"
              element={
                <ProtectedRoute>
                  <DashboardPage />
                </ProtectedRoute>
              }
            />

            <Route
              path="/schedules"
              element={
                <ProtectedRoute allowedRoles={['admin', 'dispatcher']}>
                  <SchedulesPage />
                </ProtectedRoute>
              }
            />

            <Route
              path="/users"
              element={
                <ProtectedRoute allowedRoles={['admin']}>
                  <UsersPage />
                </ProtectedRoute>
              }
            />

            <Route
              path="/stations"
              element={
                <ProtectedRoute allowedRoles={['admin', 'dispatcher']}>
                  <StationsPage />
                </ProtectedRoute>
              }
            />

            <Route
              path="/routes"
              element={
                <ProtectedRoute allowedRoles={['admin', 'dispatcher']}>
                  <RoutesPage />
                </ProtectedRoute>
              }
            />

            <Route
              path="/trips"
              element={
                <ProtectedRoute allowedRoles={['admin', 'dispatcher']}>
                  <TripsPage />
                </ProtectedRoute>
              }
            />

            <Route
              path="/audit"
              element={
                <ProtectedRoute allowedRoles={['admin']}>
                  <AuditPage />
                </ProtectedRoute>
              }
            />

            <Route
              path="/reports"
              element={
                <ProtectedRoute allowedRoles={['admin', 'dispatcher']}>
                  <ReportsPage />
                </ProtectedRoute>
              }
            />

            <Route
              path="/monitoring"
              element={
                <ProtectedRoute allowedRoles={['admin']}>
                  <MonitoringPage />
                </ProtectedRoute>
              }
            />

            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </BrowserRouter>
      </FluentProvider>
    </QueryClientProvider>
  );
}

export default App;
