<script lang="ts">
  import { useToast } from "$lib/stores/toast.svelte";
  import { useProviders } from "$lib/stores/providers.svelte";
  import {
    AlertCircle,
    Copy,
    FileText,
    Menu,
    Pause,
    Pencil,
    Play,
    Trash2,
    X,
  } from "lucide-svelte";
  import { logsApi } from "$lib/api";
  import { t } from "$lib/stores/i18n.svelte";
  import { cn } from "$lib/utils/cn";
  import { formatLogs } from "$lib/utils/logs";
  import Button from "./Button.svelte";
  import Dropdown from "./Dropdown.svelte";
  import type {
    DropdownOption,
    LogStream,
    Tunnel,
    TunnelState,
  } from "$lib/types";

  type StatusKey = "running" | "starting" | "installing" | "error" | "stopped";
  type StatusColors = { bg: string; text: string; dot: string };
  type StatusInfo = { key: StatusKey; textKey: string };
  type InstallProgress = { percent: number; step: string };

  interface TunnelCardProps {
    tunnel: Tunnel;
    onStart: (id: string) => void;
    onStop: (id: string) => void;
    onAction: (action: string, data: unknown) => void;
    index?: number;
    totalItems?: number;
    installProgress?: InstallProgress | null;
  }

  let {
    tunnel,
    onStart,
    onStop,
    onAction,
    index = 0,
    totalItems = 1,
    installProgress = null,
  }: TunnelCardProps = $props();

  const dropdownAlign = $derived(
    index === totalItems - 1 && totalItems > 1 ? "top-left" : "left",
  );

  const toast = useToast();
  const providerStore = useProviders();

  let showLogs = $state(false);
  let logs = $state("");
  let loadingLogs = $state(false);
  let logStream: LogStream | null = $state(null);
  let logPre: HTMLPreElement | undefined = $state();
  let followBottom = $state(true);

  const statusConfig: Record<StatusKey, StatusColors> = {
    running: {
      bg: "bg-status-running/40",
      text: "text-status-running",
      dot: "bg-status-running/95",
    },
    starting: {
      bg: "bg-status-starting/40",
      text: "text-status-starting",
      dot: "bg-status-starting/95",
    },
    installing: {
      bg: "bg-status-installing/40",
      text: "text-status-installing",
      dot: "bg-status-installing/95",
    },
    error: {
      bg: "bg-status-error/40",
      text: "text-status-error",
      dot: "bg-status-error/95",
    },
    stopped: {
      bg: "bg-status-stopped/40",
      text: "text-status-stopped",
      dot: "bg-status-stopped/95",
    },
  };

  const statusMap: Record<TunnelState, StatusInfo> = {
    online: { key: "running", textKey: "online" },
    starting: { key: "starting", textKey: "starting" },
    connecting: { key: "starting", textKey: "connecting" },
    installing: { key: "installing", textKey: "installing" },
    downloading: { key: "installing", textKey: "downloading" },
    need_installing: { key: "stopped", textKey: "need_installing" },
    stopping: { key: "starting", textKey: "stopping" },
    stopped: { key: "stopped", textKey: "stopped" },
    offline: { key: "stopped", textKey: "offline" },
    timeout: { key: "error", textKey: "timeout" },
    error: { key: "error", textKey: "error" },
  };

  const tunnelState = $derived(tunnel.state as TunnelState);
  const statusInfo = $derived(statusMap[tunnelState] ?? statusMap.error);
  const statusColors = $derived(statusConfig[statusInfo.key]);

  const isRunning = $derived(
    tunnelState === "online" ||
      tunnelState === "starting" ||
      tunnelState === "connecting" ||
      tunnelState === "installing" ||
      tunnelState === "downloading" ||
      tunnelState === "stopping",
  );

  const isInstalling = $derived(
    tunnelState === "installing" || tunnelState === "downloading",
  );

  const providerLabel = $derived(
    providerStore.providers.find((p) => p.id === tunnel.provider)?.name ??
      tunnel.provider,
  );

  const actionLabel = $derived(
    isInstalling
      ? t("wait")
      : tunnelState === "stopping"
        ? t("stopping")
        : t("stop"),
  );

  function scrollLogsToBottom() {
    if (!logPre || !followBottom) return;
    logPre.scrollTop = logPre.scrollHeight;
  }

  function onLogScroll() {
    if (!logPre) return;
    const distance =
      logPre.scrollHeight - logPre.scrollTop - logPre.clientHeight;
    followBottom = distance < 24;
  }

  $effect(() => {
    logs;
    queueMicrotask(scrollLogsToBottom);
  });

  function copyUrl(url: string) {
    navigator.clipboard.writeText(url);
    toast.info(t("copied"));
  }

  function closeLogs() {
    if (logStream) {
      logStream.close();
      logStream = null;
    }
    loadingLogs = false;
    showLogs = false;
  }

  function loadLogs() {
    if (showLogs) {
      closeLogs();
      return;
    }

    showLogs = true;
    loadingLogs = true;
    logs = "";
    followBottom = true;

    logsApi
      .get(tunnel.id)
      .then((initial) => {
        logs = formatLogs(initial);
        loadingLogs = false;
      })
      .catch(() => {
        logs = t("error_loading_logs");
        loadingLogs = false;
      });

    logStream = logsApi.createStream(tunnel.id, {
      onLine: (line: string) => {
        logs = logs + "\n" + formatLogs(line);
      },
      onClose: () => {
        logStream = null;
      },
    });
  }

  function handleDropdownAction(option: DropdownOption) {
    switch (option.action) {
      case "edit":
        onAction("edit", tunnel.id);
        break;
      case "logs":
        loadLogs();
        break;
      case "delete":
        onAction("delete", tunnel);
        break;
    }
  }

  const dropdownOptions: DropdownOption[] = $derived([
    { label: t("edit"), action: "edit", icon: Pencil, disabled: isRunning },
    { label: t("logs"), action: "logs", icon: FileText },
    { label: "separator", action: "separator" },
    { label: t("delete"), action: "delete", icon: Trash2, danger: true },
  ]);

  const installPercent = $derived(
    Math.trunc((installProgress?.percent ?? 0) * 100) / 100,
  );
  const installStep = $derived(installProgress?.step ?? t("installing"));
