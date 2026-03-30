import { api } from '../client';

export interface StatusResponse {
  port: number;
  version: string;
  notificationsStatus: 'pending' | 'granted' | 'rejected';
  wsClients: number;
}

export const statusApi = {
  get: (): Promise<StatusResponse> => api.get('status').json(),
};
