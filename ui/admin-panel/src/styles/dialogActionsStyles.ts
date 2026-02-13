import { makeStyles, shorthands } from '@fluentui/react-components';

/**
 * Стили для обёртки кнопок в DialogActions с правильным расстоянием
 */
export const useDialogActionsStyles = makeStyles({
  wrapper: {
    display: 'flex',
    ...shorthands.gap('12px'),
    justifyContent: 'flex-end',
    width: '100%',
    marginTop: '20px', // Отступ сверху от контента
  },
});
