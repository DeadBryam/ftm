<script lang="ts">
  import { onDestroy, onMount } from "svelte";
  import { browser } from "$app/environment";
  import { onNavigate } from "$app/navigation";
  import Header from "$lib/components/Header.svelte";
  import Toasts from "$lib/components/Toasts.svelte";
  import { subscribeWsMessages } from "$lib/api/ws";
  import { useI18n } from "$lib/stores/i18n.svelte";
  import { useTheme } from "$lib/stores/theme.svelte";

  import "../styles/app.css";

  let { children } = $props();

  const i18n = useI18n();
  const theme = useTheme();

  if (browser) theme.init();

  onNavigate((navigation) => {
    if (!document.startViewTransition || prefersReducedMotion()) return;
    return new Promise((resolve) => {
      document.startViewTransition(async () => {
        resolve();
        await navigation.complete;
      });
    });
  });

  function prefersReducedMotion(): boolean {
    return window.matchMedia("(prefers-reduced-motion: reduce)").matches;
  }

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
  class="mx-auto flex h-dvh w-full max-w-app flex-col overflow-hidden px-app-x py-app-y sm:px-5"
>
  {#if i18n.ready}
    <Header />
    <div class="min-h-0 flex-1">
      {@render children()}
    </div>
  {:else}
    <div class="flex min-h-0 flex-1 flex-col gap-3" aria-busy="true" aria-live="polite">
      <div class="flex items-center gap-3 border-b border-border pb-3">
        <div class="h-9 w-9 animate-pulse rounded-control bg-border"></div>
        <div class="flex flex-1 flex-col gap-1.5">
          <div class="h-5 w-40 max-w-full animate-pulse rounded-sm bg-border"></div>
          <div class="h-3 w-56 max-w-full animate-pulse rounded-sm bg-border/70"></div>
        </div>
      </div>
      <div class="grid min-h-0 flex-1 gap-app md:grid-cols-[minmax(0,18rem)_1fr]">
        <div class="min-h-0 rounded-panel border border-border bg-card"></div>
        <div class="min-h-0 rounded-card border border-border bg-card"></div>
      </div>
    </div>
  {/if}
</div>

<Toasts />
