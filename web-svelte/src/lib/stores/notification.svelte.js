let permission = $state('default');
let useOSNotifications = $state(false);
let enabled = $state(false);

const notificationStore = {
  get permission() { return permission; },
  get useOSNotifications() { return useOSNotifications; },
  get enabled() { return enabled; },
  
  init() {
    const saved = localStorage.getItem('ftm-notification-pref');
    if (saved === 'granted') {
      permission = 'granted';
      useOSNotifications = true;
    }
    enabled = localStorage.getItem('ftm-notifications-enabled') === 'true';
  },
  
  async requestPermission() {
    if (!('Notification' in window)) return false;
    
    if (permission === 'default') {
      const result = await Notification.requestPermission();
      permission = result;
      useOSNotifications = result === 'granted';
    }
    
    if (permission === 'granted') {
      enabled = true;
      localStorage.setItem('ftm-notification-pref', 'granted');
      localStorage.setItem('ftm-notifications-enabled', 'true');
    }
    
    return permission === 'granted';
  },
  
  enable() {
    if (permission === 'granted') {
      enabled = true;
      localStorage.setItem('ftm-notifications-enabled', 'true');
    }
  },
  
  disable() {
    enabled = false;
    localStorage.setItem('ftm-notifications-enabled', 'false');
  },
  
  notify(title, body) {
    if (!enabled) return;
    
    if (useOSNotifications && permission === 'granted') {
      new Notification(title, { body });
    }
  },
  
  notifyOnline(name, url) { this.notify('Tunnel Active', `${name} - ${url}`); },
  notifyError(name, err) { this.notify('Tunnel Error', `${name}: ${err}`); },
  notifyExpiring(name, mins) { this.notify('Tunnel Expiring', `${name}: ${mins} min remaining`); },
  notifyExpired(name) { this.notify('Tunnel Expired', `${name} session ended`); }
};

export function useNotifications() {
  return notificationStore;
}
