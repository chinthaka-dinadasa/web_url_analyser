import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  vus: 20,
  duration: '15s',
};

// Array of different URLs to test
const testUrls = [
    'https://javatodev.com',
    'https://google.com', 
    'https://github.com',
    'https://stackoverflow.com',
    'https://reddit.com',
    'https://twitter.com',
    'https://linkedin.com',
    'https://youtube.com',
    'https://amazon.com',
    'https://netflix.com',
    'https://javatodev.com'
];

export default function () {
  const randomUrl = testUrls[Math.floor(Math.random() * testUrls.length)];
  
  const payload = JSON.stringify({
    url: randomUrl
  });
  
  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };
  
  const res = http.post('http://localhost:8080/process-web-url', payload, params);
  
  check(res, {
    'status is 200': (r) => r.status === 200,
    'response time OK': (r) => r.timings.duration < 2000,
  });
  
  sleep(0.5);
}