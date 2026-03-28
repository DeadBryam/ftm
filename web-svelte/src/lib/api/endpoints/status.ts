import { api } from '../client';

export interface StatusResponse {
  port: number;
  version: string;
  notificationsStatus: 'pending' | 'granted' | 'rejected';
}

export async function getStatus(): Promise<StatusResponse> {
  return api.get('status').json();
}
