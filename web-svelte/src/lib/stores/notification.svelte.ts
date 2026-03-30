import { useToast, type ToastType } from './toast.svelte';
import { notificationsApi } from '$lib/api';

type NotificationStatus = 'pending' | 'granted' | 'rejected';

let status = $state<NotificationStatus>('pending');
let useOSNotifications = $state(false);
let soundEnabled = $state(true);

const toast = useToast();

let audioContext: AudioContext | null = null;

async function initAudio() {
  if (typeof window === 'undefined') return;
  if (audioContext) return;

  const AudioContextClass = (window.AudioContext || (window as typeof window & { webkitAudioContext?: typeof AudioContext }).webkitAudioContext);
  if (!AudioContextClass) return;

  try {
    audioContext = new AudioContextClass();
    if (audioContext.state === 'suspended') {
      await audioContext.resume();
    }
  } catch {
    audioContext = null;
  }
}

async function playSound(type: string) {
  if (!soundEnabled || typeof window === 'undefined') return;

  await initAudio();
  if (!audioContext) return;

  const oscillator = audioContext.createOscillator();
  const gainNode = audioContext.createGain();

  oscillator.connect(gainNode);
  gainNode.connect(audioContext.destination);

  if (type === 'success') {
    oscillator.frequency.setValueAtTime(523.25, audioContext.currentTime);
    oscillator.frequency.setValueAtTime(659.25, audioContext.currentTime + 0.1);
    gainNode.gain.setValueAtTime(0.3, audioContext.currentTime);
    gainNode.gain.exponentialRampToValueAtTime(0.01, audioContext.currentTime + 0.3);
    oscillator.start(audioContext.currentTime);
    oscillator.stop(audioContext.currentTime + 0.3);
  } else if (type === 'error') {
    oscillator.frequency.setValueAtTime(200, audioContext.currentTime);
    gainNode.gain.setValueAtTime(0.3, audioContext.currentTime);
    gainNode.gain.exponentialRampToValueAtTime(0.01, audioContext.currentTime + 0.2);
    oscillator.start(audioContext.currentTime);
    oscillator.stop(audioContext.currentTime + 0.2);
  } else if (type === 'warning') {
    oscillator.frequency.setValueAtTime(440, audioContext.currentTime);
    oscillator.frequency.setValueAtTime(880, audioContext.currentTime + 0.15);
    gainNode.gain.setValueAtTime(0.25, audioContext.currentTime);
    gainNode.gain.exponentialRampToValueAtTime(0.01, audioContext.currentTime + 0.3);
    oscillator.start(audioContext.currentTime);
    oscillator.stop(audioContext.currentTime + 0.3);
  } else if (type === 'alert') {
    oscillator.frequency.setValueAtTime(880, audioContext.currentTime);
    oscillator.frequency.setValueAtTime(440, audioContext.currentTime + 0.1);
    oscillator.frequency.setValueAtTime(880, audioContext.currentTime + 0.2);
    gainNode.gain.setValueAtTime(0.35, audioContext.currentTime);
    gainNode.gain.exponentialRampToValueAtTime(0.01, audioContext.currentTime + 0.3);
    oscillator.start(audioContext.currentTime);
    oscillator.stop(audioContext.currentTime + 0.3);
  } else {
    oscillator.frequency.setValueAtTime(440, audioContext.currentTime);
    gainNode.gain.setValueAtTime(0.2, audioContext.currentTime);
    gainNode.gain.exponentialRampToValueAtTime(0.01, audioContext.currentTime + 0.15);
    oscillator.start(audioContext.currentTime);
    oscillator.stop(audioContext.currentTime + 0.15);
  }
}

const notificationStore = {
  get status() { return status; },
  get enabled() { return status === 'granted'; },
  get soundEnabled() { return soundEnabled; },
  set soundEnabled(value: boolean) {
    soundEnabled = value;
    if (typeof window !== 'undefined') {
      try {
        localStorage.setItem('ftm-sound-enabled', value ? 'true' : 'false');
      } catch {}
    }
  },

  setStatus(newStatus: NotificationStatus) {
    status = newStatus;
  },

  async requestPermission(): Promise<boolean> {
    if (!('Notification' in window)) {
      return false;
    }

    const result = await Notification.requestPermission();
    const granted = result === 'granted';

    if (granted) {
      status = 'granted';
      useOSNotifications = true;
    } else {
      status = 'rejected';
    }

    try {
      await notificationsApi.updateStatus(granted ? 'granted' : 'rejected');
    } catch {
      console.error('Failed to update notification status on server');
    }

    return granted;
  },

  reject() {
    status = 'rejected';
    notificationsApi.updateStatus('rejected').catch(() => {
      console.error('Failed to update notification status on server');
    });
  },

  notify(title: string, body: string, type: ToastType = 'info') {
    if (soundEnabled) {
      playSound(type);
    }

    if (status !== 'granted') return;

    if (useOSNotifications && 'Notification' in window && Notification.permission === 'granted') {
      new Notification(title, { body });
    } else {
      toast.show(body, type);
    }
  },

  notifyOnline(name: string, url: string) {
    this.notify('Tunnel Active', `${name} - ${url}`, 'success');
  },

  notifyError(name: string, err: string) {
    this.notify('Tunnel Error', `${name}: ${err}`, 'error');
  },

  notifyExpiring(name: string, mins: number) {
    const title = mins <= 1 ? 'Last Minute!' : 'Tunnel Expiring';
    const type = mins <= 1 ? 'alert' : 'warning';
    this.notify(title, `${name}: ${mins} min remaining`, type as ToastType);
  },

  notifyExpired(name: string) {
    this.notify('Tunnel Expired', `${name} session ended`, 'error');
  }
};

export function useNotifications() {
  return notificationStore;
}
