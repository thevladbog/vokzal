import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { FluentProvider, webLightTheme } from '@fluentui/react-components';
import { LoginPage } from '@/pages/LoginPage';
import { SalePage } from '@/pages/SalePage';
import { RefundPage } from '@/pages/RefundPage';
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
                isAuthenticated ? <SalePage /> : <Navigate to="/login" replace />
              }
            />

            <Route
              path="/refund"
              element={
                isAuthenticated ? <RefundPage /> : <Navigate to="/login" replace />
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
