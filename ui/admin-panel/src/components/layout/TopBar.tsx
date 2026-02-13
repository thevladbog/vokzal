import React from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import {
  makeStyles,
  tokens,
  Avatar,
  Menu,
  MenuTrigger,
  MenuPopover,
  MenuList,
  MenuItem,
  Button,
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbButton,
  Switch,
  shorthands,
} from '@fluentui/react-components';
import { SignOut20Regular, WeatherMoon20Regular, WeatherSunny20Regular, Navigation20Regular } from '@fluentui/react-icons';
import { useTheme } from '@/hooks/useTheme';
import { useAuthStore } from '@/stores/authStore';
import { useSidebarStore } from '@/stores/sidebarStore';
import { authService } from '@/services/auth';
import { ALL_NAV_ITEMS } from '@/constants/navigation';
import i18n from '@/i18n';

const useStyles = makeStyles({
  topbar: {
    height: '64px',
    backgroundColor: tokens.colorNeutralBackground1,
    ...shorthands.borderBottom('1px', 'solid', tokens.colorNeutralStroke2),
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    ...shorthands.padding('0', '24px'),
    ...shorthands.gap('16px'),
    '@media (max-width: 768px)': {
      ...shorthands.padding('0', '16px'),
      height: '56px',
    },
  },
  left: {
    display: 'flex',
    alignItems: 'center',
    ...shorthands.gap('16px'),
    ...shorthands.flex(1),
    minWidth: 0,
  },
  hamburgerMobile: {
    display: 'none',
    '@media (max-width: 768px)': {
      display: 'flex',
      backgroundColor: 'transparent',
      ...shorthands.border('none'),
      ...shorthands.padding('8px'),
      cursor: 'pointer',
      ...shorthands.borderRadius(tokens.borderRadiusMedium),
      alignItems: 'center',
      justifyContent: 'center',
      color: tokens.colorNeutralForeground2,
      ...shorthands.transition('all', '150ms', 'ease-in-out'),
      ':hover': {
        backgroundColor: tokens.colorNeutralBackground2Hover,
      },
      ':active': {
        transform: 'scale(0.95)',
      },
    },
  },
  breadcrumb: {
    ...shorthands.overflow('hidden'),
  },
  right: {
    display: 'flex',
    alignItems: 'center',
    ...shorthands.gap('16px'),
    '@media (max-width: 768px)': {
      ...shorthands.gap('8px'),
    },
  },
  themeSwitch: {
    display: 'flex',
    alignItems: 'center',
    ...shorthands.gap('8px'),
    '@media (max-width: 480px)': {
      ...shorthands.gap('4px'),
      '& svg': {
        display: 'none',
      },
    },
  },
  langButton: {
    minWidth: '48px',
  },
  userMenu: {
    display: 'flex',
    alignItems: 'center',
    ...shorthands.gap('8px'),
    cursor: 'pointer',
  },
  userName: {
    color: tokens.colorNeutralForeground1,
    fontWeight: tokens.fontWeightSemibold,
    '@media (max-width: 768px)': {
      display: 'none',
    },
  },
});

interface TopBarProps {
  showBreadcrumbs?: boolean;
}

export const TopBar: React.FC<TopBarProps> = ({ showBreadcrumbs = true }) => {
  const styles = useStyles();
  const { t } = useTranslation();
  const location = useLocation();
  const navigate = useNavigate();
  const { mode, toggleMode } = useTheme();
  const user = useAuthStore((state) => state.user);
  const toggleSidebar = useSidebarStore((state) => state.toggleCollapsed);

  const handleLogout = async () => {
    await authService.logout();
    // Explicit redirect to login page for immediate feedback
    window.location.href = '/login';
  };

  const handleLanguageChange = () => {
    const newLang = i18n.language === 'ru' ? 'en' : 'ru';
    i18n.changeLanguage(newLang);
  };

  // Find current page info for breadcrumbs
  const currentNavItem = ALL_NAV_ITEMS.find((item) => item.path === location.pathname);

  const getBreadcrumbs = () => {
    if (location.pathname === '/') {
      return [{ label: t('layout.nav.dashboard'), path: '/' }];
    }

    if (currentNavItem) {
      return [
        { label: t('layout.nav.dashboard'), path: '/' },
        { label: t(currentNavItem.labelKey), path: currentNavItem.path },
      ];
    }

    return [{ label: t('layout.nav.dashboard'), path: '/' }];
  };

  const breadcrumbs = getBreadcrumbs();

  return (
    <header className={styles.topbar}>
      <div className={styles.left}>
        <button className={styles.hamburgerMobile} onClick={toggleSidebar} aria-label="Toggle sidebar">
          <Navigation20Regular />
        </button>
        {showBreadcrumbs && (
          <Breadcrumb className={styles.breadcrumb}>
            {breadcrumbs.map((crumb, index) => (
              <BreadcrumbItem key={crumb.path}>
                {index === breadcrumbs.length - 1 ? (
                  <BreadcrumbButton current>{crumb.label}</BreadcrumbButton>
                ) : (
                  <BreadcrumbButton onClick={() => navigate(crumb.path)}>{crumb.label}</BreadcrumbButton>
                )}
              </BreadcrumbItem>
            ))}
          </Breadcrumb>
        )}
      </div>

      <div className={styles.right}>
        {/* Theme toggle */}
        <div className={styles.themeSwitch}>
          {mode === 'light' ? <WeatherSunny20Regular /> : <WeatherMoon20Regular />}
          <Switch checked={mode === 'dark'} onChange={toggleMode} aria-label={t('layout.theme')} />
        </div>

        {/* Language toggle */}
        <Button
          appearance="subtle"
          size="small"
          className={styles.langButton}
          onClick={handleLanguageChange}
          aria-label={t('common.language')}
        >
          {i18n.language.toUpperCase()}
        </Button>

        {/* User menu */}
        <Menu>
          <MenuTrigger>
            <Button
              appearance="transparent"
              className={styles.userMenu}
              aria-label={`${t('layout.profile')}: ${user?.full_name || user?.fio || user?.username}`}
            >
              <span className={styles.userName}>{user?.full_name || user?.fio || user?.username}</span>
              <Avatar
                name={user?.full_name || user?.fio || user?.username}
                size={32}
                color="brand"
              />
            </Button>
          </MenuTrigger>
          <MenuPopover>
            <MenuList>
              <MenuItem onClick={handleLogout} icon={<SignOut20Regular />}>
                {t('layout.logout')}
              </MenuItem>
            </MenuList>
          </MenuPopover>
        </Menu>
      </div>
    </header>
  );
};
