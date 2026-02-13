import React, { ReactNode } from 'react';
import { makeStyles, tokens, shorthands } from '@fluentui/react-components';
import { Sidebar, SIDEBAR_WIDTH, SIDEBAR_COLLAPSED_WIDTH } from './Sidebar';
import { TopBar } from './TopBar';
import { useSidebarStore } from '@/stores/sidebarStore';

const useStyles = makeStyles({
  layout: {
    display: 'flex',
    minHeight: '100vh',
    backgroundColor: tokens.colorNeutralBackground2,
  },
  mainWrapper: {
    ...shorthands.flex(1),
    display: 'flex',
    flexDirection: 'column',
    ...shorthands.transition('margin-left', '200ms', 'ease-in-out'),
    marginLeft: `${SIDEBAR_WIDTH}px`,
    '@media (max-width: 768px)': {
      marginLeft: '0',
    },
  },
  mainWrapperCollapsed: {
    marginLeft: `${SIDEBAR_COLLAPSED_WIDTH}px`,
    '@media (max-width: 768px)': {
      marginLeft: '0',
    },
  },
  content: {
    ...shorthands.flex(1),
    ...shorthands.padding('24px'),
    maxWidth: '100%',
    ...shorthands.overflow('auto'),
    '@media (max-width: 768px)': {
      ...shorthands.padding('16px'),
    },
  },
  contentInner: {
    maxWidth: '1600px',
    marginLeft: 'auto',
    marginRight: 'auto',
    width: '100%',
  },
});

interface AppLayoutProps {
  children: ReactNode;
  showBreadcrumbs?: boolean;
}

export const AppLayout: React.FC<AppLayoutProps> = ({ children, showBreadcrumbs = true }) => {
  const styles = useStyles();
  const isCollapsed = useSidebarStore((state) => state.isCollapsed);

  return (
    <div className={styles.layout}>
      <Sidebar />
      <div className={`${styles.mainWrapper} ${isCollapsed ? styles.mainWrapperCollapsed : ''}`}>
        <TopBar showBreadcrumbs={showBreadcrumbs} />
        <main className={styles.content}>
          <div className={styles.contentInner}>{children}</div>
        </main>
      </div>
    </div>
  );
};
