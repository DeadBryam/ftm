<script lang="ts">
  import { onMount } from "svelte";
  import {
    Check,
    Copy,
    Pause,
    Play,
    QrCode as QrCodeIcon,
    Unplug,
  } from "lucide-svelte";
  import { useTunnels } from "$lib/stores/tunnels.svelte";
  import { useProviders } from "$lib/stores/providers.svelte";
  import { useToast } from "$lib/stores/toast.svelte";
  import { useClock } from "$lib/stores/clock.svelte";
  import { t } from "$lib/stores/i18n.svelte";
  import { cn } from "$lib/utils/cn";
  import { formatDuration } from "$lib/utils/duration";
  import { isInstallingState, isRunningState } from "$lib/utils/status";
  import { providerErrorHint } from "$lib/utils/providerError";
  import { copyText } from "$lib/utils/clipboard";
  import Button from "./Button.svelte";
  import QrCode from "./QrCode.svelte";
  import StatusBadge from "./StatusBadge.svelte";
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
  const isRunning = $derived(!!tunnel && isRunningState(tunnelState));
  const isInstalling = $derived(!!tunnel && isInstallingState(tunnelState));
  const stale = $derived(isRunning && !store.connected);
  const errorHint = $derived(providerErrorHint(tunnel?.errorMessage));

  const providerLabel = $derived(
    providerStore.providers.find((p) => p.id === tunnel?.provider)?.name ??
      tunnel?.provider ??
      "",
  );

  const uptime = $derived(
    tunnel?.startedAt && isRunning && !stale
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

  let copied = $state(false);
  let copiedTimer: ReturnType<typeof setTimeout> | null = null;

  async function copyUrl() {
    if (!tunnel?.publicUrl) return;

    if (!(await copyText(tunnel.publicUrl))) {
      toast.error(t("copy_failed"));
      return;
    }

    toast.success(t("overview_copied"));

    copied = true;
    if (copiedTimer) clearTimeout(copiedTimer);
    copiedTimer = setTimeout(() => {
      copied = false;
    }, 3000);
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

<div class="panel flex h-dvh min-h-0 flex-col gap-2.5 bg-bg p-3 text-text">
  {#if !tunnel}
    <div
      class="flex h-full flex-col items-center justify-center gap-2 text-center text-text-muted"
    >
      <Unplug size={22} />
      <p class="m-0 text-xs leading-relaxed">{t("pip_gone")}</p>
    </div>
  {:else}
    <div class="flex shrink-0 items-start justify-between gap-2">
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
      <span class="shrink-0">
        <StatusBadge state={tunnelState} {stale} />
      </span>
    </div>

    <dl
      class="m-0 grid shrink-0 grid-cols-3 divide-x divide-border-light rounded-control border border-border-light bg-card"
    >
      <div class="min-w-0 px-2 py-1.5">
        <dt class="m-0 truncate text-[10px] text-text-muted">
          {t("detail_uptime")}
        </dt>
        <dd class="m-0 truncate font-mono text-xs text-text-heading">
          {uptime || "—"}
        </dd>
      </div>
      <div class="min-w-0 px-2 py-1.5">
        <dt class="m-0 truncate text-[10px] text-text-muted">
          {t("detail_sessions")}
        </dt>
        <dd class="m-0 font-mono text-xs text-text-heading">
          {tunnelState === "online" ? (tunnel.activeSessions ?? 0) : "—"}
        </dd>
      </div>
      <div class="min-w-0 px-2 py-1.5">
        <dt class="m-0 truncate text-[10px] text-text-muted">
          {remaining > 0 ? t("detail_expires") : t("detail_visitors")}
        </dt>
        <dd class="m-0 truncate font-mono text-xs text-text-heading">
          {#if remaining > 0}
            {formatDuration(remaining)}
          {:else}
            {tunnelState === "online" ? (tunnel.visitors ?? 0) : "—"}
          {/if}
        </dd>
      </div>
    </dl>

    {#if tunnel.publicUrl}
      <button
        type="button"
        onclick={copyUrl}
        class="flex w-full shrink-0 cursor-pointer items-center gap-2 rounded-control bg-url-bg px-2.5 py-2 text-left transition-colors hover:bg-hover"
      >
        <span class="min-w-0 flex-1 truncate font-mono text-[11px] text-url-text"
          >{tunnel.publicUrl}</span
        >
        {#if copied}
          <Check size={12} class="shrink-0 text-status-running" />
        {:else}
          <Copy size={12} class="shrink-0 text-text-muted" />
        {/if}
      </button>
    {/if}

    {#if tunnel.errorMessage}
      <p
        class="m-0 line-clamp-2 shrink-0 rounded-control border border-status-error/40 bg-status-error/10 px-2 py-1 text-[11px] break-words text-status-error"
        title={tunnel.errorMessage}
      >
        {errorHint ? t(errorHint) : tunnel.errorMessage}
      </p>
    {/if}

    <div class="fill flex min-h-0 flex-1 flex-col items-center justify-center gap-1.5">
      {#if tunnel.publicUrl}
        <QrCode
          value={tunnel.publicUrl}
          size={176}
          class="h-full min-h-0 w-auto max-w-full object-contain"
        />
        <p class="m-0 shrink-0 text-center text-[11px] text-text-muted">
          {t("overview_share")}
        </p>
      {:else}
        <div
          class="flex aspect-square h-full min-h-0 items-center justify-center rounded-control border border-dashed border-border text-text-muted/50"
        >
          <QrCodeIcon size={32} strokeWidth={1.5} />
        </div>
        <p class="m-0 max-w-[24ch] shrink-0 text-center text-[11px] leading-relaxed text-text-muted">
          {t("pip_idle")}
        </p>
      {/if}
    </div>

    <Button
      variant={isRunning ? "error" : "success"}
      icon={isRunning ? Pause : Play}
      size="md"
      class="w-full shrink-0"
      disabled={isInstalling || tunnelState === "stopping"}
      onclick={toggle}
    >
      {isRunning ? actionLabel : t("start")}
    </Button>
  {/if}
</div>

<style>
  .panel {
    container: panel / size;
  }

  @container panel (height < 300px) {
    .fill > :global(*) {
      display: none;
    }
  }
</style>
