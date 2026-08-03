<script lang="ts">
  import { page } from '$app/stores';
  import { cn } from '$lib/utils/cn';
  import { t } from '$lib/stores/i18n.svelte';
  import { Settings } from 'lucide-svelte';

  let isSettings = $derived($page.url.pathname === '/settings');
</script>

<header class="z-10 mb-3 flex shrink-0 items-center justify-between border-b border-border pb-3">
  <div class="flex min-w-0 items-center gap-3">
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
  </div>

  <a
    href="/settings"
    class={cn(
      "rounded-control p-2 transition-colors",
      isSettings ? "bg-primary/20 text-primary" : "hover:bg-secondary"
    )}
    aria-label={t('settings')}
  >
    <Settings size={18} />
  </a>
</header>