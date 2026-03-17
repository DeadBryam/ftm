document.addEventListener('DOMContentLoaded', function() {
    
    document.body.addEventListener('htmx:afterRequest', function(evt) {
        const target = evt.detail.target;
        
        if (evt.detail.failed && evt.detail.xhr.status === 202) {
            const response = evt.detail.xhr.responseText;
            if (response.includes('INSTALLING PROVIDER')) {
                showInstallingToast();
                setTimeout(() => {
                    target.setAttribute('hx-get', '/api/tunnels');
                    target.setAttribute('hx-trigger', 'load');
                    htmx.trigger(target, 'load');
                }, 3000);
            }
        }
    });
    
    document.body.addEventListener('htmx:afterSwap', function(evt) {
        const target = evt.detail.target;
        
        if (target.id === 'tunnels-container') {
            const emptyState = target.querySelector('#empty-state');
            const items = target.querySelectorAll('.connection-item');
            
            if (items.length > 0 && emptyState) {
                emptyState.remove();
            }
            
            if (items.length === 0 && !emptyState) {
                target.innerHTML = `
                    <div class="empty-state" id="empty-state">
                        <div class="empty-state-icon">📡</div>
                        <h3>No connections yet</h3>
                        <p>Create your first connection to share your Foundry world with players.</p>
                    </div>
                `;
            }
            
            updateConnectionCount();
        }
    });
    
    document.body.addEventListener('htmx:confirm', function(evt) {
        if (evt.detail.elt.getAttribute('hx-delete')) {
            if (!confirm('Delete this connection?')) {
                evt.preventDefault();
            }
        }
    });
    
    document.body.addEventListener('htmx:afterDelete', function(evt) {
        const container = document.getElementById('tunnels-container');
        const items = container.querySelectorAll('.connection-item');
        
        if (items.length === 0) {
            container.innerHTML = `
                <div class="empty-state" id="empty-state">
                    <div class="empty-state-icon">📡</div>
                    <h3>No connections yet</h3>
                    <p>Create your first connection to share your Foundry world with players.</p>
                </div>
            `;
        }
        updateConnectionCount();
    });
});

function updateConnectionCount() {
    const container = document.getElementById('tunnels-container');
    const count = container.querySelectorAll('.connection-item').length;
    const badge = document.getElementById('connection-count');
    if (badge) badge.textContent = count;
}

function showInstallingToast() {
    const toast = document.createElement('div');
    toast.className = 'toast installing';
    toast.innerHTML = '⏳ Installing provider, please wait...';
    toast.style.cssText = `
        position: fixed;
        bottom: 20px;
        right: 20px;
        background: #92400e;
        color: white;
        padding: 12px 20px;
        border-radius: 8px;
        box-shadow: 0 4px 6px rgba(0,0,0,0.2);
        z-index: 1000;
        animation: slideIn 0.3s ease-out;
    `;
    document.body.appendChild(toast);
    
    setTimeout(() => {
        toast.style.animation = 'slideOut 0.3s ease-out';
        setTimeout(() => toast.remove(), 300);
    }, 5000);
}

function initDragAndDrop() {
    
}
