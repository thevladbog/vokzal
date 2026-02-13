import { createContext } from 'react';
import type { VokzalTheme } from '@/constants/theme';

export interface ThemeContextValue {
  theme: VokzalTheme;
  mode: 'light' | 'dark';
  setMode: (mode: 'light' | 'dark') => void;
  toggleMode: () => void;
}

export const ThemeContext = createContext<ThemeContextValue | undefined>(undefined);
