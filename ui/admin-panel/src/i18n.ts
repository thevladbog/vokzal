import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import ru from './locales/ru.json';
import en from './locales/en.json';

const resources = {
  ru: { translation: ru },
  en: { translation: en },
};

const savedLang = localStorage.getItem('admin-panel-lang') as 'ru' | 'en' | null;
const lng = savedLang && (savedLang === 'ru' || savedLang === 'en') ? savedLang : 'ru';

i18n.use(initReactI18next).init({
  resources,
  lng,
  fallbackLng: 'ru',
  interpolation: {
    escapeValue: false,
  },
});

i18n.on('languageChanged', (lng) => {
  localStorage.setItem('admin-panel-lang', lng);
});

export default i18n;
