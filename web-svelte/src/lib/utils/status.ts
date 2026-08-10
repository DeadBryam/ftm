import type { TunnelState } from '$lib/types';

export type StatusKey = 'running' | 'starting' | 'installing' | 'error' | 'stopped';

export interface StatusColors {
  bg: string;
  text: string;
  dot: string;
}

export interface StatusInfo {
  key: StatusKey;
  textKey: string;
}

export const STATUS_COLORS: Record<StatusKey, StatusColors> = {
  running: {
    bg: 'bg-status-running/40',
    text: 'text-status-running',
    dot: 'bg-status-running/95'
  },
  starting: {
    bg: 'bg-status-starting/40',
    text: 'text-status-starting',
    dot: 'bg-status-starting/95'
  },
  installing: {
    bg: 'bg-status-installing/40',
    text: 'text-status-installing',
    dot: 'bg-status-installing/95'
  },
  error: {
    bg: 'bg-status-error/40',
    text: 'text-status-error',
    dot: 'bg-status-error/95'
  },
  stopped: {
    bg: 'bg-status-stopped/40',
    text: 'text-status-stopped',
    dot: 'bg-status-stopped/95'
  }
};

const STATUS_MAP: Record<TunnelState, StatusInfo> = {
  online: { key: 'running', textKey: 'online' },
  starting: { key: 'starting', textKey: 'starting' },
  connecting: { key: 'starting', textKey: 'connecting' },
  installing: { key: 'installing', textKey: 'installing' },
  downloading: { key: 'installing', textKey: 'downloading' },
  need_installing: { key: 'stopped', textKey: 'need_installing' },
  stopping: { key: 'starting', textKey: 'stopping' },
  stopped: { key: 'stopped', textKey: 'stopped' },
  offline: { key: 'stopped', textKey: 'offline' },
  timeout: { key: 'error', textKey: 'timeout' },
  error: { key: 'error', textKey: 'error' }
};

const RUNNING_STATES: TunnelState[] = [
  'online',
  'starting',
  'connecting',
  'installing',
  'downloading',
  'stopping'
];

export function statusInfo(state: TunnelState): StatusInfo {
  return STATUS_MAP[state] ?? STATUS_MAP.error;
}

export function statusColors(state: TunnelState): StatusColors {
  return STATUS_COLORS[statusInfo(state).key];
}

export function isRunningState(state: TunnelState): boolean {
  return RUNNING_STATES.includes(state);
}

export function isInstallingState(state: TunnelState): boolean {
  return state === 'installing' || state === 'downloading';
}
