<script lang="ts">
  import { t } from "$lib/stores/i18n.svelte";
  import { cn } from "$lib/utils/cn";
  import { STATUS_FILL, statusInfo } from "$lib/utils/status";

  let {
    state,
    percent = null,
    stale = false,
  }: {
    state: string;
    percent?: number | null;
    stale?: boolean;
  } = $props();

  const info = $derived(statusInfo(state));
  const idle = $derived(info.key === "stopped" || stale);
</script>

{#if idle}
  <span class="text-2xs font-medium text-text-muted">{t(info.textKey)}</span>
{:else}
  <span
    class={cn(
      "inline-flex items-center gap-1.5 rounded-control px-2 py-0.5 text-2xs font-medium text-btn-text",
      STATUS_FILL[info.key],
    )}
  >
    <span>{t(info.textKey)}</span>
    {#if percent !== null}
      <span class="font-semibold">{percent}%</span>
    {/if}
  </span>
{/if}
