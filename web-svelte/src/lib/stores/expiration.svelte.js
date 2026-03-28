import { useNotifications } from './notification.svelte.js';

const DEFAULT_THRESHOLDS = [30, 15, 10, 5, 1];

let timers = $state({});

export function useExpirationMonitor() {
    const notifications = useNotifications();
    
    function getThresholds() {
        const saved = localStorage.getItem('ftm-expiration-thresholds');
        return saved ? JSON.parse(saved) : DEFAULT_THRESHOLDS;
    }
    
    function setThresholds(thresholds) {
        localStorage.setItem('ftm-expiration-thresholds', JSON.stringify(thresholds));
    }
    
    function start(tunnel) {
        if (!tunnel.expiresAt || !notifications.enabled) return;
        
        const thresholds = getThresholds();
        const expiresAt = new Date(tunnel.expiresAt).getTime();
        
        thresholds.forEach(mins => {
            const triggerAt = expiresAt - (mins * 60 * 1000);
            const now = Date.now();
            
            if (triggerAt > now) {
                const delay = triggerAt - now;
                const key = `${tunnel.id}-${mins}`;
                
                timers[key] = setTimeout(() => {
                    notifications.notifyExpiring(tunnel.name, mins);
                }, delay);
            }
        });
    }
    
    function stop(tunnelId) {
        Object.keys(timers)
            .filter(k => k.startsWith(tunnelId))
            .forEach(k => {
                clearTimeout(timers[k]);
                delete timers[k];
            });
    }
    
    function stopAll() {
        Object.values(timers).forEach(t => clearTimeout(t));
        timers = {};
    }
    
    return { start, stop, stopAll, getThresholds, setThresholds };
}