</script>

<div
  class="cursor-default rounded-card border border-border bg-card transition-all duration-150"
>
  <div class="flex flex-col">
    <div
      class="flex flex-row items-start justify-between gap-2 p-2.5 sm:items-stretch"
    >
      <div class="min-w-0 flex-1">
        <div
          class="mb-0.5 overflow-hidden text-ellipsis whitespace-nowrap text-sm font-semibold"
        >
          {tunnel.name}
        </div>
        <div class="mb-1.5 text-xs text-muted">
          {providerLabel} · {t("port")} {tunnel.port}
        </div>
        <div
          class={cn(
            "inline-flex items-center gap-1.5 rounded-control px-2 py-0.5 text-xs font-medium",
            statusColors.bg,
            statusColors.text,
          )}
        >
          <span class={cn("h-1.5 w-1.5 rounded-full", statusColors.dot)}></span>
          <span>{t(statusInfo.textKey)}</span>
          {#if tunnelState === "installing" && installProgress}
            <span class="ml-1 font-semibold">{installPercent}%</span>
          {/if}
        </div>
        {#if tunnelState === "installing" && installProgress}
          <div class="mt-2 h-1 w-full overflow-hidden rounded bg-border">
            <div
              class="h-full rounded bg-status-installing"
              style="width: {installPercent}%"
            ></div>
          </div>
          <div
            class="mt-1 overflow-hidden text-ellipsis whitespace-nowrap text-[11px] text-muted sm:text-[10px]"
          >
            {installStep}
          </div>
        {/if}
      </div>
      <div class="relative flex shrink-0 gap-2">
        {#if isRunning}
          <Button
            variant="error"
            icon={Pause}
            onclick={() => onStop(tunnel.id)}
            disabled={isInstalling || tunnelState === "stopping"}
          >
            {actionLabel}
          </Button>
        {:else}
          <Button
            variant="success"
            icon={Play}
            onclick={() => onStart(tunnel.id)}
          >
            {t("start")}
          </Button>
        {/if}
        <Dropdown
          options={dropdownOptions}
          onSelect={handleDropdownAction}
          align={dropdownAlign}
          ariaLabel={t("tunnel_options")}
        >
          {#snippet children()}
            <Menu size={16} />
          {/snippet}
        </Dropdown>
      </div>
    </div>

    {#if tunnel.publicUrl}
      <button
        type="button"
        class={cn(
          "flex w-full cursor-pointer items-center gap-2 border-t border-t-status-stopped bg-url-bg px-2.5 py-2",
          "transition-colors hover:bg-hover",
          {
            "rounded-b-card": !(tunnel.errorMessage || showLogs),
          },
        )}
        onclick={() => tunnel.publicUrl && copyUrl(tunnel.publicUrl)}
      >
        <span class="flex h-4 w-4 shrink-0 text-muted"><Copy size={14} /></span>
        <span
          class="min-w-0 flex-1 overflow-hidden text-ellipsis whitespace-nowrap text-start font-mono text-xs text-primary"
          >{tunnel.publicUrl}</span
        >
        <span class="shrink-0 text-[10px] text-muted">{t("click_to_copy")}</span>
      </button>
    {/if}

    {#if tunnel.errorMessage}
      <div
        class={cn(
          "flex items-center gap-2 border-t border-t-status-error/70 bg-status-error/15 px-2.5 py-2 text-status-error",
          {
            "rounded-b-card": !showLogs,
          },
        )}
      >
        <span class="h-4 w-4 shrink-0"><AlertCircle size={16} /></span>
        <span class="font-mono text-xs break-words">{tunnel.errorMessage}</span>
      </div>
    {/if}

    {#if showLogs}
      <div class="overflow-hidden rounded-b-card bg-logs-bg">
        <div
          class="flex items-center justify-between border-b border-border px-2.5 py-1.5"
        >
          <span class="text-[11px] font-medium text-muted">{t("live_logs")}</span>
          <Button variant="ghost" size="sm" icon={X} onclick={closeLogs}>
            {t("close")}
          </Button>
        </div>
        {#if loadingLogs}
          <div
            class="flex items-center justify-center gap-3 p-6 text-status-stopped sm:gap-2.5 sm:p-4"
          >
            <span>{t("loading")}</span>
          </div>
        {:else}
          <pre
            bind:this={logPre}
            onscroll={onLogScroll}
            class="m-0 max-h-[300px] overflow-y-auto overflow-x-hidden whitespace-pre-wrap break-all p-4 font-mono text-[12px] leading-relaxed text-logs-text sm:p-3.5 sm:text-[11px]">{logs ||
              t("no_logs")}</pre>
        {/if}
      </div>
    {/if}
  </div>
</div>
