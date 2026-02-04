import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

// Кастомные метрики
const loginFailureRate = new Rate('login_failures');
const loginDuration = new Trend('login_duration');

// Конфигурация нагрузки
export const options = {
  stages: [
    { duration: '30s', target: 10 },  // Разогрев: 0 -> 10 пользователей за 30с
    { duration: '1m', target: 50 },   // Рост: 10 -> 50 пользователей за 1м
    { duration: '3m', target: 50 },   // Стабильная нагрузка: 50 пользователей 3м
    { duration: '30s', target: 0 },   // Снижение: 50 -> 0 за 30с
  ],
  thresholds: {
    'http_req_duration': ['p(95)<500', 'p(99)<1000'], // 95% < 500ms, 99% < 1s
    'http_req_failed': ['rate<0.01'],                 // <1% ошибок
    'login_failures': ['rate<0.05'],                  // <5% неудачных логинов
    'login_duration': ['p(95)<300'],                  // 95% логинов < 300ms
  },
};

const BASE_URL = __ENV.API_URL || 'http://localhost:8080/api/v1';

// Тестовые данные
const users = [
  { login: 'admin', password: 'admin123', role: 'admin' },
  { login: 'cashier1', password: 'cashier123', role: 'cashier' },
  { login: 'cashier2', password: 'cashier123', role: 'cashier' },
  { login: 'controller1', password: 'controller123', role: 'controller' },
];

export default function () {
  // Выбираем случайного пользователя
  const user = users[Math.floor(Math.random() * users.length)];

  // 1. Логин
  const loginStart = new Date();
  const loginRes = http.post(
    `${BASE_URL}/auth/login`,
    JSON.stringify({
      username: user.login,
      password: user.password,
    }),
    {
      headers: { 'Content-Type': 'application/json' },
      tags: { name: 'login' },
    }
  );

  const loginSuccess = check(loginRes, {
    'login status is 200': (r) => r.status === 200,
    'login returns access token': (r) => {
      if (!r.body) return false;
      try {
        const body = JSON.parse(r.body);
        return body && body.data && !!body.data.access_token;
      } catch {
        return false;
      }
    },
    'login returns refresh token': (r) => {
      if (!r.body) return false;
      try {
        const body = JSON.parse(r.body);
        return body && body.data && !!body.data.refresh_token;
      } catch {
        return false;
      }
    },
  });

  // Записываем метрики логина
  loginFailureRate.add(!loginSuccess);
  loginDuration.add(new Date() - loginStart);

  if (!loginSuccess) {
    sleep(1);
    return;
  }

  // Извлекаем токены (API: access_token, refresh_token)
  let accessToken, refreshToken;
  try {
    const body = JSON.parse(loginRes.body);
    accessToken = body && body.data && body.data.access_token;
    refreshToken = body && body.data && body.data.refresh_token;
  } catch {
    sleep(1);
    return;
  }
  if (!accessToken || !refreshToken) {
    sleep(1);
    return;
  }

  // 2. Проверка профиля
  const meRes = http.get(`${BASE_URL}/auth/me`, {
    headers: {
      'Authorization': `Bearer ${accessToken}`,
    },
    tags: { name: 'get_profile' },
  });

  check(meRes, {
    'get profile status is 200': (r) => r.status === 200,
    'profile has correct role': (r) => {
      if (!r.body) return false;
      try {
        const body = JSON.parse(r.body);
        return body && body.data && body.data.role === user.role;
      } catch {
        return false;
      }
    },
  });

  sleep(1);

  // 3. Refresh токена (50% вероятность)
  if (Math.random() > 0.5) {
    const refreshRes = http.post(
      `${BASE_URL}/auth/refresh`,
      JSON.stringify({ refreshToken }),
      {
        headers: { 'Content-Type': 'application/json' },
        tags: { name: 'refresh_token' },
      }
    );

    check(refreshRes, {
      'refresh status is 200': (r) => r.status === 200,
      'refresh returns new access token': (r) => {
        if (!r.body) return false;
        try {
          const body = JSON.parse(r.body);
          return body && body.data && !!body.data.access_token;
        } catch {
          return false;
        }
      },
    });

    sleep(1);
  }

  // 4. Logout (backend expects X-Refresh-Token header)
  const logoutRes = http.post(
    `${BASE_URL}/auth/logout`,
    '{}',
    {
      headers: {
        'Content-Type': 'application/json',
        'X-Refresh-Token': refreshToken,
      },
      tags: { name: 'logout' },
    }
  );

  check(logoutRes, {
    'logout status is 200': (r) => r.status === 200,
  });

  sleep(1);
}

// Функция для setup (выполняется один раз перед тестами)
export function setup() {
  console.log('Starting auth load test');
  console.log(`Target API: ${BASE_URL}`);
  return { startTime: new Date() };
}

// Функция для teardown (выполняется один раз после тестов)
export function teardown(data) {
  const duration = (new Date() - data.startTime) / 1000;
  console.log(`Test completed in ${duration.toFixed(2)} seconds`);
}
