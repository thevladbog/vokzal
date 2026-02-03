/// <reference types="cypress" />

describe('Посадка пассажиров', () => {
  let tripId: string;
  let ticketIds: string[] = [];
  let ticketQRCodes: string[] = [];

  before(() => {
    // Создаём тестовый рейс и билеты
    cy.loginAsCashier();

    const departureTime = new Date(Date.now() + 2 * 60 * 60 * 1000).toISOString();

    cy.apiRequest('POST', '/trips', {
      routeId: 'test-route-boarding',
      departureTime,
      totalSeats: 50,
      platform: '5',
      busNumber: 'А123БВ',
    }).then((tripResponse) => {
      tripId = tripResponse.body.data.id;

      // Создаём несколько билетов на этот рейс
      const passengers = [
        { lastName: 'Иванов', firstName: 'Иван', seat: 1 },
        { lastName: 'Петров', firstName: 'Пётр', seat: 2 },
        { lastName: 'Сидоров', firstName: 'Сидор', seat: 3 },
      ];

      passengers.forEach((passenger) => {
        cy.apiRequest('POST', '/passengers', passenger).then((passengerResponse) => {
          const passengerId = passengerResponse.body.data.id;

          cy.apiRequest('POST', '/tickets', {
            tripId,
            passengerId,
            seatNumber: passenger.seat,
            price: 1500,
          }).then((ticketResponse) => {
            ticketIds.push(ticketResponse.body.data.id);
            ticketQRCodes.push(ticketResponse.body.data.qrCode);
          });
        });
      });
    });
  });

  beforeEach(() => {
    cy.loginAsController();
  });

  it('выбор активного рейса для контроля', () => {
    cy.visit('/controller');

    // Проверяем список активных рейсов
    cy.get('[data-cy=trip-list]').should('be.visible');
    cy.get('[data-cy=trip-card]').should('have.length.greaterThan', 0);

    // Выбираем наш тестовый рейс
    cy.get('[data-cy=trip-card]').contains('А123БВ').click();

    // Проверяем, что перешли на страницу сканирования
    cy.url().should('include', '/scan');
    cy.get('[data-cy=trip-info]').should('be.visible');
    cy.get('[data-cy=boarding-stats]').should('be.visible');
  });

  it('успешное сканирование QR-кода билета', () => {
    cy.visit(`/controller/scan/${tripId}`);

    // Эмулируем сканирование (в реальности это камера)
    cy.window().then((win) => {
      win.postMessage({ type: 'QR_SCANNED', qrCode: ticketQRCodes[0] }, '*');
    });

    // Проверяем успешную посадку
    cy.get('[data-cy=scan-result-success]').should('be.visible');
    cy.contains(/посадка выполнена/i).should('be.visible');
    
    // Проверяем отображение информации о пассажире
    cy.get('[data-cy=passenger-name]').should('contain', 'Иванов');
    cy.get('[data-cy=seat-number]').should('contain', '1');

    // Проверяем обновление статистики
    cy.get('[data-cy=boarded-count]').should('contain', '1');
  });

  it('повторное сканирование того же билета', () => {
    cy.visit(`/controller/scan/${tripId}`);

    // Первое сканирование
    cy.window().then((win) => {
      win.postMessage({ type: 'QR_SCANNED', qrCode: ticketQRCodes[1] }, '*');
    });

    cy.get('[data-cy=scan-result-success]').should('be.visible');

    // Ждём несколько секунд и сканируем повторно
    cy.wait(3000);

    cy.window().then((win) => {
      win.postMessage({ type: 'QR_SCANNED', qrCode: ticketQRCodes[1] }, '*');
    });

    // Должно появиться предупреждение
    cy.get('[data-cy=scan-result-warning]').should('be.visible');
    cy.contains(/уже прошёл посадку/i).should('be.visible');
    
    // Проверяем время первой посадки
    cy.get('[data-cy=boarding-time]').should('be.visible');
  });

  it('сканирование недействительного QR-кода', () => {
    cy.visit(`/controller/scan/${tripId}`);

    cy.window().then((win) => {
      win.postMessage({ type: 'QR_SCANNED', qrCode: 'INVALID-QR-CODE-12345' }, '*');
    });

    // Должна появиться ошибка
    cy.get('[data-cy=scan-result-error]').should('be.visible');
    cy.contains(/билет не найден/i).should('be.visible');
  });

  it('сканирование билета на другой рейс', () => {
    // Создаём билет на другой рейс
    const anotherDepartureTime = new Date(Date.now() + 3 * 60 * 60 * 1000).toISOString();

    cy.apiRequest('POST', '/trips', {
      routeId: 'another-route',
      departureTime: anotherDepartureTime,
      totalSeats: 50,
    }).then((anotherTripResponse) => {
      const anotherTripId = anotherTripResponse.body.data.id;

      cy.apiRequest('POST', '/passengers', {
        lastName: 'Васильев',
        firstName: 'Василий',
      }).then((passengerResponse) => {
        cy.apiRequest('POST', '/tickets', {
          tripId: anotherTripId,
          passengerId: passengerResponse.body.data.id,
          seatNumber: 10,
          price: 2000,
        }).then((ticketResponse) => {
          const wrongTripQR = ticketResponse.body.data.qrCode;

          // Пытаемся отсканировать этот билет на нашем рейсе
          cy.visit(`/controller/scan/${tripId}`);

          cy.window().then((win) => {
            win.postMessage({ type: 'QR_SCANNED', qrCode: wrongTripQR }, '*');
          });

          // Должна появиться ошибка
          cy.get('[data-cy=scan-result-error]').should('be.visible');
          cy.contains(/неверный рейс/i).should('be.visible');
        });
      });
    });
  });

  it('сканирование возвращённого билета', () => {
    // Создаём и возвращаем билет
    cy.apiRequest('POST', '/passengers', {
      lastName: 'Михайлов',
      firstName: 'Михаил',
    }).then((passengerResponse) => {
      cy.apiRequest('POST', '/tickets', {
        tripId,
        passengerId: passengerResponse.body.data.id,
        seatNumber: 15,
        price: 1500,
      }).then((ticketResponse) => {
        const ticketId = ticketResponse.body.data.id;
        const refundedQR = ticketResponse.body.data.qrCode;

        // Возвращаем билет
        cy.apiRequest('POST', `/tickets/${ticketId}/refund`, {
          reason: 'Тест',
        });

        // Пытаемся отсканировать
        cy.visit(`/controller/scan/${tripId}`);

        cy.window().then((win) => {
          win.postMessage({ type: 'QR_SCANNED', qrCode: refundedQR }, '*');
        });

        // Должна появиться ошибка
        cy.get('[data-cy=scan-result-error]').should('be.visible');
        cy.contains(/билет возвращён/i).should('be.visible');
      });
    });
  });

  it('отображение статистики посадки', () => {
    cy.visit(`/controller/scan/${tripId}`);

    // Проверяем начальную статистику
    cy.get('[data-cy=sold-tickets]').should('be.visible');
    cy.get('[data-cy=boarded-tickets]').should('be.visible');
    cy.get('[data-cy=progress-bar]').should('be.visible');

    // Сканируем несколько билетов
    ticketQRCodes.forEach((qrCode, index) => {
      cy.window().then((win) => {
        win.postMessage({ type: 'QR_SCANNED', qrCode }, '*');
      });

      cy.wait(3000); // Ждём между сканированиями

      // Проверяем обновление счётчика
      cy.get('[data-cy=boarded-count]').should('contain', (index + 1).toString());
    });

    // Проверяем прогресс
    cy.get('[data-cy=progress-bar]').invoke('attr', 'value').should('be.greaterThan', 0);
  });

  it('звуковые и вибро сигналы при сканировании', () => {
    cy.visit(`/controller/scan/${tripId}`);

    // Мокаем Audio и vibrate API
    cy.window().then((win) => {
      cy.spy(win.navigator, 'vibrate').as('vibrate');
      // Audio тяжело проверить в тесте, но можем проверить, что он создаётся
    });

    // Успешное сканирование
    cy.window().then((win) => {
      win.postMessage({ type: 'QR_SCANNED', qrCode: ticketQRCodes[0] }, '*');
    });

    cy.get('[data-cy=scan-result-success]').should('be.visible');
    
    // Проверяем вибрацию (если поддерживается)
    cy.get('@vibrate').should('be.called');
  });

  it('история последних сканирований', () => {
    cy.visit(`/controller/scan/${tripId}`);

    // Сканируем несколько билетов
    ticketQRCodes.slice(0, 3).forEach((qrCode) => {
      cy.window().then((win) => {
        win.postMessage({ type: 'QR_SCANNED', qrCode }, '*');
      });
      cy.wait(2000);
    });

    // Проверяем историю
    cy.get('[data-cy=recent-scans]').should('be.visible');
    cy.get('[data-cy=scan-history-item]').should('have.length', 3);

    // Проверяем информацию в истории
    cy.get('[data-cy=scan-history-item]').first().within(() => {
      cy.get('[data-cy=passenger-name]').should('not.be.empty');
      cy.get('[data-cy=boarding-time]').should('be.visible');
    });
  });

  it('офлайн режим PWA', () => {
    cy.visit(`/controller/scan/${tripId}`);

    // Переходим в офлайн режим
    cy.window().then((win) => {
      cy.stub(win.navigator, 'onLine').value(false);
    });

    // Проверяем индикатор офлайн режима
    cy.contains(/офлайн режим/i).should('be.visible');

    // Сканирование должно работать (данные кэшированы)
    cy.window().then((win) => {
      win.postMessage({ type: 'QR_SCANNED', qrCode: ticketQRCodes[0] }, '*');
    });

    // В офлайн режиме данные сохраняются локально
    cy.get('[data-cy=offline-queue]').should('be.visible');
  });
});
