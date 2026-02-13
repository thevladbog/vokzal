import { Theme, webLightTheme, webDarkTheme } from '@fluentui/react-components';

// Вокзал.ТЕХ brand colors from brand-colors.css
export const BRAND_COLORS = {
  primary: '#2563EB',
  primaryDark: '#1D4ED8',
  primaryLight: '#3B82F6',
  secondary: '#10B981',
  secondaryDark: '#059669',
  accent: '#F59E0B',
  accentDark: '#D97706',
};

// Custom theme type with brand tokens
export type VokzalTheme = Theme & {
  brandPrimary: string;
  brandSecondary: string;
  brandAccent: string;
};

// Extend Fluent UI light theme with Vokzal.TEH branding
export const vokzalLightTheme: VokzalTheme = {
  ...webLightTheme,
  brandPrimary: BRAND_COLORS.primary,
  brandSecondary: BRAND_COLORS.secondary,
  brandAccent: BRAND_COLORS.accent,
  // Override Fluent UI brand colors with Vokzal.TEH colors
  colorBrandBackground: BRAND_COLORS.primary,
  colorBrandBackgroundHover: BRAND_COLORS.primaryDark,
  colorBrandBackgroundPressed: BRAND_COLORS.primaryDark,
  colorBrandForeground1: BRAND_COLORS.primary,
  colorBrandForeground2: BRAND_COLORS.primaryLight,
};

// Extend Fluent UI dark theme with Vokzal.TEH branding
export const vokzalDarkTheme: VokzalTheme = {
  ...webDarkTheme,
  brandPrimary: BRAND_COLORS.primary,
  brandSecondary: BRAND_COLORS.secondary,
  brandAccent: BRAND_COLORS.accent,
  // Override Fluent UI brand colors with Vokzal.TEH colors
  colorBrandBackground: BRAND_COLORS.primaryLight,
  colorBrandBackgroundHover: BRAND_COLORS.primary,
  colorBrandBackgroundPressed: BRAND_COLORS.primaryDark,
  colorBrandForeground1: BRAND_COLORS.primaryLight,
  colorBrandForeground2: BRAND_COLORS.primary,
};
