import { describe, expect, it, vi, beforeEach } from 'vitest';

const payload = {
  notifications_enabled: 'granted',
  notification_sound: true,
  language: 'en',
  onboarded: true,
  autostart_supported: true,
  autostart_enabled: false,
};

vi.mock('../api', () => ({
  settingsApi: {
    get: vi.fn(async () => ({ ...payload })),
    update: vi.fn(async (partial) => ({ ...payload, ...partial })),
  },
}));

vi.mock('./notification.svelte', () => ({
  useNotifications: () => ({ applySettings: vi.fn() }),
}));

import { useSettings } from './settings.svelte';
import { settingsApi } from '../api';

describe('settings store', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('treats autostart as unsupported until the backend says otherwise', () => {
    const store = useSettings();

    expect(store.settings.autostart_supported).toBe(false);
    expect(store.settings.autostart_enabled).toBe(false);
  });

  it('loads the autostart capability from the API', async () => {
    const store = useSettings();
    await store.load();

    expect(store.settings.autostart_supported).toBe(true);
    expect(store.settings.autostart_enabled).toBe(false);
  });

  it('keeps the capability flag after toggling autostart', async () => {
    const store = useSettings();
    await store.load();
    await store.update({ autostart_enabled: true });

    expect(settingsApi.update).toHaveBeenCalledWith({ autostart_enabled: true });
    expect(store.settings.autostart_enabled).toBe(true);
    expect(store.settings.autostart_supported).toBe(true);
  });

  it('rolls the toggle back when the system refuses the change', async () => {
    const store = useSettings();
    await store.load();

    vi.mocked(settingsApi.update).mockRejectedValueOnce(new Error('disabled by user'));

    await expect(store.update({ autostart_enabled: true })).rejects.toThrow();
    expect(store.settings.autostart_enabled).toBe(false);
  });
});
