import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { makeStyles, tokens, Text, shorthands } from '@fluentui/react-components';
import { Navigation20Regular } from '@fluentui/react-icons';
import { useSidebarStore } from '@/stores/sidebarStore';
import { useAuthStore } from '@/stores/authStore';
import { filterNavigationByRole } from '@/constants/navigation';

const SIDEBAR_WIDTH = 240;
const SIDEBAR_COLLAPSED_WIDTH = 64;

const useStyles = makeStyles({
  sidebar: {
    height: '100vh',
    backgroundColor: tokens.colorNeutralBackground3,
    ...shorthands.borderRight('1px', 'solid', tokens.colorNeutralStroke2),
    display: 'flex',
    flexDirection: 'column',
    ...shorthands.transition('width', '200ms', 'ease-in-out'),
    position: 'fixed',
    left: 0,
    top: 0,
    zIndex: 1000,
    overflowX: 'hidden',
    '@media (max-width: 768px)': {
      transform: 'translateX(-100%)',
    },
  },
  sidebarExpanded: {
    width: `${SIDEBAR_WIDTH}px`,
    '@media (max-width: 768px)': {
      transform: 'translateX(0)',
      boxShadow: tokens.shadow16,
    },
  },
  sidebarCollapsed: {
    width: `${SIDEBAR_COLLAPSED_WIDTH}px`,
    '@media (max-width: 768px)': {
      transform: 'translateX(-100%)',
    },
  },
  overlay: {
    display: 'none',
    '@media (max-width: 768px)': {
      display: 'block',
      position: 'fixed',
      top: 0,
      left: 0,
      right: 0,
      bottom: 0,
      backgroundColor: 'rgba(0, 0, 0, 0.5)',
      zIndex: 999,
    },
  },
  header: {
    ...shorthands.padding('16px'),
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    minHeight: '64px',
    ...shorthands.borderBottom('1px', 'solid', tokens.colorNeutralStroke2),
  },
  logo: {
    display: 'flex',
    alignItems: 'center',
    ...shorthands.gap('12px'),
    ...shorthands.overflow('hidden'),
  },
  logoText: {
    fontSize: tokens.fontSizeBase400,
    fontWeight: tokens.fontWeightSemibold,
    color: tokens.colorBrandForeground1,
    whiteSpace: 'nowrap',
  },
  hamburger: {
    backgroundColor: 'transparent',
    ...shorthands.border('none'),
    ...shorthands.padding('8px'),
    cursor: 'pointer',
    ...shorthands.borderRadius(tokens.borderRadiusMedium),
    display: 'flex',
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
  nav: {
    ...shorthands.flex(1),
    ...shorthands.padding('16px', '0'),
    overflowY: 'auto',
    overflowX: 'hidden',
  },
  navGroup: {
    marginBottom: '24px',
  },
  navGroupTitle: {
    ...shorthands.padding('8px', '16px'),
    fontSize: tokens.fontSizeBase200,
    fontWeight: tokens.fontWeightSemibold,
    color: tokens.colorNeutralForeground3,
    textTransform: 'uppercase',
    letterSpacing: '0.5px',
    whiteSpace: 'nowrap',
    ...shorthands.overflow('hidden'),
    textOverflow: 'ellipsis',
  },
  navGroupTitleCollapsed: {
    textAlign: 'center',
    ...shorthands.padding('8px'),
  },
  navList: {
    listStyle: 'none',
    ...shorthands.padding('0'),
    ...shorthands.margin('0'),
  },
  navItem: {
    ...shorthands.margin('2px', '8px'),
  },
  navLink: {
    display: 'flex',
    alignItems: 'center',
    ...shorthands.gap('12px'),
    ...shorthands.padding('10px', '12px'),
    ...shorthands.borderRadius(tokens.borderRadiusMedium),
    textDecoration: 'none',
    color: tokens.colorNeutralForeground2,
    ...shorthands.transition('all', '150ms', 'ease-in-out'),
    ':hover': {
      backgroundColor: tokens.colorNeutralBackground2Hover,
      color: tokens.colorNeutralForeground1,
      transform: 'translateX(2px)',
    },
  },
  navLinkCollapsed: {
    justifyContent: 'center',
    ...shorthands.padding('10px'),
  },
  navLinkActive: {
    backgroundColor: tokens.colorBrandBackground2,
    color: tokens.colorBrandForeground1,
    fontWeight: tokens.fontWeightSemibold,
    ':hover': {
      backgroundColor: tokens.colorBrandBackground2Hover,
      color: tokens.colorBrandForeground1,
    },
  },
  navIcon: {
    fontSize: '20px',
    flexShrink: 0,
  },
  navLabel: {
    whiteSpace: 'nowrap',
    ...shorthands.overflow('hidden'),
    textOverflow: 'ellipsis',
  },
});

export const Sidebar: React.FC = () => {
  const styles = useStyles();
  const { t } = useTranslation();
  const location = useLocation();
  const isCollapsed = useSidebarStore((state) => state.isCollapsed);
  const toggleCollapsed = useSidebarStore((state) => state.toggleCollapsed);
  const user = useAuthStore((state) => state.user);

  const navigationGroups = filterNavigationByRole(user?.role);

  // Close sidebar on mobile when clicking overlay
  const handleOverlayClick = () => {
    if (window.innerWidth <= 768 && !isCollapsed) {
      toggleCollapsed();
    }
  };

  return (
    <>
      {!isCollapsed && <div className={styles.overlay} onClick={handleOverlayClick} />}
      <aside className={`${styles.sidebar} ${isCollapsed ? styles.sidebarCollapsed : styles.sidebarExpanded}`}>
        <div className={styles.header}>
          {!isCollapsed && (
            <div className={styles.logo}>
              <Text className={styles.logoText}>Вокзал.ТЕХ</Text>
            </div>
          )}
          <button className={styles.hamburger} onClick={toggleCollapsed} aria-label="Toggle sidebar">
            <Navigation20Regular />
          </button>
        </div>

        <nav className={styles.nav}>
          {navigationGroups.map((group) => (
            <div key={group.key} className={styles.navGroup}>
              <div className={`${styles.navGroupTitle} ${isCollapsed ? styles.navGroupTitleCollapsed : ''}`}>
                {isCollapsed ? '•' : t(group.labelKey)}
              </div>
              <ul className={styles.navList}>
                {group.items.map((item) => {
                  const isActive = location.pathname === item.path;
                  const Icon = item.icon;

                  return (
                    <li key={item.key} className={styles.navItem}>
                      <Link
                        to={item.path}
                        className={`${styles.navLink} ${isCollapsed ? styles.navLinkCollapsed : ''} ${
                          isActive ? styles.navLinkActive : ''
                        }`}
                        title={isCollapsed ? t(item.labelKey) : undefined}
                        onClick={() => {
                          // Close sidebar on mobile after navigation
                          if (window.innerWidth <= 768) {
                            toggleCollapsed();
                          }
                        }}
                      >
                        <Icon className={styles.navIcon} />
                        {!isCollapsed && <span className={styles.navLabel}>{t(item.labelKey)}</span>}
                      </Link>
                    </li>
                  );
                })}
              </ul>
            </div>
          ))}
        </nav>
      </aside>
    </>
  );
};

export { SIDEBAR_WIDTH, SIDEBAR_COLLAPSED_WIDTH };
