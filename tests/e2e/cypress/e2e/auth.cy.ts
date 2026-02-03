/// <reference types="cypress" />

describe('Авторизация', () => {
  beforeEach(() => {
    cy.visit('/login');
  });

  it('успешная авторизация администратора', () => {
    cy.get('[data-cy=login-input]').type(Cypress.env('adminLogin'));
    cy.get('[data-cy=password-input]').type(Cypress.env('adminPassword'));
    cy.get('[data-cy=login-button]').click();

    // Проверяем, что перешли на главную страницу
    cy.url().should('not.include', '/login');
    cy.url().should('include', '/dashboard');

    // Проверяем, что токен сохранён
    cy.window().then((win) => {
      const token = win.localStorage.getItem('accessToken');
      expect(token).to.not.be.null;
      expect(token).to.not.be.empty;
    });
  });

  it('неверный логин', () => {
    cy.get('[data-cy=login-input]').type('nonexistent');
    cy.get('[data-cy=password-input]').type('wrongpassword');
    cy.get('[data-cy=login-button]').click();

    // Должны остаться на странице логина
    cy.url().should('include', '/login');

    // Должно появиться сообщение об ошибке
    cy.contains(/неверный логин или пароль/i).should('be.visible');
  });

  it('неверный пароль', () => {
    cy.get('[data-cy=login-input]').type(Cypress.env('adminLogin'));
    cy.get('[data-cy=password-input]').type('wrongpassword');
    cy.get('[data-cy=login-button]').click();

    cy.url().should('include', '/login');
    cy.contains(/неверный логин или пароль/i).should('be.visible');
  });

  it('пустые поля', () => {
    cy.get('[data-cy=login-button]').click();

    // Проверяем валидацию полей
    cy.get('[data-cy=login-input]').should('have.attr', 'aria-invalid', 'true');
    cy.get('[data-cy=password-input]').should('have.attr', 'aria-invalid', 'true');
  });

  it('выход из системы', () => {
    // Авторизуемся
    cy.loginAsAdmin();

    // Выходим
    cy.logout();

    // Проверяем, что токен удалён
    cy.window().then((win) => {
      const token = win.localStorage.getItem('accessToken');
      expect(token).to.be.null;
    });
  });

  it('автоматический переход после авторизации', () => {
    // Пытаемся открыть защищённую страницу без авторизации
    cy.visit('/dashboard');

    // Должны быть перенаправлены на логин
    cy.url().should('include', '/login');

    // Авторизуемся
    cy.loginAsAdmin();

    // После успешной авторизации должны вернуться на /dashboard
    cy.url().should('include', '/dashboard');
  });

  it('обновление токена при истечении', () => {
    cy.loginAsAdmin();

    // Эмулируем истечение access token (удаляем его)
    cy.window().then((win) => {
      win.localStorage.removeItem('accessToken');
    });

    // Делаем запрос, который должен триггернуть refresh
    cy.apiRequest('GET', '/trips').then((response) => {
      expect(response.status).to.eq(200);
    });

    // Проверяем, что новый токен получен
    cy.window().then((win) => {
      const token = win.localStorage.getItem('accessToken');
      expect(token).to.not.be.null;
    });
  });
});
