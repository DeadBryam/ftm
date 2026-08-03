import { describe, expect, it, vi, beforeEach } from 'vitest';

vi.mock('$lib/api', () => ({
  providersApi: {
    getAll: vi.fn(async () => [
      { id: 'cloudflared', name: 'Cloudflared' },
      { id: 'bore', name: 'bore' },
    ]),
    detectPort: vi.fn(async () => 30000),
  },
}));

import { useProviders } from './providers.svelte';
import { providersApi } from '$lib/api';

describe('providers store', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('loads providers from the API including bore', async () => {
    const store = useProviders();
    await store.fetch();

    expect(providersApi.getAll).toHaveBeenCalledOnce();
    expect(store.providers.map((p) => p.id)).toContain('bore');
    expect(store.error).toBeNull();
    expect(store.loading).toBe(false);
  });
});
