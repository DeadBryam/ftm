import { api } from '../client';
import type { Provider } from '../types';

export const providersApi = {
  getAll: (): Promise<Provider[]> => api.get('providers').json<Provider[]>(),
};
