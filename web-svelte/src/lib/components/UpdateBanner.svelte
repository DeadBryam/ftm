<script lang="ts">
  import { onMount } from 'svelte';
  import { ArrowUpCircle, ExternalLink } from 'lucide-svelte';
  import { useUpdate } from '$lib/stores/update.svelte';
  import { t } from '$lib/stores/i18n.svelte';
  import { subscribeWsMessages } from '$lib/api/ws';
  import type { UpdateMethod } from '$lib/api';
  import Button from './Button.svelte';

  const update = useUpdate();

  onMount(() => {
    void update.check();

    const unsub = subscribeWsMessages((msg) => {
      const m = msg as {
        type?: string;
        latest?: string;
        tag?: string;
        releaseUrl?: string;
        assetName?: string;
        current?: string;
        method?: UpdateMethod;
      } | null;
      if (m?.type === 'update_available') {
        update.set({
          current: m.current ?? update.info?.current ?? '',
          latest: m.latest ?? '',
          tag: m.tag ?? '',
          assetName: m.assetName ?? '',
          releaseUrl: m.releaseUrl ?? '',
          method: m.method ?? 'self',
          hasUpdate: true,
        });
      }
    });
    return unsub;
  });

  const selfUpdatable = $derived(update.info?.method === 'self');

  const hint = $derived.by(() => {
    switch (update.info?.method) {
      case 'store':
        return t('update_manual_store');
      case 'homebrew':
        return t('update_manual_homebrew');
      case 'download':
        return t('update_manual_download');
      default:
        return '';
    }
  });
</script>

{#if update.info?.hasUpdate}
  <div
    class="ftm-enter mb-app flex flex-wrap items-center gap-3 rounded-card border border-primary/30 bg-primary/10 px-3 py-2 text-sm"
  >
    <span class="flex min-w-0 items-center gap-2 text-text-heading">
      <ArrowUpCircle size={16} class="shrink-0 text-primary" />
      <span class="truncate font-semibold">
        {t('update_web_banner', { latest: update.info.latest })}
      </span>
    </span>

    {#if hint}
      <span class="min-w-0 truncate text-xs text-text-muted">{hint}</span>
    {/if}

    <div class="ml-auto flex shrink-0 items-center gap-2">
      <a
        href={update.info.releaseUrl}
        target="_blank"
        rel="noreferrer"
        class="flex items-center gap-1 text-xs text-text-muted underline-offset-2 transition-colors hover:text-primary hover:underline"
      >
        {t('update_web_notes')}
        <ExternalLink size={12} />
      </a>

      {#if update.applying}
        <span class="text-xs text-text-muted">{t('update_web_restart')}</span>
      {:else if selfUpdatable}
        <Button variant="primary" size="sm" onclick={() => update.apply()}>
          {t('update_web_button')}
        </Button>
      {:else}
        <Button
          variant="default"
          size="sm"
          icon={ExternalLink}
          onclick={() => window.open(update.info?.releaseUrl, '_blank', 'noreferrer')}
        >
          {t('update_web_download')}
        </Button>
      {/if}
    </div>
  </div>
{/if}

{#if update.error}
  <div
    class="ftm-enter mb-app flex items-center gap-2 rounded-card border border-status-error/40 bg-status-error/10 px-3 py-2 text-sm text-status-error"
  >
    {t('update_apply_failed', { 0: update.error })}
  </div>
{/if}
