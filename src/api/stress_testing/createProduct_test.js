import http from 'k6/http';
import { sleep, check } from 'k6';

export let options = {
  stages: [
    { duration: '10s', target: 100 },   // ramp up
    { duration: '5s', target: 10 },   // stay steady
    { duration: '10s', target: 0 },    // ramp down
  ],
};

export default function () {
  const payload = JSON.stringify({
    name: `Test Product ${Math.floor(Math.random() * 10000)}`,
    cost: 20.34,
  });

  const headers = { 'Content-Type': 'application/json' };

  const res = http.post('http://localhost:2020/api/v1/product', payload, { headers });

  check(res, {
    'status is 200 or 201': (r) => r.status === 200 || r.status === 201,
    'response time < 500ms': (r) => r.timings.duration < 500,
  });

  sleep(1);
}
