<script lang="ts">
  import { useToast } from "$lib/stores/toast.svelte";
  import { useProviders } from "$lib/stores/providers.svelte";
  import { AlertCircle, Copy, Pause, Play } from "lucide-svelte";
  import { t } from "$lib/stores/i18n.svelte";
  import { useClock } from "$lib/stores/clock.svelte";
  import { cn } from "$lib/utils/cn";
  import { formatDuration } from "$lib/utils/duration";
  import {
    isInstallingState,
    isRunningState,
    statusColors as statusColorsFor,
    statusInfo as statusInfoFor,
  } from "$lib/utils/status";
  import { onMount } from "svelte";
  import Button from "./Button.svelte";
  import type { Tunnel, TunnelState } from "$lib/types";

  type InstallProgress = { percent: number; step: string };

  interface TunnelCardProps {
    tunnel: Tunnel;
    onStart: (id: string) => void;
    onStop: (id: string) => void;
    onAction: (action: string, data: unknown) => void;
    index?: number;
    totalItems?: number;
    selected?: boolean;
    installProgress?: InstallProgress | null;
  }

  let {
    tunnel,
    onStart,
    onStop,
    onAction,
    index = 0,
    totalItems = 1,
    selected = false,
    installProgress = null,
  }: TunnelCardProps = $props();

  const toast = useToast();
  const providerStore = useProviders();
  const clock = useClock();

  onMount(() => clock.subscribe());

  const tunnelState = $derived(tunnel.state as TunnelState);
  const statusInfo = $derived(statusInfoFor(tunnelState));
  const statusColors = $derived(statusColorsFor(tunnelState));

  const isRunning = $derived(isRunningState(tunnelState));

  const isInstalling = $derived(isInstallingState(tunnelState));

  const providerLabel = $derived(
    providerStore.providers.find((p) => p.id === tunnel.provider)?.name ??
      tunnel.provider,
  );

  const uptime = $derived(
    tunnel.startedAt && isRunning
      ? formatDuration(clock.now - tunnel.startedAt)
      : "",
  );

  const remaining = $derived(
    tunnel.expiresAt && isRunning ? tunnel.expiresAt - clock.now : 0,
  );

  const expiryLabel = $derived(
    tunnel.expiresAt && isRunning
      ? remaining > 0
        ? t("card_expires", { 0: formatDuration(remaining) })
        : t("card_expired")
      : "",
  );

  const actionLabel = $derived(
    isInstalling
      ? t("wait")
      : tunnelState === "stopping"
        ? t("stopping")
        : t("stop"),
  );

  function copyUrl(url: string) {
    navigator.clipboard.writeText(url);
    toast.info(t("copied"));
  }

  const installPercent = $derived(
    Math.trunc((installProgress?.percent ?? 0) * 100) / 100,
  );
  const installStep = $derived(installProgress?.step ?? t("installing"));
</script>

<!-- svelte-ignore a11y_no_noninteractive_element_to_interactive_role -->
<div
  role="button"
  tabindex="0"
  aria-pressed={selected}
  onclick={() => onAction("select", tunnel.id)}
  onkeydown={(e) =>
    (e.key === "Enter" || e.key === " ") &&
    (e.preventDefault(), onAction("select", tunnel.id))}
  class={cn(
    "cursor-pointer rounded-card border bg-card transition-colors duration-150",
    selected
      ? "border-primary bg-primary/5"
      : "border-border hover:border-primary/40",
  )}
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
        <div class="mb-1.5 text-xs text-text-muted">
          {providerLabel} · <span class="font-mono">localhost:{tunnel.port}</span>
        </div>
        <div class="flex flex-wrap items-center gap-x-2 gap-y-1">
          <div
            class={cn(
              "inline-flex items-center gap-1.5 rounded-control px-2 py-0.5 text-xs font-medium",
              statusColors.bg,
              statusColors.text,
            )}
          >
            <span class={cn("h-1.5 w-1.5 rounded-full", statusColors.dot)}
            ></span>
            <span>{t(statusInfo.textKey)}</span>
            {#if tunnelState === "installing" && installProgress}
              <span class="ml-1 font-semibold">{installPercent}%</span>
            {/if}
          </div>
          {#if uptime}
            <span class="text-xs text-text-muted">
              {t("card_uptime", { 0: uptime })}
            </span>
          {/if}
          {#if expiryLabel}
            <span
              class={cn(
                "text-xs",
                remaining > 0 ? "text-text-muted" : "text-status-error",
              )}
            >
              {expiryLabel}
            </span>
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
      <div class="flex shrink-0 gap-2">
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
      </div>
    </div>

    {#if tunnel.publicUrl}
      <button
        type="button"
        class={cn(
          "flex w-full cursor-pointer items-center gap-2 border-t border-t-status-stopped bg-url-bg px-2.5 py-2",
          "transition-colors hover:bg-hover",
          { "rounded-b-card": !tunnel.errorMessage },
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
        class="flex items-center gap-2 rounded-b-card border-t border-t-status-error/70 bg-status-error/15 px-2.5 py-2 text-status-error"
      >
        <span class="h-4 w-4 shrink-0"><AlertCircle size={16} /></span>
        <span class="font-mono text-xs break-words">{tunnel.errorMessage}</span>
      </div>
    {/if}

  </div>
</div>
