/// <reference types="cypress" />

describe('Возврат билета', () => {
  let ticketId: string;
  let ticketNumber: string;

  beforeEach(() => {
    cy.loginAsCashier();
    cy.visit('/pos/refund');
  });

  // Создаём тестовый билет перед тестами
  before(() => {
    cy.loginAsCashier();
    
    // Создаём билет через API для возврата
    cy.apiRequest('POST', '/tickets', {
      tripId: 'test-trip-123',
      passengerId: 'test-passenger-123',
      seatNumber: 15,
      price: 1500,
      paymentMethod: 'cash',
    }).then((response) => {
      ticketId = response.body.data.id;
      ticketNumber = response.body.data.ticketNumber;
    });
  });

  it('полный цикл возврата билета', () => {
    // Шаг 1: Поиск билета
    cy.get('[data-cy=ticket-search-input]').type(ticketNumber);
    cy.get('[data-cy=search-ticket-button]').click();

    // Проверяем, что билет найден
    cy.get('[data-cy=ticket-details]').should('be.visible');
    cy.get('[data-cy=passenger-name]').should('not.be.empty');
    cy.get('[data-cy=route-info]').should('be.visible');

    // Шаг 2: Проверка возможности возврата
    cy.get('[data-cy=refund-available]').should('be.visible');
    cy.get('[data-cy=refund-amount]').should('contain', '₽');

    // Проверяем расчёт штрафа
    cy.get('[data-cy=penalty-amount]').should('be.visible');
    
    // Шаг 3: Подтверждение возврата
    cy.get('[data-cy=refund-reason-select]').click();
    cy.contains('По просьбе пассажира').click();
    
    cy.get('[data-cy=confirm-refund-button]').click();

    // Подтверждаем в диалоге
    cy.get('[data-cy=confirm-dialog]').should('be.visible');
    cy.get('[data-cy=confirm-yes-button]').click();

    // Шаг 4: Фискализация возврата
    cy.get('[data-cy=fiscal-progress]').should('be.visible');
    cy.contains(/фискализация возврата/i).should('be.visible');

    cy.get('[data-cy=fiscal-success]', { timeout: 30000 }).should('be.visible');

    // Шаг 5: Возврат средств
    cy.get('[data-cy=refund-method-cash]').should('be.visible');
    cy.get('[data-cy=complete-refund-button]').click();

    // Проверяем успешное завершение
    cy.get('[data-cy=refund-success]').should('be.visible');
    cy.get('[data-cy=refund-receipt]').should('be.visible');

    // Проверяем, что билет помечен как возвращённый
    cy.apiRequest('GET', `/tickets/${ticketId}`).then((response) => {
      expect(response.body.data.status).to.eq('refunded');
    });
  });

  it('возврат невозможен - менее 1 часа до отправления', () => {
    // Создаём билет на рейс, который отправляется через 30 минут
    const departureTime = new Date(Date.now() + 30 * 60 * 1000).toISOString();
    
    cy.apiRequest('POST', '/trips', {
      routeId: 'test-route',
      departureTime,
      totalSeats: 50,
    }).then((tripResponse) => {
      const tripId = tripResponse.body.data.id;

      cy.apiRequest('POST', '/tickets', {
        tripId,
        passengerId: 'test-passenger-456',
        seatNumber: 20,
        price: 2000,
      }).then((ticketResponse) => {
        const testTicketNumber = ticketResponse.body.data.ticketNumber;

        // Пытаемся вернуть билет
        cy.get('[data-cy=ticket-search-input]').type(testTicketNumber);
        cy.get('[data-cy=search-ticket-button]').click();

        // Проверяем, что возврат недоступен
        cy.contains(/возврат невозможен/i).should('be.visible');
        cy.contains(/менее 1 часа/i).should('be.visible');
        cy.get('[data-cy=confirm-refund-button]').should('be.disabled');
      });
    });
  });

  it('билет уже возвращён', () => {
    // Создаём и сразу возвращаем билет через API
    cy.apiRequest('POST', '/tickets', {
      tripId: 'test-trip-789',
      passengerId: 'test-passenger-789',
      seatNumber: 25,
      price: 1800,
    }).then((ticketResponse) => {
      const testTicketId = ticketResponse.body.data.id;
      const testTicketNumber = ticketResponse.body.data.ticketNumber;

      // Возвращаем через API
      cy.apiRequest('POST', `/tickets/${testTicketId}/refund`, {
        reason: 'Тест',
      });

      // Пытаемся вернуть повторно через UI
      cy.get('[data-cy=ticket-search-input]').type(testTicketNumber);
      cy.get('[data-cy=search-ticket-button]').click();

      // Проверяем сообщение
      cy.contains(/билет уже возвращён/i).should('be.visible');
      cy.get('[data-cy=refund-date]').should('be.visible');
      cy.get('[data-cy=confirm-refund-button]').should('not.exist');
    });
  });

  it('билет не найден', () => {
    cy.get('[data-cy=ticket-search-input]').type('INVALID-TICKET-NUMBER');
    cy.get('[data-cy=search-ticket-button]').click();

    cy.contains(/билет не найден/i).should('be.visible');
    cy.get('[data-cy=ticket-details]').should('not.exist');
  });

  it('расчёт штрафа в зависимости от времени до отправления', () => {
    // Билет на рейс через 25 часов (штраф 0%)
    const departureTime1 = new Date(Date.now() + 25 * 60 * 60 * 1000).toISOString();
    
    cy.apiRequest('POST', '/trips', { routeId: 'test-route', departureTime: departureTime1, totalSeats: 50 })
      .then((tripResponse) => {
        return cy.apiRequest('POST', '/tickets', {
          tripId: tripResponse.body.data.id,
          passengerId: 'test-passenger-001',
          seatNumber: 1,
          price: 1000,
        });
      })
      .then((ticketResponse) => {
        cy.get('[data-cy=ticket-search-input]').clear().type(ticketResponse.body.data.ticketNumber);
        cy.get('[data-cy=search-ticket-button]').click();

        // Проверяем штраф 0%
        cy.get('[data-cy=penalty-amount]').should('contain', '0');
        cy.get('[data-cy=refund-amount]').should('contain', '1 000');
      });

    // Билет на рейс через 10 часов (штраф 10%)
    const departureTime2 = new Date(Date.now() + 10 * 60 * 60 * 1000).toISOString();
    
    cy.apiRequest('POST', '/trips', { routeId: 'test-route', departureTime: departureTime2, totalSeats: 50 })
      .then((tripResponse) => {
        return cy.apiRequest('POST', '/tickets', {
          tripId: tripResponse.body.data.id,
          passengerId: 'test-passenger-002',
          seatNumber: 2,
          price: 1000,
        });
      })
      .then((ticketResponse) => {
        cy.visit('/pos/refund');
        cy.get('[data-cy=ticket-search-input]').clear().type(ticketResponse.body.data.ticketNumber);
        cy.get('[data-cy=search-ticket-button]').click();

        // Проверяем штраф 10%
        cy.get('[data-cy=penalty-amount]').should('contain', '100');
        cy.get('[data-cy=refund-amount]').should('contain', '900');
      });
  });

  it('печать чека возврата', () => {
    cy.get('[data-cy=ticket-search-input]').type(ticketNumber);
    cy.get('[data-cy=search-ticket-button]').click();
    cy.get('[data-cy=confirm-refund-button]').click();
    cy.get('[data-cy=confirm-yes-button]').click();

    // Ждём фискализации
    cy.get('[data-cy=fiscal-success]', { timeout: 30000 }).should('be.visible');

    // Печатаем чек
    cy.get('[data-cy=print-refund-receipt-button]').click();
    cy.get('[data-cy=print-progress]').should('be.visible');
    cy.get('[data-cy=print-success]', { timeout: 10000 }).should('be.visible');
  });
});
