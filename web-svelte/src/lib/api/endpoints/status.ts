import { api } from '../client';

export interface StatusResponse {
  port: number;
  version: string;
}

export const statusApi = {
  get: (): Promise<StatusResponse> => api.get('status').json(),
};
