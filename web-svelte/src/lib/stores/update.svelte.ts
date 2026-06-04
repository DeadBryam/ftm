import { writable } from 'svelte/store';
import { updateApi, type UpdateInfo } from '$lib/api';

interface UpdateState {
  info: UpdateInfo | null;
  applying: boolean;
  error: string | null;
}

const state = writable<UpdateState>({ info: null, applying: false, error: null });

export const updateStore = {
  subscribe: state.subscribe,
  async check() {
    try {
      const info = await updateApi.get();
      state.update((s) => ({ ...s, info, error: null }));
    } catch {
      // silent: no update info or offline
    }
  },
  async apply() {
    state.update((s) => ({ ...s, applying: true, error: null }));
    try {
      await updateApi.apply();
    } catch (e) {
      state.update((s) => ({
        ...s,
        applying: false,
        error: e instanceof Error ? e.message : String(e),
      }));
    }
  },
  set(info: UpdateInfo) {
    state.update((s) => ({ ...s, info }));
  },
};
