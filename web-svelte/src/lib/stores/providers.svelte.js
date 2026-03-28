import { providersApi } from '$lib/api';
import { api } from '$lib/api';

let providers = $state([]);
let loading = $state(false);
let error = $state(null);

export function useProviders() {
  return {
    get providers() { return providers; },
    get loading() { return loading; },
    get error() { return error; },
    
    async fetch() {
      loading = true;
      error = null;
      try {
        providers = await providersApi.getAll();
      } catch (e) {
        error = e.message;
      } finally {
        loading = false;
      }
    }
  };
}

export async function detectPort() {
  try {
    const data = await api.get('detect-port').json();
    return data.suggested || 30000;
  } catch {
    return 30000;
  }
}
