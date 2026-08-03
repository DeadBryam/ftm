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

<div class="max-w-[1000px] mx-auto p-10 min-h-dvh flex max-md:max-h-dvh">
  <!--
    Held until the catalogue arrives. Rendering first meant the opening frame
    showed raw keys -- "port", "start", "connections" -- which then flipped to
    real text. `ready` is set even if the request fails, so this cannot hang.
  -->
  {#if i18n.ready}
    {@render children()}
  {/if}
</div>

<Toasts />
