import { api } from '../client';

export interface StatusResponse {
  port: number;
  version: string;
  notificationsStatus: 'pending' | 'granted' | 'rejected';
}

export const statusApi = {
  get: (): Promise<StatusResponse> => api.get('status').json(),
};
