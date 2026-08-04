<script lang="ts">
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { t } from '$lib/stores/i18n.svelte';
  import { Settings, ChevronLeft } from 'lucide-svelte';

  const isHome = $derived($page.url.pathname === '/');
  const isSettings = $derived($page.url.pathname === '/settings');

  function handleBack() {
    if (history.length > 1) {
      history.back();
    } else {
      goto('/');
    }
  }
</script>

<header class="z-10 mb-3 flex shrink-0 items-center justify-between border-b border-border pb-3">
  <a
    href="/"
    class="flex min-w-0 items-center gap-3 rounded-control pr-1 transition-opacity hover:opacity-80"
    aria-label={t('go_home')}
  >
    <img
      src="/favicon.png"
      alt={t('app_name')}
      class="h-9 w-9 shrink-0 rounded-control object-cover sm:h-10 sm:w-10"
    />
    <div class="min-w-0">
      <h1 class="m-0 font-serif text-xl font-bold tracking-tight text-text-heading sm:text-2xl">
        {t('app_name')}
      </h1>
      <p class="m-0 truncate text-xs font-medium text-text-muted">{t('app_tagline')}</p>
    </div>
  </a>

  {#if isSettings}
    <button
      type="button"
      onclick={handleBack}
      class="inline-flex cursor-pointer items-center gap-1.5 rounded-control px-3 py-1.5 text-sm text-text-muted transition-colors hover:bg-secondary hover:text-text-heading"
    >
      <ChevronLeft size={14} />
      {t('go_home')}
    </button>
  {:else if isHome}
    <a
      href="/settings"
      class="rounded-control p-2 transition-colors hover:bg-secondary"
      aria-label={t('settings')}
    >
      <Settings size={18} />
    </a>
  {/if}
</header>