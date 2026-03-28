import { api } from '../client';

export interface LogStreamOptions {
  onLine: (line: string) => void;
  onError?: (error: Event) => void;
  onClose?: () => void;
}


export function createLogStream(
  tunnelId: string,
  options: LogStreamOptions
): { close: () => void } {
  const es = new EventSource(`/api/logs/${tunnelId}/stream`);

  es.onmessage = (e: MessageEvent) => {
    options.onLine(e.data);
  };

  es.onerror = (e) => {
    if (options.onError) {
      options.onError(e);
    }
  };

  es.close = () => {
    if (options.onClose) {
      options.onClose();
    }
  };

  return {
    close: () => es.close(),
  };
}


export async function getLogs(tunnelId: string): Promise<string> {
  const res = await api.get(`logs/${tunnelId}`);
  return res.text();
}
