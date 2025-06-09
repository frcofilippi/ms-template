import http from 'k6/http';
import { sleep, check } from 'k6';

export let options = {
  stages: [
    { duration: '10s', target: 50 },   // ramp up
    { duration: '10s', target: 100 },   // stay steady
    { duration: '10s', target: 0 },    // ramp down
  ],
};

const productIds = [115, 114]; 

export default function () {
  const id = productIds[Math.floor(Math.random() * productIds.length)];
  const url = `http://localhost:2020/api/v1/product/114`;

  const res = http.get(url);

  check(res, {
    'status is 200': (r) => r.status === 200,
    'response is not empty': (r) => r.body.length > 0,
  });

  sleep(1);  
}
