<script>
  import { onMount } from 'svelte';
  import { useTunnels } from '$lib/stores/tunnels.svelte';
  import { useToast } from '$lib/stores/toast.svelte';
  import { useProviders, detectPort } from '$lib/stores/providers.svelte';

  const store = useTunnels();
  const toast = useToast();
  const providerStore = useProviders();

  let formData = $state({ name: '', provider: 'cloudflared', localPort: 30000 });

  onMount(async () => {
    const detectedPort = await detectPort();
    formData.localPort = detectedPort;
  });

  async function handleSubmit(e) {
    e.preventDefault();
    const name = formData.name;
    await store.create(formData);
    formData = { name: '', provider: 'cloudflared', localPort: formData.localPort };
    toast.success(`Connection "${name}" created`);
  }
</script>

<section class="panel create-panel">
  <div class="panel-header">
    <h2>New Connection</h2>
  </div>
  <div class="panel-body">
    <form class="create-form" onsubmit={handleSubmit}>
      <div class="field-group">
        <label for="name">Connection Name</label>
        <input type="text" id="name" bind:value={formData.name} placeholder="e.g. Storm King's Thunder" required autocomplete="off" />
      </div>

      <div class="field-row">
        <div class="field-group">
          <label for="port">Port</label>
          <input type="number" id="port" bind:value={formData.localPort} min="1" max="65535" required />
        </div>
        <div class="field-group">
          <label for="provider">Provider</label>
          <div class="select-wrap">
            <select id="provider" bind:value={formData.provider} required>
              {#each providerStore.providers as p}
                <option value={p.id}>{p.name}</option>
              {/each}
            </select>
          </div>
        </div>
      </div>

      <button type="submit" class="btn btn-primary btn-full">
        Create Connection
      </button>
    </form>
  </div>
</section>

<style>
  .panel {
    background: var(--card-bg);
    border-radius: 12px;
    box-shadow: 0 1px 3px rgba(0,0,0,0.05);
    border: 1px solid var(--border-color);
    display: flex;
    flex-direction: column;
    overflow: hidden;
    transition: box-shadow 0.2s ease, transform 0.2s ease;
    opacity: 0;
    transform: translateY(30px);
    animation: panelIn 0.5s cubic-bezier(0.16, 1, 0.3, 1) forwards;
    animation-delay: 0.1s;
    min-height: 0;
  }

  .panel:hover {
    box-shadow: 0 8px 24px rgba(0,0,0,0.08), 0 2px 6px rgba(0,0,0,0.04);
    transform: translateY(-1px);
  }

  @keyframes panelIn {
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .panel-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 14px 18px;
    border-bottom: 1px solid var(--border-light);
    background: var(--url-bg);
    flex-shrink: 0;
  }

  .panel-header h2 {
    font-family: 'Crimson Pro', Georgia, serif;
    font-size: 17px;
    font-weight: 600;
    color: var(--text-heading);
    margin: 0;
  }

  .panel-body {
    padding: 18px;
    flex: 1;
    overflow-y: auto;
    min-height: 0;
  }

  .field-group {
    margin-bottom: 14px;
  }

  .field-group input {
    height: 42px;
  }

  .field-row {
    display: grid;
    grid-template-columns: 90px 1fr;
    gap: 10px;
  }

  @media (max-width: 500px) {
    .field-row {
      grid-template-columns: 1fr;
    }
  }

  label {
    display: block;
    font-size: 12px;
    font-weight: 500;
    color: var(--text-muted);
    margin-bottom: 5px;
  }

  input, select {
    width: 100%;
    padding: 8px 10px;
    border: 1px solid var(--border-color);
    border-radius: 8px !important;
    font-size: 13px;
    font-family: inherit;
    background: var(--card-bg);
    box-sizing: border-box;
  }

  input:focus, select:focus {
    outline: none;
    border-color: var(--primary-color);
    box-shadow: 0 0 0 3px rgba(0,0,0,0.08);
  }

  .select-wrap {
    position: relative;
  }

  .select-wrap select {
    height: 42px;
    appearance: none;
    cursor: pointer;
  }

  .btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    padding: 10px 16px;
    border-radius: 8px;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    border: 1px solid #d6d3d1;
    background: var(--card-bg);
    color: #44403c;
    transition: all 0.2s ease;
  }

  .btn:hover {
    background: #f5f5f4;
  }

  .btn-primary {
    background: linear-gradient(135deg, #92400e 0%, #b45309 100%);
    color: white;
    border-color: #92400e;
    box-shadow: 0 4px 12px rgba(146, 64, 14, 0.3), 0 2px 4px rgba(146, 64, 14, 0.2);
  }

  .btn-primary:hover {
    background: linear-gradient(135deg, #78350f 0%, #92400e 100%);
    border-color: #78350f;
    box-shadow: 0 6px 16px rgba(146, 64, 14, 0.35), 0 3px 6px rgba(146, 64, 14, 0.25);
    transform: translateY(-1px);
  }

  .btn-full {
    width: 100%;
    padding: 12px;
  }

  @media (prefers-reduced-motion: reduce) {
    .panel {
      animation: none;
      opacity: 1;
      transform: none;
    }
  }
</style>
