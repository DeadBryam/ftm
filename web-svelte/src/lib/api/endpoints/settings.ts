import { api } from '../client';

export interface Settings {
  notifications_enabled: boolean;
  notification_sound: boolean;
}

export const settingsApi = {
  get: (): Promise<Settings> => api.get('settings').json(),
  update: (settings: Partial<Settings>): Promise<Settings> =>
    api.patch('settings', { json: settings }).json(),
};
