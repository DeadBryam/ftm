<script lang="ts">
  import { Copy, Radio } from "lucide-svelte";
  import { useTunnels } from "$lib/stores/tunnels.svelte";
  import { useToast } from "$lib/stores/toast.svelte";
  import { t } from "$lib/stores/i18n.svelte";
  import { cn } from "$lib/utils/cn";

  const store = useTunnels();
  const toast = useToast();

  const ACTIVE = ["online", "starting", "connecting"];

  const online = $derived(
    store.tunnels.filter((tunnel) => ACTIVE.includes(tunnel.state)),
  );
  const shareUrl = $derived(
    online.find((tunnel) => tunnel.publicUrl)?.publicUrl ?? "",
  );

  async function copy() {
    if (!shareUrl) return;
    await navigator.clipboard.writeText(shareUrl);
    toast.success(t("overview_copied"));
  }
</script>

<section
  class="mb-3 flex shrink-0 flex-wrap items-center gap-x-4 gap-y-2 rounded-card border border-border bg-card px-3 py-2"
>
  <div class="flex items-center gap-2">
    <span
      class={cn(
        "flex h-2 w-2 rounded-full",
        online.length > 0 ? "bg-status-running" : "bg-status-stopped",
      )}
    ></span>
    <span class="text-sm font-semibold text-text-heading">
      {online.length}
      <span class="font-normal text-text-muted">{t("overview_online")}</span>
    </span>
    <span class="text-text-muted">·</span>
    <span class="text-sm text-text-muted">
      {store.tunnels.length}
      {t("overview_total")}
    </span>
  </div>

  {#if shareUrl}
    <button
      type="button"
      onclick={copy}
      class="flex min-w-0 flex-1 basis-64 cursor-pointer items-center gap-2 rounded-control bg-url-bg px-2.5 py-1.5 text-left transition-colors hover:bg-hover"
      title={t("overview_share")}
    >
      <Radio size={14} class="shrink-0 text-status-running" />
      <span
        class="min-w-0 flex-1 truncate font-mono text-xs text-url-text">{shareUrl}</span
      >
      <Copy size={13} class="shrink-0 text-text-muted" />
    </button>
  {:else}
    <span class="flex-1 basis-64 truncate text-xs text-text-muted">
      {t("overview_idle")}
    </span>
  {/if}
</section>
