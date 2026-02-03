/// <reference types="cypress" />

describe('Продажа билета', () => {
  beforeEach(() => {
    // Авторизуемся как кассир
    cy.loginAsCashier();
    cy.visit('/pos');
  });

  it('полный цикл продажи билета', () => {
    // Шаг 1: Поиск рейса
    cy.get('[data-cy=from-station-select]').click();
    cy.contains('Ростов-на-Дону').click();

    cy.get('[data-cy=to-station-select]').click();
    cy.contains('Москва').click();

    cy.get('[data-cy=date-picker]').click();
    // Выбираем завтрашний день
    cy.get('.calendar-day').contains(new Date().getDate() + 1).click();

    cy.get('[data-cy=search-button]').click();

    // Проверяем, что рейсы загрузились
    cy.get('[data-cy=trip-list]').should('be.visible');
    cy.get('[data-cy=trip-card]').should('have.length.greaterThan', 0);

    // Шаг 2: Выбор рейса
    cy.get('[data-cy=trip-card]').first().click();
    cy.get('[data-cy=select-trip-button]').click();

    // Шаг 3: Выбор места
    cy.get('[data-cy=seat-map]').should('be.visible');
    cy.get('[data-cy=seat-available]').first().click();

    // Проверяем, что место выбрано
    cy.get('[data-cy=selected-seat]').should('be.visible');
    cy.get('[data-cy=continue-button]').should('be.enabled').click();

    // Шаг 4: Ввод данных пассажира
    cy.get('[data-cy=passenger-lastname]').type('Иванов');
    cy.get('[data-cy=passenger-firstname]').type('Иван');
    cy.get('[data-cy=passenger-middlename]').type('Иванович');

    cy.get('[data-cy=document-type-select]').click();
    cy.contains('Паспорт РФ').click();

    cy.get('[data-cy=document-series]').type('4510');
    cy.get('[data-cy=document-number]').type('123456');

    cy.get('[data-cy=passenger-phone]').type('+79991234567');
    cy.get('[data-cy=passenger-email]').type('ivanov@example.com');

    cy.get('[data-cy=continue-button]').click();

    // Шаг 5: Подтверждение и оплата
    cy.get('[data-cy=ticket-summary]').should('be.visible');
    cy.get('[data-cy=total-price]').should('contain', '₽');

    // Выбираем способ оплаты
    cy.get('[data-cy=payment-method-cash]').click();

    // Вводим сумму наличными
    cy.get('[data-cy=cash-amount]').clear().type('2000');
    cy.get('[data-cy=calculate-change]').click();

    // Проверяем расчёт сдачи
    cy.get('[data-cy=change-amount]').should('be.visible');

    // Подтверждаем оплату
    cy.get('[data-cy=confirm-payment-button]').click();

    // Шаг 6: Фискализация
    cy.get('[data-cy=fiscal-progress]').should('be.visible');
    cy.contains(/фискализация/i).should('be.visible');

    // Ждём завершения фискализации (до 30 секунд)
    cy.get('[data-cy=fiscal-success]', { timeout: 30000 }).should('be.visible');

    // Шаг 7: Печать билета
    cy.get('[data-cy=print-ticket-button]').should('be.enabled').click();
    cy.get('[data-cy=print-progress]').should('be.visible');
    cy.get('[data-cy=print-success]', { timeout: 10000 }).should('be.visible');

    // Проверяем, что билет создан
    cy.get('[data-cy=ticket-number]').should('be.visible');
    cy.get('[data-cy=qr-code]').should('be.visible');

    // Проверяем, что можем вернуться к продаже
    cy.get('[data-cy=new-sale-button]').click();
    cy.get('[data-cy=from-station-select]').should('be.visible');
  });

  it('продажа билета с оплатой картой', () => {
    // Выбираем рейс (упрощённый сценарий)
    cy.get('[data-cy=quick-sale-button]').click();
    cy.get('[data-cy=trip-card]').first().click();
    cy.get('[data-cy=seat-available]').first().click();
    cy.get('[data-cy=continue-button]').click();

    // Минимальные данные пассажира
    cy.get('[data-cy=passenger-lastname]').type('Петров');
    cy.get('[data-cy=passenger-firstname]').type('Пётр');
    cy.get('[data-cy=document-series]').type('4511');
    cy.get('[data-cy=document-number]').type('654321');
    cy.get('[data-cy=continue-button]').click();

    // Оплата картой
    cy.get('[data-cy=payment-method-card]').click();
    cy.get('[data-cy=confirm-payment-button]').click();

    // Эмуляция оплаты через терминал
    cy.get('[data-cy=terminal-status]').should('be.visible');
    cy.contains(/ожидание оплаты/i).should('be.visible');

    // Имитируем успешную оплату (в реальности происходит через терминал)
    cy.get('[data-cy=payment-success]', { timeout: 60000 }).should('be.visible');

    // Проверяем фискализацию
    cy.get('[data-cy=fiscal-success]', { timeout: 30000 }).should('be.visible');
  });

  it('отмена продажи на этапе выбора места', () => {
    // Начинаем продажу
    cy.get('[data-cy=quick-sale-button]').click();
    cy.get('[data-cy=trip-card]').first().click();
    cy.get('[data-cy=seat-available]').first().click();

    // Отменяем
    cy.get('[data-cy=cancel-button]').click();
    cy.get('[data-cy=confirm-cancel-button]').click();

    // Проверяем, что вернулись к поиску
    cy.get('[data-cy=from-station-select]').should('be.visible');
  });

  it('валидация данных пассажира', () => {
    // Доходим до ввода данных
    cy.get('[data-cy=quick-sale-button]').click();
    cy.get('[data-cy=trip-card]').first().click();
    cy.get('[data-cy=seat-available]').first().click();
    cy.get('[data-cy=continue-button]').click();

    // Пытаемся продолжить с пустыми полями
    cy.get('[data-cy=continue-button]').click();

    // Проверяем валидацию
    cy.get('[data-cy=passenger-lastname]').should('have.attr', 'aria-invalid', 'true');
    cy.get('[data-cy=passenger-firstname]').should('have.attr', 'aria-invalid', 'true');
    cy.get('[data-cy=document-number]').should('have.attr', 'aria-invalid', 'true');

    // Заполняем корректно
    cy.get('[data-cy=passenger-lastname]').type('Сидоров');
    cy.get('[data-cy=passenger-firstname]').type('Сидор');
    cy.get('[data-cy=document-series]').type('4512');
    cy.get('[data-cy=document-number]').type('789012');

    // Проверяем, что валидация прошла
    cy.get('[data-cy=continue-button]').click();
    cy.get('[data-cy=ticket-summary]').should('be.visible');
  });

  it('недостаточно мест на рейсе', () => {
    // Создаём ситуацию, когда мест нет (через API)
    cy.apiRequest('POST', '/trips', {
      routeId: 'test-route',
      departureTime: new Date().toISOString(),
      totalSeats: 1,
      availableSeats: 0,
    });

    // Пытаемся найти этот рейс
    cy.get('[data-cy=from-station-select]').click();
    cy.contains('Ростов-на-Дону').click();
    cy.get('[data-cy=to-station-select]').click();
    cy.contains('Москва').click();
    cy.get('[data-cy=search-button]').click();

    // Проверяем, что рейс недоступен
    cy.contains(/мест нет/i).should('be.visible');
    cy.get('[data-cy=select-trip-button]').should('be.disabled');
  });

  it('экран покупателя отображает информацию', () => {
    // Проверяем, что второе окно (экран покупателя) открывается
    cy.window().then((win) => {
      cy.stub(win, 'open').as('windowOpen');
    });

    cy.get('[data-cy=open-customer-display]').click();
    cy.get('@windowOpen').should('be.called');

    // В реальной ситуации нужно тестировать второе окно отдельно
  });
});
