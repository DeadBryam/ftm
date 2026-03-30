import { api } from '../client';

type NotificationsUpdateStatus = 'granted' | 'rejected';

interface NotificationsUpdateResponse {
  success: boolean;
  status: 'pending' | 'granted' | 'rejected';
}

export const notificationsApi = {
  updateStatus: (status: NotificationsUpdateStatus): Promise<NotificationsUpdateResponse> =>
    api.post('notifications', { json: { status } }).json(),
};
