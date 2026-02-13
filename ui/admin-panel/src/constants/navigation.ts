import {
  Home20Regular,
  CalendarLtr20Regular,
  Map20Regular,
  VehicleBus20Regular,
  VehicleCar20Regular,
  PersonAccounts20Regular,
  Location20Regular,
  DocumentTable20Regular,
  DataArea20Regular,
  People20Regular,
  DocumentBulletList20Regular,
} from '@fluentui/react-icons';
import type { FluentIcon } from '@fluentui/react-icons';

export interface NavItem {
  key: string;
  label: string;
  labelKey: string; // i18n key
  icon: FluentIcon;
  path: string;
  roles?: string[]; // undefined = все роли
}

export interface NavGroup {
  key: string;
  label: string;
  labelKey: string; // i18n key
  items: NavItem[];
}

// Конфигурация навигационного меню с группировкой
export const NAVIGATION_GROUPS: NavGroup[] = [
  {
    key: 'overview',
    label: 'Обзор',
    labelKey: 'layout.groups.overview',
    items: [
      {
        key: 'dashboard',
        label: 'Главная',
        labelKey: 'layout.nav.dashboard',
        icon: Home20Regular,
        path: '/',
      },
    ],
  },
  {
    key: 'transport',
    label: 'Управление транспортом',
    labelKey: 'layout.groups.transport',
    items: [
      {
        key: 'schedules',
        label: 'Расписания',
        labelKey: 'nav.schedules',
        icon: CalendarLtr20Regular,
        path: '/schedules',
        roles: ['admin', 'dispatcher'],
      },
      {
        key: 'routes',
        label: 'Маршруты',
        labelKey: 'nav.routes',
        icon: Map20Regular,
        path: '/routes',
        roles: ['admin', 'dispatcher'],
      },
      {
        key: 'trips',
        label: 'Рейсы',
        labelKey: 'nav.trips',
        icon: VehicleCar20Regular,
        path: '/trips',
        roles: ['admin', 'dispatcher'],
      },
      {
        key: 'buses',
        label: 'Автобусы',
        labelKey: 'nav.buses',
        icon: VehicleBus20Regular,
        path: '/buses',
        roles: ['admin', 'dispatcher'],
      },
      {
        key: 'drivers',
        label: 'Водители',
        labelKey: 'nav.drivers',
        icon: PersonAccounts20Regular,
        path: '/drivers',
        roles: ['admin', 'dispatcher'],
      },
      {
        key: 'stations',
        label: 'Станции',
        labelKey: 'nav.stations',
        icon: Location20Regular,
        path: '/stations',
        roles: ['admin', 'dispatcher'],
      },
    ],
  },
  {
    key: 'reports',
    label: 'Отчётность',
    labelKey: 'layout.groups.reports',
    items: [
      {
        key: 'reports',
        label: 'Отчёты',
        labelKey: 'nav.reports',
        icon: DocumentTable20Regular,
        path: '/reports',
        roles: ['admin', 'dispatcher'],
      },
      {
        key: 'monitoring',
        label: 'Мониторинг',
        labelKey: 'nav.monitoring',
        icon: DataArea20Regular,
        path: '/monitoring',
        roles: ['admin'],
      },
    ],
  },
  {
    key: 'admin',
    label: 'Администрирование',
    labelKey: 'layout.groups.admin',
    items: [
      {
        key: 'users',
        label: 'Пользователи',
        labelKey: 'nav.users',
        icon: People20Regular,
        path: '/users',
        roles: ['admin'],
      },
      {
        key: 'audit',
        label: 'Журнал аудита',
        labelKey: 'nav.audit',
        icon: DocumentBulletList20Regular,
        path: '/audit',
        roles: ['admin'],
      },
    ],
  },
];

// Flat list для быстрого поиска
export const ALL_NAV_ITEMS: NavItem[] = NAVIGATION_GROUPS.flatMap((group) => group.items);

// Функция для проверки доступа к навигационному пункту
export const hasAccess = (item: NavItem, userRole?: string): boolean => {
  if (!item.roles || item.roles.length === 0) {
    return true; // доступно всем
  }
  if (!userRole) {
    return false;
  }
  return item.roles.includes(userRole);
};

// Функция для фильтрации групп по правам доступа
export const filterNavigationByRole = (userRole?: string): NavGroup[] => {
  return NAVIGATION_GROUPS.map((group) => ({
    ...group,
    items: group.items.filter((item) => hasAccess(item, userRole)),
  })).filter((group) => group.items.length > 0); // убираем пустые группы
};
