<script lang="ts">
  import { onDestroy, onMount } from "svelte";
  import Toasts from "$lib/components/Toasts.svelte";
  import { subscribeWsMessages } from "$lib/api/ws";
  import { useI18n } from "$lib/stores/i18n.svelte";

  import "../styles/app.css";

  let { children } = $props();

  const i18n = useI18n();

  let unsubscribeWs: (() => void) | null = null;

  onMount(async () => {
    unsubscribeWs = subscribeWsMessages(() => {});
    await i18n.init();
  });

  onDestroy(() => {
    if (unsubscribeWs) {
      unsubscribeWs();
      unsubscribeWs = null;
    }
  });
</script>

<div
  class="mx-auto flex h-dvh w-full max-w-[var(--app-max-width)] flex-col overflow-hidden px-[var(--app-pad-x)] py-[var(--app-pad-y)] sm:px-5"
>
  {#if i18n.ready}
    {@render children()}
  {:else}
    <div class="flex min-h-0 flex-1 flex-col gap-3" aria-busy="true" aria-live="polite">
      <div class="flex items-center gap-3 border-b border-border pb-3">
        <div class="h-9 w-9 animate-pulse rounded-[var(--radius-control)] bg-border"></div>
        <div class="flex flex-1 flex-col gap-1.5">
          <div class="h-5 w-40 max-w-full animate-pulse rounded-sm bg-border"></div>
          <div class="h-3 w-56 max-w-full animate-pulse rounded-sm bg-border/70"></div>
        </div>
      </div>
      <div class="grid min-h-0 flex-1 gap-[var(--app-gap)] md:grid-cols-[minmax(0,18rem)_1fr]">
        <div class="min-h-0 rounded-[var(--radius-panel)] border border-border bg-card"></div>
        <div class="min-h-0 rounded-[var(--radius-card)] border border-border bg-card"></div>
      </div>
    </div>
  {/if}
</div>

<Toasts />
