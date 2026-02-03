# Вокзал.ТЕХ - Brand Assets

Фирменные материалы для использования во всех UI приложениях платформы.

## 📁 Структура

```
ui/shared/assets/
├── logo.svg              # Основной логотип (SVG)
├── logo.png              # Основной логотип (PNG, 4000x4000)
├── brand-colors.css      # CSS переменные с цветами и токенами дизайна
└── README.md             # Этот файл
```

## 🎨 Логотип

### Использование в React/TypeScript

```tsx
import logo from '@/shared/assets/logo.svg';

function Header() {
  return (
    <img 
      src={logo} 
      alt="Вокзал.ТЕХ" 
      style={{ height: '40px' }}
    />
  );
}
```

### Минимальные размеры
- **Web**: 120px по ширине
- **Mobile**: 80px по ширине
- **Favicon**: 32x32px

## 🎨 Цвета

Импортируйте CSS файл в корневой компонент:

```tsx
import '@/shared/assets/brand-colors.css';
```

### Основные цвета

```css
/* Primary - Синий */
background-color: var(--color-primary);        /* #2563EB */
background-color: var(--color-primary-dark);   /* #1D4ED8 */

/* Secondary - Зелёный */
background-color: var(--color-secondary);      /* #10B981 */

/* Accent - Оранжевый */
background-color: var(--color-accent);         /* #F59E0B */
```

### Использование в styled-components

```tsx
import styled from 'styled-components';

const Button = styled.button`
  background-color: var(--color-primary);
  color: var(--color-white);
  padding: var(--space-3) var(--space-6);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-sm);
  
  &:hover {
    background-color: var(--color-primary-dark);
  }
`;
```

### Использование в Tailwind CSS

Если используете Tailwind, добавьте в `tailwind.config.js`:

```js
module.exports = {
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: '#2563EB',
          dark: '#1D4ED8',
          light: '#3B82F6',
        },
        secondary: {
          DEFAULT: '#10B981',
          dark: '#059669',
        },
        accent: {
          DEFAULT: '#F59E0B',
          dark: '#D97706',
        },
      },
    },
  },
};
```

## 🔧 Design Tokens

Все переменные доступны через CSS Custom Properties:

### Spacing
```css
margin: var(--space-4);      /* 16px */
padding: var(--space-6);     /* 24px */
gap: var(--space-2);         /* 8px */
```

### Typography
```css
font-family: var(--font-family-base);  /* Inter */
font-size: var(--font-size-base);      /* 16px */
font-weight: var(--font-weight-semibold); /* 600 */
line-height: var(--line-height-normal);   /* 1.5 */
```

### Border Radius
```css
border-radius: var(--radius-sm);   /* 4px */
border-radius: var(--radius-md);   /* 8px */
border-radius: var(--radius-lg);   /* 12px */
border-radius: var(--radius-full); /* 9999px */
```

### Shadows
```css
box-shadow: var(--shadow-sm);   /* Subtle */
box-shadow: var(--shadow-md);   /* Medium */
box-shadow: var(--shadow-lg);   /* Large */
```

## 📱 Адаптивность

### Breakpoints
```css
/* Mobile */
@media (max-width: 767px) {
  /* Уменьшенные размеры логотипа */
  /* Полноширинные кнопки */
}

/* Tablet */
@media (min-width: 768px) and (max-width: 1023px) {
  /* Средние размеры */
}

/* Desktop */
@media (min-width: 1024px) {
  /* Полные размеры */
}
```

## 🎯 Готовые компоненты

### Кнопка Primary
```tsx
<button className="btn-primary">
  Продолжить
</button>
```

### Карточка
```tsx
<div className="card">
  <h3>Заголовок</h3>
  <p>Содержимое карточки</p>
</div>
```

### Input
```tsx
<input 
  type="text" 
  className="input"
  placeholder="Введите текст"
/>
```

### Badge
```tsx
<span className="badge badge-success">
  Оплачено
</span>
```

## 📚 Дополнительная информация

Полное руководство по фирменному стилю:
👉 [Brand Guidelines](../../docs/branding/BRAND_GUIDELINES.md)

## 📧 Контакты

По вопросам использования брендинга:
- Email: brand@vokzal.tech
- Telegram: @vokzal_tech_brand
