import { api } from '../client';

export interface NotificationsUpdateRequest {
  status: 'granted' | 'rejected';
}

export interface NotificationsUpdateResponse {
  success: boolean;
  status: 'pending' | 'granted' | 'rejected';
}

export async function updateNotificationsStatus(
  status: 'granted' | 'rejected'
): Promise<NotificationsUpdateResponse> {
  return api
    .post('notifications', { json: { status } })
    .json();
}
