export type ToastType = 'success' | 'error' | 'info' | 'warning' | 'alert';

export type Toast = {
  id: number;
  message: string;
  type: ToastType;
  duration?: number;
};

let toasts: Toast[] = $state([]);

function add(message: string, type: ToastType = "info", duration = 3000) {
  const next = { id: Date.now() + Math.random(), message, type, duration };
  toasts = [next, ...toasts];
}

function remove(id: number) {
  toasts = toasts.filter(t => t.id !== id);
}

export const toast = {
  get toasts() { return toasts; },
  show: add,
  success: (msg: string, duration?: number) => add(msg, "success", duration),
  error: (msg: string, duration?: number) => add(msg, "error", duration),
  info: (msg: string, duration?: number) => add(msg, "info", duration),
  remove,
};

export function useToast() { return toast; }