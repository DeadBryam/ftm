import { tunnelsApi } from '$lib/api';

let tunnels = $state([]);
let loading = $state(true);
let error = $state(null);
let eventSource = $state(null);
let installProgress = $state({});

function handleSSEEvent(msg) {
  // Evento de instalación
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
  
  // Eventos de túnel tienen `id` directo
  if (!msg.id) return;
  
  tunnels = tunnels.map(t => {
    if (t.id !== msg.id) return t;
    
    const newState = msg.state || 'stopped';
    return {
      ...t,
      state: newState,
      publicUrl: msg.publicUrl ?? t.publicUrl,
      errorMessage: msg.errorMessage ?? t.errorMessage
    };
  });
  
  // Limpiar progress cuando pasa a online
  const tunnel = tunnels.find(t => t.id === msg.id);
  if (tunnel && tunnel.state === 'online') {
    const updated = { ...installProgress };
    delete updated[tunnel.provider];
    installProgress = updated;
  }
}

function connect() {
  if (eventSource) return;
  
  loading = true;
  
  tunnelsApi.getAll()
    .then(data => {
      tunnels = data;
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
  tunnelsApi.stop(id).catch(() => {});
}

function remove(id) {
  tunnelsApi.delete(id).catch(() => {});
  tunnels = tunnels.filter(t => t.id !== id);
}

function add(data) {
  tunnelsApi.create(data).then(newTunnel => {
    tunnels = [...tunnels, newTunnel];
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
    create: add
  };
}
