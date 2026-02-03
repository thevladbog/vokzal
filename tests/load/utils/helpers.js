import http from 'k6/http';

const BASE_URL = __ENV.API_URL || 'http://localhost:8080/api/v1';

/**
 * Авторизация пользователя
 * @param {string} username - логин
 * @param {string} password - пароль
 * @returns {Object} - { success: boolean, accessToken: string, refreshToken: string }
 */
export function login(username, password) {
  const res = http.post(
    `${BASE_URL}/auth/login`,
    JSON.stringify({ login: username, password }),
    {
      headers: { 'Content-Type': 'application/json' },
    }
  );

  if (res.status === 200) {
    const body = JSON.parse(res.body);
    return {
      success: true,
      accessToken: body.data.accessToken,
      refreshToken: body.data.refreshToken,
    };
  }

  return { success: false };
}

/**
 * Генерация случайных данных пассажира
 * @returns {Object} - данные пассажира
 */
export function generatePassenger() {
  const lastNames = ['Иванов', 'Петров', 'Сидоров', 'Васильев', 'Михайлов', 'Фёдоров'];
  const firstNames = ['Иван', 'Пётр', 'Сидор', 'Василий', 'Михаил', 'Фёдор'];
  const middleNames = ['Иванович', 'Петрович', 'Сидорович', 'Васильевич', 'Михайлович'];

  return {
    lastName: lastNames[Math.floor(Math.random() * lastNames.length)],
    firstName: firstNames[Math.floor(Math.random() * firstNames.length)],
    middleName: middleNames[Math.floor(Math.random() * middleNames.length)],
    documentType: 'passport',
    documentSeries: String(Math.floor(1000 + Math.random() * 9000)),
    documentNumber: String(Math.floor(100000 + Math.random() * 900000)),
    phone: '+7999' + String(Math.floor(1000000 + Math.random() * 9000000)),
    email: `test${Math.floor(Math.random() * 10000)}@example.com`,
  };
}

/**
 * Генерация случайной даты в будущем
 * @param {number} maxDays - максимальное количество дней вперёд
 * @returns {string} - дата в формате ISO
 */
export function generateFutureDate(maxDays = 7) {
  const date = new Date();
  date.setDate(date.getDate() + Math.floor(Math.random() * maxDays) + 1);
  return date.toISOString().split('T')[0];
}

/**
 * Задержка с вариацией
 * @param {number} baseSeconds - базовая задержка в секундах
 * @param {number} variation - вариация (0-1)
 */
export function sleepRandom(baseSeconds, variation = 0.3) {
  const min = baseSeconds * (1 - variation);
  const max = baseSeconds * (1 + variation);
  const delay = min + Math.random() * (max - min);
  return delay;
}

/**
 * HTTP запрос с автоматической retry логикой
 * @param {string} method - HTTP метод
 * @param {string} url - URL
 * @param {Object} body - тело запроса
 * @param {Object} params - параметры запроса
 * @param {number} maxRetries - максимальное количество повторов
 * @returns {Object} - response
 */
export function requestWithRetry(method, url, body, params = {}, maxRetries = 3) {
  let lastResponse;
  
  for (let i = 0; i < maxRetries; i++) {
    if (method === 'GET') {
      lastResponse = http.get(url, params);
    } else if (method === 'POST') {
      lastResponse = http.post(url, body, params);
    } else if (method === 'PUT') {
      lastResponse = http.put(url, body, params);
    } else if (method === 'DELETE') {
      lastResponse = http.del(url, params);
    }

    // Успех или не retry-able ошибка
    if (lastResponse.status < 500 || i === maxRetries - 1) {
      return lastResponse;
    }

    // Exponential backoff
    sleep(Math.pow(2, i));
  }

  return lastResponse;
}
