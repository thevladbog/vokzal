import React from 'react';
import { DatePicker } from '@fluentui/react-datepicker-compat';
import { useTranslation } from 'react-i18next';
import type { DatePickerProps } from '@fluentui/react-datepicker-compat';

interface VokzalDatePickerProps extends Omit<DatePickerProps, 'strings'> {
  value?: Date | null;
  onSelectDate?: (date: Date | null | undefined) => void;
  required?: boolean;
  disabled?: boolean;
  placeholder?: string;
  minDate?: Date;
  maxDate?: Date;
}

/**
 * VokzalDatePicker - обёртка для DatePicker из Fluent UI с русской локализацией
 * 
 * Использование:
 * <Field label="Дата" required>
 *   <VokzalDatePicker
 *     value={date}
 *     onSelectDate={(date) => setDate(date)}
 *   />
 * </Field>
 */
export const VokzalDatePicker: React.FC<VokzalDatePickerProps> = ({
  value,
  onSelectDate,
  placeholder,
  required = false,
  disabled = false,
  minDate,
  maxDate,
  ...rest
}) => {
  const { i18n } = useTranslation();

  // Определяем локализацию в зависимости от текущего языка
  const isRussian = i18n.language === 'ru';

  const strings = {
    months: isRussian
      ? ['Январь', 'Февраль', 'Март', 'Апрель', 'Май', 'Июнь', 'Июль', 'Август', 'Сентябрь', 'Октябрь', 'Ноябрь', 'Декабрь']
      : ['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December'],
    shortMonths: isRussian
      ? ['Янв', 'Фев', 'Мар', 'Апр', 'Май', 'Июн', 'Июл', 'Авг', 'Сен', 'Окт', 'Ноя', 'Дек']
      : ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'],
    days: isRussian
      ? ['Воскресенье', 'Понедельник', 'Вторник', 'Среда', 'Четверг', 'Пятница', 'Суббота']
      : ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'],
    shortDays: isRussian
      ? ['Вс', 'Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб']
      : ['S', 'M', 'T', 'W', 'T', 'F', 'S'],
    goToToday: isRussian ? 'Сегодня' : 'Go to today',
    weekNumberFormatString: isRussian ? 'Неделя {0}' : 'Week number {0}',
    prevMonthAriaLabel: isRussian ? 'Предыдущий месяц' : 'Go to previous month',
    nextMonthAriaLabel: isRussian ? 'Следующий месяц' : 'Go to next month',
    prevYearAriaLabel: isRussian ? 'Предыдущий год' : 'Go to previous year',
    nextYearAriaLabel: isRussian ? 'Следующий год' : 'Go to next year',
    closeButtonAriaLabel: isRussian ? 'Закрыть' : 'Close date picker',
    monthPickerHeaderAriaLabel: isRussian ? '{0}, выберите для смены года' : '{0}, select to change the year',
    yearPickerHeaderAriaLabel: isRussian ? '{0}, выберите для смены месяца' : '{0}, select to change the month',
  };

  // Форматирование даты в зависимости от локали
  const formatDate = (date?: Date): string => {
    if (!date) return '';
    
    if (isRussian) {
      const day = date.getDate().toString().padStart(2, '0');
      const month = (date.getMonth() + 1).toString().padStart(2, '0');
      const year = date.getFullYear();
      return `${day}.${month}.${year}`;
    }
    
    // Английский формат
    const month = (date.getMonth() + 1).toString().padStart(2, '0');
    const day = date.getDate().toString().padStart(2, '0');
    const year = date.getFullYear();
    return `${month}/${day}/${year}`;
  };

  // Парсинг даты из строки
  const parseDateFromString = (dateStr: string): Date | null => {
    if (!dateStr) return null;

    // Пробуем распарсить российский формат DD.MM.YYYY
    const ruMatch = dateStr.match(/^(\d{1,2})\.(\d{1,2})\.(\d{4})$/);
    if (ruMatch) {
      const [, day, month, year] = ruMatch;
      const parsedDay = parseInt(day);
      const parsedMonth = parseInt(month);
      const parsedYear = parseInt(year);
      
      const date = new Date(parsedYear, parsedMonth - 1, parsedDay);
      
      // Валидация: проверяем, что Date не нормализовал невалидные значения
      if (
        !isNaN(date.getTime()) &&
        date.getFullYear() === parsedYear &&
        date.getMonth() === parsedMonth - 1 &&
        date.getDate() === parsedDay
      ) {
        return date;
      }
    }

    // Пробуем распарсить американский формат MM/DD/YYYY
    const usMatch = dateStr.match(/^(\d{1,2})\/(\d{1,2})\/(\d{4})$/);
    if (usMatch) {
      const [, month, day, year] = usMatch;
      const parsedDay = parseInt(day);
      const parsedMonth = parseInt(month);
      const parsedYear = parseInt(year);
      
      const date = new Date(parsedYear, parsedMonth - 1, parsedDay);
      
      // Валидация: проверяем, что Date не нормализовал невалидные значения
      if (
        !isNaN(date.getTime()) &&
        date.getFullYear() === parsedYear &&
        date.getMonth() === parsedMonth - 1 &&
        date.getDate() === parsedDay
      ) {
        return date;
      }
    }

    // Пробуем стандартный ISO формат YYYY-MM-DD
    const isoMatch = dateStr.match(/^(\d{4})-(\d{1,2})-(\d{1,2})$/);
    if (isoMatch) {
      const [, year, month, day] = isoMatch;
      const parsedDay = parseInt(day);
      const parsedMonth = parseInt(month);
      const parsedYear = parseInt(year);
      
      const date = new Date(parsedYear, parsedMonth - 1, parsedDay);
      
      // Валидация: проверяем, что Date не нормализовал невалидные значения
      if (
        !isNaN(date.getTime()) &&
        date.getFullYear() === parsedYear &&
        date.getMonth() === parsedMonth - 1 &&
        date.getDate() === parsedDay
      ) {
        return date;
      }
    }

    return null;
  };

  return (
    <DatePicker
      placeholder={placeholder || (isRussian ? 'Выберите дату' : 'Select a date')}
      value={value ?? null}
      onSelectDate={onSelectDate}
      strings={strings}
      formatDate={formatDate}
      parseDateFromString={parseDateFromString}
      firstDayOfWeek={isRussian ? 1 : 0} // Понедельник для RU, Воскресенье для EN
      showGoToToday
      allowTextInput
      required={required}
      disabled={disabled}
      minDate={minDate}
      maxDate={maxDate}
      {...rest}
    />
  );
};
