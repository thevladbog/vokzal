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

export const playSuccessSound = () => {
  const audioContext = new (window.AudioContext || (window as any).webkitAudioContext)();
  const oscillator = audioContext.createOscillator();
  const gainNode = audioContext.createGain();

  oscillator.connect(gainNode);
  gainNode.connect(audioContext.destination);

  oscillator.frequency.value = 800;
  oscillator.type = 'sine';

  gainNode.gain.setValueAtTime(0.3, audioContext.currentTime);
  gainNode.gain.exponentialRampToValueAtTime(0.01, audioContext.currentTime + 0.2);

  oscillator.start(audioContext.currentTime);
  oscillator.stop(audioContext.currentTime + 0.2);
};

export const playErrorSound = () => {
  const audioContext = new (window.AudioContext || (window as any).webkitAudioContext)();
  const oscillator = audioContext.createOscillator();
  const gainNode = audioContext.createGain();

  oscillator.connect(gainNode);
  gainNode.connect(audioContext.destination);

  oscillator.frequency.value = 400;
  oscillator.type = 'sawtooth';

  gainNode.gain.setValueAtTime(0.3, audioContext.currentTime);
  gainNode.gain.exponentialRampToValueAtTime(0.01, audioContext.currentTime + 0.3);

  oscillator.start(audioContext.currentTime);
  oscillator.stop(audioContext.currentTime + 0.3);
};

export const vibratePhone = (pattern: number | number[] = 200) => {
  if ('vibrate' in navigator) {
    navigator.vibrate(pattern);
  }
};
