import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import ru from './locales/ru.json';
import en from './locales/en.json';

const resources = {
  ru: { translation: ru },
  en: { translation: en },
};

function getSavedLanguage(): 'ru' | 'en' {
  try {
    if (typeof window === 'undefined' || typeof window.localStorage === 'undefined') {
      return 'ru';
    }
    const saved = window.localStorage.getItem('admin-panel-lang');
    if (saved === 'ru' || saved === 'en') return saved;
  } catch {
    // localStorage unavailable (private mode, quota, etc.)
  }
  return 'ru';
}

const lng = getSavedLanguage();

i18n.use(initReactI18next).init({
  resources,
  lng,
  fallbackLng: 'ru',
  interpolation: {
    escapeValue: false,
  },
});

i18n.on('languageChanged', (lng) => {
  try {
    if (typeof window !== 'undefined' && window.localStorage) {
      window.localStorage.setItem('admin-panel-lang', lng);
    }
  } catch {
    // ignore
  }
});

export default i18n;
