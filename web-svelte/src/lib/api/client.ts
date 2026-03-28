import ky from 'ky';

export const api = ky.create({
  prefixUrl: '/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
  retry: {
    limit: 3,
    methods: ['get', 'post', 'put', 'delete', 'patch'],
    statusCodes: [408, 413, 429, 500, 502, 503, 504],
    backoffLimit: 10000,
  },
});
