import { updateApi, type UpdateInfo } from '$lib/api';

let info: UpdateInfo | null = $state(null);
let applying = $state(false);
let error: string | null = $state(null);

export function useUpdate() {
  return {
    get info() { return info; },
    get applying() { return applying; },
    get error() { return error; },

    async check() {
      try {
        info = await updateApi.get();
        error = null;
      } catch {
        // silent: no update info or offline
      }
    },

    async apply() {
      applying = true;
      error = null;
      try {
        await updateApi.apply();
      } catch (e) {
        applying = false;
        error = e instanceof Error ? e.message : String(e);
      }
    },

    set(next: UpdateInfo) {
      info = next;
    },
  };
}
