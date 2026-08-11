<script lang="ts">
  import { onDestroy, onMount } from "svelte";
  import { browser } from "$app/environment";
  import { onNavigate } from "$app/navigation";
  import Toasts from "$lib/components/Toasts.svelte";
  import { subscribeWsMessages } from "$lib/api/ws";
  import { useI18n } from "$lib/stores/i18n.svelte";
  import { useTheme } from "$lib/stores/theme.svelte";
  import { prefersReducedMotion } from "$lib/utils/motion";

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

{@render children()}

<Toasts />
