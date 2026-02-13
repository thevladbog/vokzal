import { create } from 'zustand';

type ThemeMode = 'light' | 'dark';

interface ThemeStore {
  mode: ThemeMode;
  setMode: (mode: ThemeMode) => void;
  toggleMode: () => void;
}

const THEME_STORAGE_KEY = 'vokzal-theme';

function getSavedTheme(): ThemeMode {
  try {
    if (typeof window === 'undefined' || typeof window.localStorage === 'undefined') {
      return 'light';
    }
    const saved = window.localStorage.getItem(THEME_STORAGE_KEY);
    if (saved === 'light' || saved === 'dark') return saved;
  } catch {
    // localStorage unavailable
  }
  return 'light';
}

export const useThemeStore = create<ThemeStore>()((set) => ({
  mode: getSavedTheme(),

  setMode: (mode) => {
    set({ mode });
    try {
      if (typeof window !== 'undefined' && window.localStorage) {
        window.localStorage.setItem(THEME_STORAGE_KEY, mode);
      }
    } catch {
      // ignore
    }
  },

  toggleMode: () => {
    set((state) => {
      const newMode = state.mode === 'light' ? 'dark' : 'light';
      try {
        if (typeof window !== 'undefined' && window.localStorage) {
          window.localStorage.setItem(THEME_STORAGE_KEY, newMode);
        }
      } catch {
        // ignore
      }
      return { mode: newMode };
    });
  },
}));
