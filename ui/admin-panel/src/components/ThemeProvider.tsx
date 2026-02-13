import React, { ReactNode } from 'react';
import { FluentProvider } from '@fluentui/react-components';
import { useThemeStore } from '@/stores/themeStore';
import { vokzalLightTheme, vokzalDarkTheme } from '@/constants/theme';
import { ThemeContext, ThemeContextValue } from '@/contexts/ThemeContext';

interface ThemeProviderProps {
  children: ReactNode;
}

export const ThemeProvider: React.FC<ThemeProviderProps> = ({ children }) => {
  const mode = useThemeStore((state) => state.mode);
  const setMode = useThemeStore((state) => state.setMode);
  const toggleMode = useThemeStore((state) => state.toggleMode);

  const theme = mode === 'dark' ? vokzalDarkTheme : vokzalLightTheme;

  const value: ThemeContextValue = {
    theme,
    mode,
    setMode,
    toggleMode,
  };

  return (
    <ThemeContext.Provider value={value}>
      <FluentProvider theme={theme}>{children}</FluentProvider>
    </ThemeContext.Provider>
  );
};
