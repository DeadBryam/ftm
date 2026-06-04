<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { updateStore } from '$lib/stores/update.svelte';
  import { translate } from '$lib/i18n';
  import { subscribeWsMessages } from '$lib/api/ws';

  let t = $derived($translate);

  onMount(() => {
    void updateStore.check();

    const unsub = subscribeWsMessages((msg) => {
      const m = msg as { type?: string; latest?: string; tag?: string; releaseUrl?: string; assetName?: string; current?: string } | null;
      if (m?.type === 'update_available') {
        const info = $updateStore.info;
        updateStore.set({
          current: m.current ?? info?.current ?? '',
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

{#if $updateStore.info?.hasUpdate}
  <div
    class="flex flex-wrap items-center gap-3 px-4 py-3 mb-5 rounded-lg border border-primary/30 bg-primary/10 text-sm"
  >
    <span class="font-semibold text-primary">
      ↑ {t('update_web_banner', { latest: $updateStore.info.latest })}
    </span>

    {#if $updateStore.applying}
      <span class="text-text-muted">{t('update_applying')}</span>
    {:else}
      <button
        class="px-3 py-1.5 rounded-md bg-primary text-white font-medium hover:bg-primary/90 transition-colors"
        onclick={() => updateStore.apply()}
      >
        {t('update_web_button')}
      </button>
    {/if}

    <a
      href={$updateStore.info.releaseUrl}
      target="_blank"
      rel="noreferrer"
      class="text-text-muted hover:text-primary underline"
    >
      {t('update_web_notes')}
    </a>
  </div>
{/if}

{#if $updateStore.error}
  <div
    class="flex items-center gap-3 px-4 py-3 mb-5 rounded-lg border border-red-300 bg-red-50 text-sm text-red-700"
  >
    {t('update_apply_failed', { 0: $updateStore.error })}
  </div>
{/if}
