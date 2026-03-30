<script lang="ts">
  import { onDestroy, onMount } from 'svelte';
  import { page } from '$app/stores';
  import { statusApi } from '$lib/api';
  import { cn } from '$lib/utils/cn';
  import { Settings } from 'lucide-svelte';

  let isSettings = $derived($page.url.pathname === '/settings');
  let wsClients = $state(0);
  let intervalId: ReturnType<typeof setInterval> | null = null;

  async function refreshStatus() {
    try {
      const status = await statusApi.get();
      wsClients = status.wsClients;
    } catch {
      wsClients = 0;
    }
  }

  onMount(() => {
    void refreshStatus();
    intervalId = setInterval(() => {
      void refreshStatus();
    }, 5000);
  });

  onDestroy(() => {
    if (intervalId) {
      clearInterval(intervalId);
      intervalId = null;
    }
  });
</script>

<header class="flex justify-between items-center mb-6 pb-5 border-b flex-shrink-0 z-10 border-border">
  <div class="flex items-center gap-4">
    <img src="/favicon.png" alt="Foundry Tunnel Manager" class="w-10 h-10 sm:w-12 sm:h-12 rounded-xl object-cover" />
    <div class="text">
      <h1 class="font-serif text-2xl sm:text-3xl lg:text-4xl font-bold m-0 mb-1 tracking-tight text-text-heading">Foundry Tunnel Manager</h1>
      <p class="text-xs sm:text-sm m-0 font-medium text-text-muted">Share your world with players everywhere</p>
    </div>
  </div>
  
  <div class="flex items-center gap-3">
    {#if wsClients > 0}
      <span class="px-2.5 py-1 rounded-full text-xs font-semibold bg-primary/15 text-primary">
        Sessions: {wsClients}
      </span>
    {/if}
    <a
      href="/settings"
      class={cn(
        "p-2 rounded-lg transition-colors",
        isSettings ? "bg-primary/20 text-primary" : "hover:bg-secondary"
      )}
      aria-label="Settings"
    >
      <Settings size={20} />
    </a>
  </div>
</header>
