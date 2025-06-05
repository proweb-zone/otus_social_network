import http from 'k6/http';

export let options = {
  vus: 1,
  duration: '30s',
  thresholds: {
    http_req_failed: ['rate<0.01'],
    http_req_duration: ['p(95)<200'],
  },
};

export default function () {

  let email = generateMixedString(6);

  const payload = JSON.stringify({
    "email": `${email}@gmail.com`,
    "password": "123123Vc",
    "first_name": "Семен",
    "last_name": "Семеныч",
    "birth_date": "2025-06-03",
    "gender": "man",
    "hobby": "travel, tennis",
    "city": "Moscow"
});

  let res = http.post('http://localhost:3002/user/register', payload, {
    headers: {
      'Content-Type': 'application/json',
    },
  });

  // Проверка ответа
  check(res, {
    'status was 200': (r) => r.status === 200,
  });

   sleep(2);
}

function generateMixedString(length) {
  let result = '';
  const characters = 'abcdefghijklmnopqrstuvwxyz0123456789'; // Буквы и цифры

  for (let i = 0; i < length; i++) {
    const randomIndex = Math.floor(Math.random() * characters.length);
    result += characters.charAt(randomIndex);
  }

  return result;
}
