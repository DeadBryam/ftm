<script lang="ts">
  import { onMount } from 'svelte';
  import { useUpdate } from '$lib/stores/update.svelte';
  import { t } from '$lib/stores/i18n.svelte';
  import { subscribeWsMessages } from '$lib/api/ws';
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
    class="mb-3 flex flex-wrap items-center gap-2 rounded-panel border border-primary/30 bg-primary/10 px-3 py-2 text-sm"
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
        <button
          class="px-3 py-1.5 rounded-md bg-primary text-white font-medium hover:bg-primary/90 transition-colors cursor-pointer"
          onclick={() => update.apply()}
        >
          {t('update_web_button')}
        </button>
      {/if}
    </div>
  </div>
{/if}

{#if update.error}
  <div
    class="mb-3 flex items-center gap-2 rounded-panel border border-red-300 bg-red-50 px-3 py-2 text-sm text-red-700"
  >
    {t('update_apply_failed', { 0: update.error })}
  </div>
{/if}
