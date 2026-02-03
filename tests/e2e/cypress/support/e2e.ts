// ***********************************************************
// This file is processed and loaded automatically before your test files.
//
// You can read more here:
// https://on.cypress.io/configuration
// ***********************************************************

import '@cypress/code-coverage/support';

// Custom commands
declare global {
  namespace Cypress {
    interface Chainable {
      /**
       * Авторизация пользователя
       * @example cy.login('admin', 'admin123')
       */
      login(username: string, password: string): Chainable<void>;

      /**
       * Авторизация администратора
       * @example cy.loginAsAdmin()
       */
      loginAsAdmin(): Chainable<void>;

      /**
       * Авторизация кассира
       * @example cy.loginAsCashier()
       */
      loginAsCashier(): Chainable<void>;

      /**
       * Авторизация контролёра
       * @example cy.loginAsController()
       */
      loginAsController(): Chainable<void>;

      /**
       * Выход из системы
       * @example cy.logout()
       */
      logout(): Chainable<void>;

      /**
       * Получение токена авторизации
       * @example cy.getAuthToken()
       */
      getAuthToken(): Chainable<string>;

      /**
       * API запрос с авторизацией
       * @example cy.apiRequest('GET', '/trips')
       */
      apiRequest(method: string, url: string, body?: any): Chainable<any>;
    }
  }
}

// Команда для авторизации
Cypress.Commands.add('login', (username: string, password: string) => {
  cy.visit('/login');
  cy.get('[data-cy=login-input]').type(username);
  cy.get('[data-cy=password-input]').type(password);
  cy.get('[data-cy=login-button]').click();
  cy.url().should('not.include', '/login');
});

// Команда для авторизации администратора
Cypress.Commands.add('loginAsAdmin', () => {
  cy.login(Cypress.env('adminLogin'), Cypress.env('adminPassword'));
});

// Команда для авторизации кассира
Cypress.Commands.add('loginAsCashier', () => {
  cy.login(Cypress.env('cashierLogin'), Cypress.env('cashierPassword'));
});

// Команда для авторизации контролёра
Cypress.Commands.add('loginAsController', () => {
  cy.login(Cypress.env('controllerLogin'), Cypress.env('controllerPassword'));
});

// Команда для выхода
Cypress.Commands.add('logout', () => {
  cy.get('[data-cy=user-menu]').click();
  cy.get('[data-cy=logout-button]').click();
  cy.url().should('include', '/login');
});

// Команда для получения токена
Cypress.Commands.add('getAuthToken', () => {
  return cy.window().then((win) => {
    return win.localStorage.getItem('accessToken') || '';
  });
});

// Команда для API запросов с авторизацией
Cypress.Commands.add('apiRequest', (method: string, url: string, body?: any) => {
  return cy.getAuthToken().then((token) => {
    return cy.request({
      method,
      url: `${Cypress.env('apiUrl')}${url}`,
      body,
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
  });
});

// Глобальная обработка ошибок
Cypress.on('uncaught:exception', (err) => {
  // Игнорируем ошибки ResizeObserver (часто встречаются в React приложениях)
  if (err.message.includes('ResizeObserver')) {
    return false;
  }
  return true;
});

// Очистка локального хранилища перед каждым тестом
beforeEach(() => {
  cy.clearLocalStorage();
  cy.clearCookies();
});
