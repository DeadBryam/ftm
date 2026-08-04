<script lang="ts">
  import { onMount } from 'svelte';
  import { useUpdate } from '$lib/stores/update.svelte';
  import { t } from '$lib/stores/i18n.svelte';
  import { subscribeWsMessages } from '$lib/api/ws';
  import { cn } from '$lib/utils/cn';
  import Button from './Button.svelte';

  const update = useUpdate();

  onMount(() => {
    void update.check();

    const unsub = subscribeWsMessages((msg) => {
      const m = msg as { type?: string; latest?: string; tag?: string; releaseUrl?: string; assetName?: string; current?: string } | null;
      if (m?.type === 'update_available') {
        update.set({
          current: m.current ?? update.info?.current ?? '',
          latest: m.latest ?? '',
          tag: m.tag ?? '',
          assetName: m.assetName ?? '',
          releaseUrl: m.releaseUrl ?? '',
          hasUpdate: true,
        });
      }
    });
    return unsub;
  });
</script>

{#if update.info?.hasUpdate}
  <div
    class={cn(
      "ftm-enter mb-3 flex flex-wrap items-center gap-2 rounded-panel border border-primary/30 bg-primary/15 px-3 py-2 text-sm text-text-heading",
    )}
  >
    <div class="flex items-center gap-3 flex-wrap">
      <span class="font-semibold text-primary">
        ↑ {t('update_web_banner', { latest: update.info.latest })}
      </span>
      <a
        href={update.info.releaseUrl}
        target="_blank"
        rel="noreferrer"
        class="text-text-muted hover:text-primary underline"
      >
        {t('update_web_notes')}
      </a>
    </div>

    <div class="ml-auto">
      {#if update.applying}
        <span class="text-text-muted">{t('update_applying')}</span>
      {:else}
        <Button
          variant="primary"
          size="sm"
          onclick={() => update.apply()}
        >
          {t('update_web_button')}
        </Button>
      {/if}
    </div>
  </div>
{/if}

{#if update.error}
  <div
    class={cn(
      "ftm-enter mb-3 flex items-center gap-2 rounded-panel border border-status-error/40 bg-status-error/10 px-3 py-2 text-sm text-status-error",
    )}
  >
    {t('update_apply_failed', { 0: update.error })}
  </div>
{/if}
