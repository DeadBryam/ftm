import { tunnelsApi } from '$lib/api';
import { useNotifications } from './notification.svelte.js';
import { useExpirationMonitor } from './expiration.svelte.js';

let tunnels = $state([]);
let loading = $state(true);
let error = $state(null);
let eventSource = $state(null);
let installProgress = $state({});

const notifications = useNotifications();
const expirationMonitor = useExpirationMonitor();

function handleSSEEvent(msg) {
  
  if (msg.type === 'install') {
    installProgress = {
      ...installProgress,
      [msg.provider]: {
        step: msg.step,
        percent: msg.percent,
        downloaded: msg.downloaded,
        total: msg.total
      }
    };
    return;
  }
  
  
  if (!msg.id) return;
  
  const oldTunnel = tunnels.find(t => t.id === msg.id);
  
  tunnels = tunnels.map(t => {
    if (t.id !== msg.id) return t;
    
    const newState = msg.state || 'stopped';
    return {
      ...t,
      name: msg.name ?? t.name,
      provider: msg.provider ?? t.provider,
      port: msg.port ?? t.port,
      state: newState,
      publicUrl: msg.publicUrl ?? t.publicUrl,
      errorMessage: msg.errorMessage ?? t.errorMessage,
      expiresAt: msg.expiresAt ?? t.expiresAt
    };
  });
  
  const tunnel = tunnels.find(t => t.id === msg.id);
  
  if (!notifications.enabled || !tunnel) return;
  
  if (msg.state === 'online' && (!oldTunnel || oldTunnel.state !== 'online')) {
    notifications.notifyOnline(tunnel.name, msg.publicUrl);
    if (msg.expiresAt) {
      expirationMonitor.start(tunnel);
    }
  }
  
  if (msg.state === 'error' && oldTunnel?.state !== 'error') {
    notifications.notifyError(tunnel.name, msg.errorMessage);
  }
  
  if (msg.state === 'timeout' && oldTunnel?.state !== 'timeout') {
    notifications.notify('Timeout', `${tunnel.name} could not connect`);
  }
  
  if (msg.state === 'stopped' && oldTunnel?.state === 'online') {
    expirationMonitor.stop(msg.id);
  }
  
  if (tunnel && tunnel.state === 'online') {
    const updated = { ...installProgress };
    delete updated[tunnel.provider];
    installProgress = updated;
  }
}

function connect() {
  if (eventSource) return;
  
  loading = true;
  notifications.init();
  
  tunnelsApi.getAll()
    .then(data => {
      tunnels = data;
      tunnels.forEach(t => {
        if (t.state === 'online' && t.expiresAt) {
          expirationMonitor.start(t);
        }
      });
      loading = false;
    })
    .catch(e => {
      error = e.message;
      loading = false;
    });
  
  const es = new EventSource('/api/events');
  
  es.onmessage = (e) => {
    try {
      const msg = JSON.parse(e.data);
      handleSSEEvent(msg);
    } catch {}
  };
  
  es.onerror = () => {
    error = 'Connection lost. Retrying...';
    setTimeout(() => {
      if (es.readyState === EventSource.CLOSED) {
        eventSource = null;
        connect();
      }
    }, 3000);
  };
  
  eventSource = es;
}

function disconnect() {
  if (eventSource) {
    eventSource.close();
    eventSource = null;
  }
  expirationMonitor.stopAll();
}

function start(id) {
  tunnels = tunnels.map(t => t.id === id ? { ...t, state: 'starting' } : t);
  
  tunnelsApi.start(id)
    .then(data => {
      if (data.status === 'installing') {
        tunnels = tunnels.map(t => t.id === id ? { ...t, state: 'installing' } : t);
      }
    })
    .catch(e => {
      tunnels = tunnels.map(t => t.id === id ? { ...t, state: 'error', errorMessage: e.message } : t);
    });
}

function stop(id) {
  tunnels = tunnels.map(t => t.id === id ? { ...t, state: 'stopped', publicUrl: null } : t);
  expirationMonitor.stop(id);
  tunnelsApi.stop(id).catch(() => {});
}

function remove(id) {
  expirationMonitor.stop(id);
  tunnelsApi.delete(id).catch(() => {});
  tunnels = tunnels.filter(t => t.id !== id);
}

function add(data) {
  tunnelsApi.create(data).then(newTunnel => {
    tunnels = [...tunnels, newTunnel];
  });
}

function update(id, data) {
  return tunnelsApi.update(id, data).then(updated => {
    tunnels = tunnels.map(t => t.id === id ? updated : t);
    return updated;
  });
}

function getById(id) {
  return tunnels.find(t => t.id === id);
}

export function useTunnels() {
  return {
    get tunnels() { return tunnels; },
    get loading() { return loading; },
    get error() { return error; },
    get installProgress() { return installProgress; },
    
    connect,
    disconnect,
    start,
    stop,
    delete: remove,
    create: add,
    update,
    getById
  };
}
