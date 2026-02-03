import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend, Counter } from 'k6/metrics';
import { login } from '../utils/helpers.js';

// Кастомные метрики
const searchFailureRate = new Rate('search_failures');
const searchDuration = new Trend('search_duration');
const tripsFound = new Counter('trips_found');

// Конфигурация нагрузки
export const options = {
  stages: [
    { duration: '1m', target: 20 },   // 0 -> 20 пользователей
    { duration: '5m', target: 100 },  // 20 -> 100 пользователей
    { duration: '10m', target: 100 }, // Держим 100 пользователей
    { duration: '2m', target: 0 },    // 100 -> 0
  ],
  thresholds: {
    'http_req_duration': ['p(95)<800', 'p(99)<1500'],
    'http_req_failed': ['rate<0.02'],
    'search_failures': ['rate<0.05'],
    'search_duration': ['p(95)<500', 'avg<300'],
  },
};

const BASE_URL = __ENV.API_URL || 'http://localhost:8080/api/v1';

// Популярные направления для поиска
const popularRoutes = [
  { from: 'Ростов-на-Дону', to: 'Москва' },
  { from: 'Ростов-на-Дону', to: 'Санкт-Петербург' },
  { from: 'Москва', to: 'Ростов-на-Дону' },
  { from: 'Краснодар', to: 'Сочи' },
  { from: 'Ростов-на-Дону', to: 'Волгоград' },
  { from: 'Краснодар', to: 'Ростов-на-Дону' },
];

export default function () {
  // Авторизуемся как пассажир (опционально, можно искать без авторизации)
  const useAuth = Math.random() > 0.3; // 70% авторизованных
  let accessToken = null;

  if (useAuth) {
    const authResult = login('passenger', 'passenger123');
    if (!authResult.success) {
      sleep(1);
      return;
    }
    accessToken = authResult.accessToken;
  }

  // Выбираем случайное направление
  const route = popularRoutes[Math.floor(Math.random() * popularRoutes.length)];

  // Случайная дата (от сегодня до +7 дней)
  const date = new Date();
  date.setDate(date.getDate() + Math.floor(Math.random() * 7));
  const dateStr = date.toISOString().split('T')[0];

  // 1. Поиск станций отправления
  const fromStationRes = http.get(
    `${BASE_URL}/stations?name=${encodeURIComponent(route.from)}`,
    {
      headers: accessToken ? { 'Authorization': `Bearer ${accessToken}` } : {},
      tags: { name: 'search_stations_from' },
    }
  );

  const fromStationCheck = check(fromStationRes, {
    'from station search status is 200': (r) => r.status === 200,
    'from station found': (r) => {
      const body = JSON.parse(r.body);
      return body.data && body.data.length > 0;
    },
  });

  if (!fromStationCheck) {
    sleep(1);
    return;
  }

  const fromStationId = JSON.parse(fromStationRes.body).data[0].id;

  sleep(0.5);

  // 2. Поиск станций назначения
  const toStationRes = http.get(
    `${BASE_URL}/stations?name=${encodeURIComponent(route.to)}`,
    {
      headers: accessToken ? { 'Authorization': `Bearer ${accessToken}` } : {},
      tags: { name: 'search_stations_to' },
    }
  );

  const toStationCheck = check(toStationRes, {
    'to station search status is 200': (r) => r.status === 200,
    'to station found': (r) => {
      const body = JSON.parse(r.body);
      return body.data && body.data.length > 0;
    },
  });

  if (!toStationCheck) {
    sleep(1);
    return;
  }

  const toStationId = JSON.parse(toStationRes.body).data[0].id;

  sleep(0.5);

  // 3. Основной поиск рейсов
  const searchStart = new Date();
  const searchRes = http.get(
    `${BASE_URL}/trips?fromStationId=${fromStationId}&toStationId=${toStationId}&date=${dateStr}`,
    {
      headers: accessToken ? { 'Authorization': `Bearer ${accessToken}` } : {},
      tags: { name: 'search_trips' },
    }
  );

  const searchSuccess = check(searchRes, {
    'trip search status is 200': (r) => r.status === 200,
    'trip search returns data': (r) => {
      const body = JSON.parse(r.body);
      return body.data !== undefined;
    },
  });

  // Записываем метрики
  searchFailureRate.add(!searchSuccess);
  searchDuration.add(new Date() - searchStart);

  if (searchSuccess) {
    const trips = JSON.parse(searchRes.body).data;
    tripsFound.add(trips.length);

    // 4. Если найдены рейсы, запрашиваем детали случайного рейса (30% вероятность)
    if (trips.length > 0 && Math.random() > 0.7) {
      const randomTrip = trips[Math.floor(Math.random() * trips.length)];
      
      sleep(1);

      const tripDetailsRes = http.get(
        `${BASE_URL}/trips/${randomTrip.id}`,
        {
          headers: accessToken ? { 'Authorization': `Bearer ${accessToken}` } : {},
          tags: { name: 'get_trip_details' },
        }
      );

      check(tripDetailsRes, {
        'trip details status is 200': (r) => r.status === 200,
        'trip details has route info': (r) => {
          const body = JSON.parse(r.body);
          return body.data && body.data.route;
        },
      });
    }
  }

  sleep(2);
}

export function setup() {
  console.log('Starting search trips load test');
  console.log(`Target API: ${BASE_URL}`);
  
  // Прогреваем кэш
  http.get(`${BASE_URL}/stations`);
  
  return { startTime: new Date() };
}

export function teardown(data) {
  const duration = (new Date() - data.startTime) / 1000;
  console.log(`Test completed in ${duration.toFixed(2)} seconds`);
}
