import { defineConfig } from 'cypress';

export default defineConfig({
  e2e: {
    baseUrl: 'http://localhost:3000',
    viewportWidth: 1920,
    viewportHeight: 1080,
    video: true,
    screenshotOnRunFailure: true,
    defaultCommandTimeout: 10000,
    requestTimeout: 10000,
    responseTimeout: 10000,
    
    env: {
      apiUrl: 'http://localhost:8080/api/v1',
      adminLogin: 'admin',
      adminPassword: 'admin123',
      cashierLogin: 'cashier',
      cashierPassword: 'cashier123',
      controllerLogin: 'controller',
      controllerPassword: 'controller123',
    },

    setupNodeEvents(on, config) {
      // implement node event listeners here
      require('@cypress/code-coverage/task')(on, config);
      return config;
    },

    specPattern: 'cypress/e2e/**/*.cy.{js,jsx,ts,tsx}',
    supportFile: 'cypress/support/e2e.ts',
  },
});
