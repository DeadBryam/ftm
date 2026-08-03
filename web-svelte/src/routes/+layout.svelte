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
  class="mx-auto flex min-h-dvh w-full max-w-[var(--app-max-width)] flex-col px-4 py-6 sm:px-6 sm:py-8 lg:px-10"
>
  {#if i18n.ready}
    {@render children()}
  {:else}
    <div class="flex flex-1 flex-col gap-6" aria-busy="true" aria-live="polite">
      <div class="flex items-center gap-4 border-b border-border pb-5">
        <div class="h-12 w-12 animate-pulse rounded-xl bg-border"></div>
        <div class="flex flex-1 flex-col gap-2">
          <div class="h-7 w-48 max-w-full animate-pulse rounded-md bg-border"></div>
          <div class="h-4 w-64 max-w-full animate-pulse rounded-md bg-border/70"></div>
        </div>
      </div>
      <div class="grid flex-1 gap-5 lg:grid-cols-[minmax(0,22rem)_1fr]">
        <div class="h-64 animate-pulse rounded-[var(--radius-panel)] bg-card border border-border"></div>
        <div class="h-64 animate-pulse rounded-[var(--radius-card)] bg-card border border-border"></div>
      </div>
    </div>
  {/if}
</div>

<Toasts />
