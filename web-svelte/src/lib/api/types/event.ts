export interface TunnelEvent {
  id: string;
  state?: string;
  publicUrl?: string;
  errorMessage?: string;
  provider?: string;
}

export interface InstallEvent {
  type: 'install';
  provider: string;
  step: string;
  percent: number;
  downloaded: number;
  total: number;
}

export interface UnknownEvent {
  type?: 'unknown';
  [key: string]: unknown;
}

export type WSEvent = TunnelEvent | InstallEvent | UnknownEvent;

export interface InstallProgress {
  step: string;
  percent: number;
  downloaded: number;
  total: number;
}
