<script lang="ts">
  import { onMount } from "svelte";
  import { Copy, Pause, Play, Unplug } from "lucide-svelte";
  import { useTunnels } from "$lib/stores/tunnels.svelte";
  import { useProviders } from "$lib/stores/providers.svelte";
  import { useToast } from "$lib/stores/toast.svelte";
  import { useClock } from "$lib/stores/clock.svelte";
  import { t } from "$lib/stores/i18n.svelte";
  import { cn } from "$lib/utils/cn";
  import { formatDuration } from "$lib/utils/duration";
  import {
    isInstallingState,
    isRunningState,
    statusColors as statusColorsFor,
    statusInfo as statusInfoFor,
  } from "$lib/utils/status";
  import Button from "./Button.svelte";
  import type { TunnelState } from "$lib/types";

  let { tunnelId }: { tunnelId: string } = $props();

  const store = useTunnels();
  const providerStore = useProviders();
  const toast = useToast();
  const clock = useClock();

  onMount(() => {
    store.connect();
    if (providerStore.providers.length === 0) providerStore.fetch();
    return clock.subscribe();
  });

  const tunnel = $derived(store.getById(tunnelId));
  const tunnelState = $derived((tunnel?.state ?? "stopped") as TunnelState);
  const statusInfo = $derived(statusInfoFor(tunnelState));
  const statusColors = $derived(statusColorsFor(tunnelState));
  const isRunning = $derived(!!tunnel && isRunningState(tunnelState));
  const isInstalling = $derived(!!tunnel && isInstallingState(tunnelState));

  const providerLabel = $derived(
    providerStore.providers.find((p) => p.id === tunnel?.provider)?.name ??
      tunnel?.provider ??
      "",
  );

  const uptime = $derived(
    tunnel?.startedAt && isRunning
      ? formatDuration(clock.now - tunnel.startedAt)
      : "",
  );

  const remaining = $derived(
    tunnel?.expiresAt && isRunning ? tunnel.expiresAt - clock.now : 0,
  );

  const actionLabel = $derived(
    isInstalling
      ? t("wait")
      : tunnelState === "stopping"
        ? t("stopping")
        : t("stop"),
  );

  async function copyUrl() {
    if (!tunnel?.publicUrl) return;
    await navigator.clipboard.writeText(tunnel.publicUrl);
    toast.success(t("overview_copied"));
  }

  async function toggle() {
    if (!tunnel) return;
    try {
      if (isRunning) {
        await store.stop(tunnel.id);
      } else {
        await store.start(tunnel.id);
      }
    } catch (err) {
      toast.error((err as Error).message);
    }
  }
</script>

<div class="flex h-dvh min-h-0 flex-col gap-2 bg-bg p-2.5 text-text">
  {#if !tunnel}
    <div
      class="flex h-full flex-col items-center justify-center gap-2 text-center text-text-muted"
    >
      <Unplug size={22} />
      <p class="m-0 text-xs leading-relaxed">{t("pip_gone")}</p>
    </div>
  {:else}
    <div class="min-w-0">
      <h1
        class="m-0 truncate text-sm font-semibold text-text-heading"
        title={tunnel.name}
      >
        {tunnel.name}
      </h1>
      <p class="m-0 truncate text-[11px] text-text-muted">
        {providerLabel} · <span class="font-mono">localhost:{tunnel.port}</span>
      </p>
    </div>

    <div class="flex flex-wrap items-center gap-x-2 gap-y-1">
      <span
        class={cn(
          "inline-flex items-center gap-1.5 rounded-control px-2 py-0.5 text-[11px] font-medium",
          statusColors.bg,
          statusColors.text,
        )}
      >
        <span class={cn("h-1.5 w-1.5 rounded-full", statusColors.dot)}></span>
        {t(statusInfo.textKey)}
      </span>
      {#if uptime}
        <span class="font-mono text-[11px] text-text-muted">{uptime}</span>
      {/if}
      {#if remaining > 0}
        <span class="font-mono text-[11px] text-text-muted">
          {t("card_expires", { 0: formatDuration(remaining) })}
        </span>
      {/if}
    </div>

    {#if tunnelState === "online"}
      <div class="flex gap-2 text-[11px] text-text-muted">
        <span>
          {t("detail_sessions")}
          <span class="font-mono text-text">{tunnel.activeSessions ?? 0}</span>
        </span>
        <span>
          {t("detail_visitors")}
          <span class="font-mono text-text">{tunnel.visitors ?? 0}</span>
        </span>
      </div>
    {/if}

    {#if tunnel.publicUrl}
      <button
        type="button"
        onclick={copyUrl}
        class="flex w-full cursor-pointer items-center gap-2 rounded-control bg-url-bg px-2 py-1.5 text-left transition-colors hover:bg-hover"
      >
        <span class="min-w-0 flex-1 truncate font-mono text-[11px] text-url-text"
          >{tunnel.publicUrl}</span
        >
        <Copy size={12} class="shrink-0 text-text-muted" />
      </button>
    {/if}

    {#if tunnel.errorMessage}
      <p
        class="m-0 line-clamp-2 rounded-control border border-status-error/40 bg-status-error/10 px-2 py-1 font-mono text-[11px] break-words text-status-error"
      >
        {tunnel.errorMessage}
      </p>
    {/if}

    <div class="mt-auto">
      <Button
        variant={isRunning ? "error" : "success"}
        icon={isRunning ? Pause : Play}
        class="w-full"
        disabled={isInstalling || tunnelState === "stopping"}
        onclick={toggle}
      >
        {isRunning ? actionLabel : t("start")}
      </Button>
    </div>
  {/if}
</div>
