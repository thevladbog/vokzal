import { format, parseISO } from 'date-fns';
import { ru } from 'date-fns/locale';

export const formatDate = (dateString: string): string => {
  try {
    return format(parseISO(dateString), 'd MMMM yyyy', { locale: ru });
  } catch {
    return dateString;
  }
};

export const formatTime = (dateString: string): string => {
  try {
    return format(parseISO(dateString), 'HH:mm', { locale: ru });
  } catch {
    return dateString;
  }
};

export const formatDateTime = (dateString: string): string => {
  try {
    return format(parseISO(dateString), 'd MMMM yyyy, HH:mm', { locale: ru });
  } catch {
    return dateString;
  }
};

export const formatDuration = (minutes: number): string => {
  const hours = Math.floor(minutes / 60);
  const mins = minutes % 60;
  
  if (hours === 0) {
    return `${mins} мин`;
  }
  
  if (mins === 0) {
    return `${hours} ч`;
  }
  
  return `${hours} ч ${mins} мин`;
};

export const formatPrice = (amount: number): string => {
  return new Intl.NumberFormat('ru-RU', {
    style: 'currency',
    currency: 'RUB',
    minimumFractionDigits: 0,
    maximumFractionDigits: 2,
  }).format(amount);
};
