# Contributing to Вокзал.ТЕХ

Спасибо за интерес к проекту! Эта инструкция поможет вам начать контрибьютить.

## 📋 Содержание

- [Процесс разработки](#процесс-разработки)
- [Стандарты кода](#стандарты-кода)
- [Тестирование](#тестирование)
- [Отправка изменений](#отправка-изменений)
- [CI/CD](#cicd)

## 🔄 Процесс разработки

### 1. Fork и Clone

```bash
# Fork репозиторий через GitHub UI
git clone https://github.com/YOUR_USERNAME/vokzal.git
cd vokzal
git remote add upstream https://github.com/vokzal-tech/vokzal.git
```

### 2. Создайте ветку

```bash
git checkout -b feature/my-new-feature
# или
git checkout -b fix/bug-description
```

**Именование веток:**
- `feature/*` — новая функциональность
- `fix/*` — исправление багов
- `docs/*` — изменения в документации
- `refactor/*` — рефакторинг без изменения функциональности
- `test/*` — добавление/изменение тестов

### 3. Локальная разработка

```bash
# Запустите инфраструктуру
docker-compose -f infra/docker/docker-compose.yml up -d

# Go сервисы
cd services/auth
go run cmd/main.go

# React приложения
cd ui/admin-panel
npm install
npm run dev
```

### 4. Запустите тесты

```bash
# Go тесты
cd services/auth
go test ./...

# React тесты
cd ui/admin-panel
npm test
```

## 📝 Стандарты кода

### Go

Следуйте официальному [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).

**Основные правила:**
- Используйте `gofmt` для форматирования
- Запускайте `golangci-lint` перед коммитом
- Покрытие тестами: минимум 70%
- Все публичные функции должны иметь комментарии
- Используйте контекст для всех операций I/O

**Пример:**

```go
// CreateTicket создаёт новый билет для указанного рейса
func (s *TicketService) CreateTicket(ctx context.Context, req *CreateTicketRequest) (*Ticket, error) {
    if err := s.validate(req); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }
    
    ticket, err := s.repo.Create(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("failed to create ticket: %w", err)
    }
    
    return ticket, nil
}
```

### TypeScript/React

Следуйте [Airbnb JavaScript Style Guide](https://github.com/airbnb/javascript).

**Основные правила:**
- Используйте TypeScript strict mode
- Функциональные компоненты + hooks
- Именуйте компоненты в PascalCase
- Props интерфейсы должны заканчиваться на `Props`
- Используйте React Query для server state

**Пример:**

```typescript
interface TicketListProps {
  stationId: string;
  onSelect?: (ticket: Ticket) => void;
}

export const TicketList: React.FC<TicketListProps> = ({ stationId, onSelect }) => {
  const { data, isLoading } = useQuery({
    queryKey: ['tickets', stationId],
    queryFn: () => fetchTickets(stationId),
  });

  if (isLoading) return <Spinner />;

  return (
    <div className="ticket-list">
      {data?.map((ticket) => (
        <TicketCard key={ticket.id} ticket={ticket} onClick={onSelect} />
      ))}
    </div>
  );
};
```

### SQL

- Используйте snake_case для таблиц и колонок
- Всегда указывайте explicit типы
- Добавляйте индексы для foreign keys
- Включайте `up` и `down` миграции

### Commit Messages

Следуйте [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Типы:**
- `feat`: новая функциональность
- `fix`: исправление бага
- `docs`: изменения в документации
- `style`: форматирование, пробелы (не меняет код)
- `refactor`: рефакторинг кода
- `test`: добавление тестов
- `chore`: изменения в build процессе, зависимостях

**Примеры:**

```
feat(ticket): add ticket refund functionality

Implemented ticket refund logic with validation
and fiscal service integration.

Closes #123
```

```
fix(auth): prevent token expiration during active session

Added token refresh mechanism that triggers 5 minutes
before expiration.

Fixes #456
```

## 🧪 Тестирование

### Unit тесты

**Go:**
```bash
go test -v ./...
go test -race ./...
go test -coverprofile=coverage.out ./...
```

**React:**
```bash
npm test
npm test -- --coverage
```

### E2E тесты

```bash
cd tests/e2e
npm install
npm run test
```

### Линтеры

```bash
# Go
golangci-lint run

# TypeScript
npm run lint

# Все сразу через pre-commit
pre-commit run --all-files
```

## 📤 Отправка изменений

### 1. Pre-commit проверки

Установите pre-commit hooks:

```bash
pip install pre-commit
pre-commit install
```

Hooks автоматически запустятся при коммите.

### 2. Push в ваш fork

```bash
git push origin feature/my-new-feature
```

### 3. Создайте Pull Request

1. Перейдите на GitHub
2. Нажмите "New Pull Request"
3. Заполните шаблон:
   - Описание изменений
   - Связанные issues
   - Скриншоты (для UI)
   - Чеклист тестирования

### 4. Code Review

- Ответьте на комментарии reviewers
- Внесите изменения при необходимости
- Убедитесь, что все CI checks прошли

## 🚀 CI/CD

### Автоматические проверки

При каждом Pull Request запускаются:

1. **Lint & Type Check**
   - golangci-lint для Go
   - ESLint + TypeScript для UI
   
2. **Unit Tests**
   - Go tests с coverage
   - Jest tests с coverage
   
3. **Build**
   - Docker образы для сервисов
   - Production build для UI
   
4. **Security Scan**
   - Trivy vulnerability scanning
   - CodeQL analysis

### Требования для merge

- ✅ Все CI checks прошли
- ✅ Code review одобрен (минимум 1 reviewer)
- ✅ Coverage не упал ниже 70%
- ✅ Нет конфликтов с main
- ✅ Commit messages соответствуют стандарту

## 🐛 Reporting Bugs

Используйте GitHub Issues со следующей информацией:

```markdown
**Описание бага**
Краткое описание проблемы.

**Шаги для воспроизведения**
1. Перейти в '...'
2. Нажать на '...'
3. Увидеть ошибку

**Ожидаемое поведение**
Что должно было произойти.

**Скриншоты**
Приложите скриншоты если возможно.

**Окружение:**
- OS: [e.g. Windows 11]
- Browser: [e.g. Chrome 120]
- Version: [e.g. v1.2.3]

**Дополнительный контекст**
Любая другая полезная информация.
```

## 💡 Предложение новой функциональности

Используйте GitHub Issues:

```markdown
**Описание функциональности**
Чёткое описание того, что вы хотите добавить.

**Зачем это нужно?**
Объясните проблему, которую это решает.

**Предлагаемое решение**
Как вы видите реализацию?

**Альтернативы**
Какие другие решения вы рассматривали?

**Дополнительный контекст**
Mockups, примеры из других систем и т.д.
```

## 📚 Дополнительные ресурсы

- [Документация проекта](./docs/)
- [CI/CD Pipeline](./.github/CI_CD.md)
- [API Documentation](https://docs.vokzal.tech)
- [Архитектурные решения](./docs/architecture/)

## 🙏 Благодарности

Спасибо за вклад в Вокзал.ТЕХ! 

Список всех контрибьюторов: [CONTRIBUTORS.md](./CONTRIBUTORS.md)

## 📞 Вопросы?

- Telegram: @vokzal_tech
- Email: dev@vokzal.tech
- GitHub Discussions: [Обсуждения](https://github.com/vokzal-tech/vokzal/discussions)

---

© 2026 Вокзал.ТЕХ
