import ky from 'ky';

export const api = ky.create({
  prefixUrl: '/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
  retry: {
    // Reads only. Retrying a write means a failed "start tunnel" fires the
    // provider up to three times, and a retried delete can remove a tunnel the
    // user recreated in between. Only GET is safe to repeat.
    limit: 2,
    methods: ['get'],
    statusCodes: [408, 413, 429, 500, 502, 503, 504],
    backoffLimit: 10000,
  },
});
