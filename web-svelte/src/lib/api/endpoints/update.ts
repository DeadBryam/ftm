import { api } from '../client';

export interface UpdateInfo {
  current: string;
  latest: string;
  tag: string;
  assetName: string;
  releaseUrl: string;
  hasUpdate: boolean;
}

export const updateApi = {
  get: (): Promise<UpdateInfo> => api.get('update').json(),
  apply: (): Promise<{ ok: boolean }> => api.post('update', { json: {} }).json(),
};
