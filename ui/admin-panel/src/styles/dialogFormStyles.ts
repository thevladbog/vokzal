import { makeStyles, tokens } from '@fluentui/react-components';

/**
 * Стили для единообразного оформления форм в диалогах
 */
export const useDialogFormStyles = makeStyles({
  formContainer: {
    display: 'flex',
    flexDirection: 'column',
    gap: '20px',
  },
  formRow: {
    display: 'flex',
    flexDirection: 'column',
    gap: '8px',
  },
  formRowInline: {
    display: 'grid',
    gridTemplateColumns: '140px 1fr',
    gap: '12px',
    alignItems: 'center',
  },
  formLabel: {
    fontWeight: tokens.fontWeightSemibold,
    fontSize: tokens.fontSizeBase300,
  },
  checkboxGroup: {
    display: 'flex',
    flexWrap: 'wrap',
    gap: '12px',
    marginTop: '8px',
  },
  hint: {
    fontSize: tokens.fontSizeBase200,
    color: tokens.colorNeutralForeground3,
    marginTop: '4px',
  },
  error: {
    fontSize: tokens.fontSizeBase200,
    color: tokens.colorPaletteRedForeground1,
    marginTop: '4px',
  },
  formActions: {
    display: 'flex',
    justifyContent: 'flex-end',
    gap: '12px',
    marginTop: '8px',
  },
});
