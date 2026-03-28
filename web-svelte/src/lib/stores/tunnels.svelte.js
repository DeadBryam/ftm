import { tunnelsApi } from '$lib/api';
import { useNotifications } from './notification.svelte.js';
import { useExpirationMonitor } from './expiration.svelte.js';

let tunnels = $state([]);
let loading = $state(true);
let error = $state(null);
let socket = $state(null);
let installProgress = $state({});

const notifications = useNotifications();
const expirationMonitor = useExpirationMonitor();

function handleMessage(msg) {
  if (msg.type === 'install') {
    installProgress = { ...installProgress, [msg.provider]: msg };
    return;
  }
  
  if (!msg.id) return;
  
  const idx = tunnels.findIndex(t => t.id === msg.id);
  
  if (msg.name !== undefined || msg.provider !== undefined || msg.port !== undefined) {
    tunnels = tunnels.map(t => {
      if (t.id === msg.id) {
        return {
          ...t,
          name: msg.name ?? t.name,
          provider: msg.provider ?? t.provider,
          port: msg.port ?? t.port,
          state: msg.state ?? t.state,
          publicUrl: msg.publicUrl ?? t.publicUrl,
          errorMessage: msg.errorMessage ?? t.errorMessage,
          expiresAt: msg.expiresAt ?? t.expiresAt
        };
      }
      return t;
    });
    return;
  }
  
  if (idx === -1) return;
  
  const oldTunnel = tunnels[idx];
  const newState = msg.state || 'stopped';
  
  if (oldTunnel.state === newState && 
      oldTunnel.publicUrl === (msg.publicUrl ?? oldTunnel.publicUrl) &&
      oldTunnel.errorMessage === (msg.errorMessage ?? oldTunnel.errorMessage)) {
    return;
  }
  
  tunnels = tunnels.map(t => {
    if (t.id === msg.id) {
      return {
        ...t,
        state: newState,
        publicUrl: msg.publicUrl ?? t.publicUrl,
        errorMessage: msg.errorMessage ?? t.errorMessage,
        expiresAt: msg.expiresAt ?? t.expiresAt
      };
    }
    return t;
  });
  
  const updatedTunnel = tunnels.find(t => t.id === msg.id);
  if (!updatedTunnel) return;
  
  if (newState === 'online' && oldTunnel.state !== 'online') {
    notifications.notifyOnline(updatedTunnel.name, msg.publicUrl);
    if (msg.expiresAt) expirationMonitor.start(updatedTunnel);
    return;
  }
  
  if (newState === 'stopped' && oldTunnel.state === 'online') {
    expirationMonitor.stop(msg.id);
    notifications.notify('Tunnel Stopped', `${updatedTunnel.name} has been stopped`, 'info');
    return;
  }
  
  if (newState === 'error' && oldTunnel.state !== 'error') {
    notifications.notifyError(updatedTunnel.name, msg.errorMessage);
    return;
  }
  
  if (newState === 'timeout' && oldTunnel.state !== 'timeout') {
    notifications.notify('Timeout', `${updatedTunnel.name} could not connect`, 'error');
    return;
  }
  
  if (newState === 'online') {
    const updatedProgress = { ...installProgress };
    delete updatedProgress[updatedTunnel.provider];
    installProgress = updatedProgress;
  }
}

function connect() {
  if (socket && socket.readyState === WebSocket.OPEN) return;
  
  loading = true;
  notifications.init();
  
  tunnelsApi.getAll()
    .then(data => {
      tunnels = data;
      data.forEach(t => {
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
  
  const ws = new WebSocket(`ws://${window.location.host}/ws/events`);
  
  ws.onopen = () => {
    console.log('[WS] Connected');
  };
  
  ws.onmessage = (e) => {
    try {
      handleMessage(JSON.parse(e.data));
    } catch (err) {
      console.error('[WS] Parse error:', err);
    }
  };
  
  ws.onclose = () => {
    console.log('[WS] Disconnected');
    socket = null;
    setTimeout(connect, 3000);
  };
  
  ws.onerror = (e) => {
    console.error('[WS] Error:', e);
    error = 'Connection error';
  };
  
  socket = ws;
}

function disconnect() {
  if (socket) {
    socket.close();
    socket = null;
  }
  expirationMonitor.stopAll();
}

function start(id) {
  tunnels = tunnels.map(t => t.id === id ? { ...t, state: 'starting' } : t);
  tunnelsApi.start(id).catch(e => {
    tunnels = tunnels.map(t => t.id === id ? { ...t, state: 'error', errorMessage: e.message } : t);
  });
}

function stop(id) {
  tunnelsApi.stop(id).catch(() => {
    tunnels = tunnels.map(t => t.id === id ? { ...t, state: 'stopped', publicUrl: null } : t);
  });
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
    update
  };
}
