<script lang="ts">
  import {
    ChevronDown,
    Copy,
    MousePointerClick,
    Pencil,
    QrCode as QrCodeIcon,
    Trash2,
  } from "lucide-svelte";
  import { useProviders } from "$lib/stores/providers.svelte";
  import { useToast } from "$lib/stores/toast.svelte";
  import { useClock } from "$lib/stores/clock.svelte";
  import { logsApi } from "$lib/api";
  import { t } from "$lib/stores/i18n.svelte";
  import { cn } from "$lib/utils/cn";
  import { formatDuration } from "$lib/utils/duration";
  import { formatLogs } from "$lib/utils/logs";
  import Button from "./Button.svelte";
  import IconButton from "./IconButton.svelte";
  import QrCode from "./QrCode.svelte";
  import type { LogStream, Tunnel } from "$lib/types";
  import { onDestroy, onMount } from "svelte";

  let {
    tunnel,
    onAction,
  }: {
    tunnel: Tunnel | null;
    onAction: (action: string, data: unknown) => void;
  } = $props();

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

  const LOGS_OPEN_KEY = "ftm-detail-logs-open";

  let qrDataUrl = $state("");
  let logsOpen = $state(true);
  let logs = $state("");
  let logPre: HTMLPreElement | undefined = $state();
  let followBottom = $state(true);
  let stream: LogStream | null = null;

  const isRunning = $derived(
    ["online", "starting", "connecting"].includes(tunnel?.state ?? ""),
  );

  onMount(() => {
    logsOpen = localStorage.getItem(LOGS_OPEN_KEY) !== "false";
  });

  function toggleLogs() {
    logsOpen = !logsOpen;
    localStorage.setItem(LOGS_OPEN_KEY, String(logsOpen));
  }

  function closeStream() {
    stream?.close();
    stream = null;
  }

  onDestroy(closeStream);

  $effect(() => {
    const id = tunnel?.id;
    const running = isRunning;
    const open = logsOpen;

    closeStream();
    logs = "";
    followBottom = true;

    if (!id || !running || !open) return;

    logsApi
      .get(id)
      .then((initial) => {
        if (tunnel?.id === id) logs = formatLogs(initial);
      })
      .catch(() => {
        if (tunnel?.id === id) logs = t("error_loading_logs");
      });

    stream = logsApi.createStream(id, {
      onLine: (line) => {
        logs = logs + "\n" + formatLogs(line);
      },
      onClose: () => {
        stream = null;
      },
    });

    return closeStream;
  });

  $effect(() => {
    logs;
    queueMicrotask(() => {
      if (logPre && followBottom) logPre.scrollTop = logPre.scrollHeight;
    });
  });

  function onLogScroll() {
    if (!logPre) return;
    followBottom =
      logPre.scrollHeight - logPre.scrollTop - logPre.clientHeight < 24;
  }

  async function copy() {
    if (!tunnel?.publicUrl) return;
    await navigator.clipboard.writeText(tunnel.publicUrl);
    toast.success(t("overview_copied"));
  }

  function downloadQr(blob: Blob) {
    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = `${tunnel?.name ?? "tunnel"}-qr.png`;
    link.click();
    URL.revokeObjectURL(url);
    toast.success(t("qr_downloaded"));
  }

  async function copyQr() {
    if (!qrDataUrl) return;

    const blob = await (await fetch(qrDataUrl)).blob();

    if (typeof ClipboardItem === "undefined" || !navigator.clipboard?.write) {
      downloadQr(blob);
      return;
    }

    try {
      await navigator.clipboard.write([
        new ClipboardItem({ "image/png": blob }),
      ]);
      toast.success(t("qr_copied"));
    } catch {
      downloadQr(blob);
    }
  }
</script>

<section
  class="ftm-enter flex h-full min-h-0 flex-col overflow-hidden rounded-card border border-border bg-card"
>
  <div
    class="flex shrink-0 items-center justify-between gap-2 border-b border-border-light bg-url-bg px-3 py-2"
  >
    <h2 class="m-0 truncate text-sm font-semibold text-text-heading">
      {t("detail_title")}
    </h2>
    {#if tunnel}
      <div class="flex shrink-0 gap-1">
        <IconButton
          icon={Pencil}
          label={t("edit")}
          size="sm"
          disabled={isRunning}
          onclick={() => onAction("edit", tunnel.id)}
        />
        <IconButton
          icon={Trash2}
          label={t("delete")}
          variant="danger"
          size="sm"
          disabled={isRunning}
          onclick={() => onAction("delete", tunnel)}
        />
      </div>
    {/if}
  </div>

  <div class="relative min-h-0 flex-1 overflow-hidden">
    <div class="panel-pattern" aria-hidden="true"></div>

    <div class="relative z-10 h-full min-h-0 overflow-y-auto p-3">
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
        {#if isOnline}
          <dt class="text-text-muted">{t("detail_sessions")}</dt>
          <dd class="m-0 font-mono text-text">{tunnel.activeSessions ?? 0}</dd>
          <dt class="text-text-muted">{t("detail_visitors")}</dt>
          <dd class="m-0 font-mono text-text">{tunnel.visitors ?? 0}</dd>
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

        <div
          class="flex flex-col items-center gap-2 rounded-control border border-border-light bg-bg/40 p-3"
        >
          <QrCode value={tunnel.publicUrl} size={140} bind:dataUrl={qrDataUrl} />
          <p class="m-0 text-center text-xs leading-relaxed text-text-muted">
            {t("overview_share")}
          </p>
          <Button
            variant="default"
            icon={QrCodeIcon}
            disabled={!qrDataUrl}
            onclick={copyQr}
          >
            {t("qr_copy")}
          </Button>
        </div>
      {/if}

      {#if tunnel.errorMessage}
        <p
          class="m-0 mb-3 rounded-control border border-status-error/40 bg-status-error/10 px-2.5 py-1.5 font-mono text-xs break-words text-status-error"
        >
          {tunnel.errorMessage}
        </p>
      {/if}

      <div class="mt-3 border-t border-border-light pt-3">
        <button
          type="button"
          onclick={toggleLogs}
          aria-expanded={logsOpen}
          class="mb-1.5 flex w-full cursor-pointer items-center justify-between gap-2 text-xs font-medium text-text-muted transition-colors hover:text-text"
        >
          <span>{t("detail_logs")}</span>
          <span class="flex items-center gap-1">
            <span class="text-[11px] font-normal">
              {logsOpen ? t("detail_logs_hide") : t("detail_logs_show")}
            </span>
            <ChevronDown
              size={14}
              class={cn("transition-transform", !logsOpen && "-rotate-90")}
            />
          </span>
        </button>
        {#if logsOpen && !isRunning}
          <p class="m-0 text-xs text-text-muted italic">
            {t("detail_logs_idle")}
          </p>
        {:else if logsOpen}
          <pre
            bind:this={logPre}
            onscroll={onLogScroll}
            class="m-0 max-h-48 overflow-x-hidden overflow-y-auto rounded-control bg-logs-bg p-2 font-mono text-[11px] leading-relaxed break-all whitespace-pre-wrap text-logs-text">{logs ||
              t("no_logs")}</pre>
        {/if}
      </div>

    {/if}
    </div>
  </div>
</section>
