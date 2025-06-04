import http from 'k6/http';
import { sleep } from 'k6';
import { check, group } from 'k6';

export const options = {
  vus: 1000, // Количество виртуальных пользователей
  duration: '1m', // Длительность теста в минутах
  thresholds: {
    http_req_failed: ['rate<0.01'], // Не более 1% ошибок
    http_req_duration: ['p(95)<200'], // 95% запросов должны выполняться менее чем за 200 мс
  },
};

export default function () {
  group('API uptime check', () => {
    const response = http.get('http://localhost:3002/user/get/2');
    check(response, {
      "status code should be 200": res => res.status === 200,
    });

    const response1 = http.get('http://localhost:3002/user/search/лев абра');
    check(response1, {
      "status code should be 200": res => res.status === 200,
    });
  });

  sleep(1);
}
