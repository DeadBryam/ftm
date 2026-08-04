<script lang="ts">
  import { Copy, MousePointerClick } from "lucide-svelte";
  import { useProviders } from "$lib/stores/providers.svelte";
  import { useToast } from "$lib/stores/toast.svelte";
  import { useClock } from "$lib/stores/clock.svelte";
  import { t } from "$lib/stores/i18n.svelte";
  import { cn } from "$lib/utils/cn";
  import { formatDuration } from "$lib/utils/duration";
  import QrCode from "./QrCode.svelte";
  import type { Tunnel } from "$lib/types";
  import { onMount } from "svelte";

  let { tunnel }: { tunnel: Tunnel | null } = $props();

  const providerStore = useProviders();
  const toast = useToast();
  const clock = useClock();

  onMount(() => clock.subscribe());

  const providerLabel = $derived(
    providerStore.providers.find((p) => p.id === tunnel?.provider)?.name ??
      tunnel?.provider ??
      "",
  );

  const isOnline = $derived(tunnel?.state === "online");

  const uptime = $derived(
    tunnel?.startedAt && isOnline
      ? formatDuration(clock.now - tunnel.startedAt)
      : "",
  );

  const remaining = $derived(
    tunnel?.expiresAt && isOnline ? tunnel.expiresAt - clock.now : 0,
  );

  async function copy() {
    if (!tunnel?.publicUrl) return;
    await navigator.clipboard.writeText(tunnel.publicUrl);
    toast.success(t("overview_copied"));
  }
</script>

<section
  class="ftm-enter flex h-full min-h-0 flex-col overflow-hidden rounded-card border border-border bg-card"
>
  <div
    class="flex shrink-0 items-center justify-between border-b border-border-light bg-url-bg px-3 py-2"
  >
    <h2 class="m-0 text-sm font-semibold text-text-heading">
      {t("detail_title")}
    </h2>
  </div>

  <div class="min-h-0 flex-1 overflow-y-auto p-3">
    {#if !tunnel}
      <div
        class="flex h-full flex-col items-center justify-center gap-2 px-3 text-center text-text-muted"
      >
        <MousePointerClick size={26} />
        <p class="m-0 text-xs leading-relaxed">{t("detail_empty")}</p>
      </div>
    {:else}
      <h3
        class="m-0 mb-1 truncate text-base font-semibold text-text-heading"
        title={tunnel.name}
      >
        {tunnel.name}
      </h3>
      <p class="m-0 mb-3 text-xs text-text-muted">
        {providerLabel} · <span class="font-mono">localhost:{tunnel.port}</span>
      </p>

      <dl class="m-0 mb-3 grid grid-cols-[auto_1fr] gap-x-3 gap-y-1 text-xs">
        <dt class="text-text-muted">{t("status_label")}</dt>
        <dd
          class={cn(
            "m-0 font-medium",
            isOnline ? "text-status-running" : "text-text",
          )}
        >
          {t(tunnel.state)}
        </dd>
        {#if uptime}
          <dt class="text-text-muted">{t("detail_uptime")}</dt>
          <dd class="m-0 font-mono text-text">{uptime}</dd>
        {/if}
        {#if remaining > 0}
          <dt class="text-text-muted">{t("detail_expires")}</dt>
          <dd class="m-0 font-mono text-text">{formatDuration(remaining)}</dd>
        {/if}
      </dl>

      {#if tunnel.publicUrl}
        <button
          type="button"
          onclick={copy}
          class="mb-3 flex w-full cursor-pointer items-center gap-2 rounded-control bg-url-bg px-2.5 py-2 text-left transition-colors hover:bg-hover"
        >
          <span class="min-w-0 flex-1 truncate font-mono text-xs text-url-text"
            >{tunnel.publicUrl}</span
          >
          <Copy size={13} class="shrink-0 text-text-muted" />
        </button>

        <div class="flex flex-col items-center gap-2">
          <QrCode value={tunnel.publicUrl} size={148} />
          <p class="m-0 text-center text-xs leading-relaxed text-text-muted">
            {t("overview_share")}
          </p>
        </div>
      {/if}

      {#if tunnel.errorMessage}
        <p
          class="m-0 rounded-control border border-status-error/40 bg-status-error/10 px-2.5 py-1.5 font-mono text-xs break-words text-status-error"
        >
          {tunnel.errorMessage}
        </p>
      {/if}
    {/if}
  </div>
</section>
