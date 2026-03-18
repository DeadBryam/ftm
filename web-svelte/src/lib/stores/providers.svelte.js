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
        const res = await fetch('/api/providers');
        providers = await res.json();
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
    const res = await fetch('/api/detect-port');
    const data = await res.json();
    return data.suggested || 30000;
  } catch {
    return 30000;
  }
}
