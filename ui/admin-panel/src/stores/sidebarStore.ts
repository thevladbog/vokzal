import { create } from 'zustand';

interface SidebarStore {
  isCollapsed: boolean;
  setCollapsed: (collapsed: boolean) => void;
  toggleCollapsed: () => void;
}

const SIDEBAR_STORAGE_KEY = 'vokzal-sidebar-collapsed';

function getSavedCollapsedState(): boolean {
  try {
    if (typeof window === 'undefined' || typeof window.localStorage === 'undefined') {
      return false;
    }
    const saved = window.localStorage.getItem(SIDEBAR_STORAGE_KEY);
    return saved === 'true';
  } catch {
    return false;
  }
}

export const useSidebarStore = create<SidebarStore>()((set) => ({
  isCollapsed: getSavedCollapsedState(),

  setCollapsed: (collapsed) => {
    set({ isCollapsed: collapsed });
    try {
      if (typeof window !== 'undefined' && window.localStorage) {
        window.localStorage.setItem(SIDEBAR_STORAGE_KEY, String(collapsed));
      }
    } catch {
      // ignore
    }
  },

  toggleCollapsed: () => {
    set((state) => {
      const newCollapsed = !state.isCollapsed;
      try {
        if (typeof window !== 'undefined' && window.localStorage) {
          window.localStorage.setItem(SIDEBAR_STORAGE_KEY, String(newCollapsed));
        }
      } catch {
        // ignore
      }
      return { isCollapsed: newCollapsed };
    });
  },
}));
