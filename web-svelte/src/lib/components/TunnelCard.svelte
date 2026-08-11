<script lang="ts">
  import { useToast } from "$lib/stores/toast.svelte";
  import { useProviders } from "$lib/stores/providers.svelte";
  import { useTunnels } from "$lib/stores/tunnels.svelte";
  import { AlertCircle, Check, Copy, Pause, Play } from "lucide-svelte";
  import { t } from "$lib/stores/i18n.svelte";
  import { useClock } from "$lib/stores/clock.svelte";
  import { cn } from "$lib/utils/cn";
  import { formatDuration } from "$lib/utils/duration";
  import { providerErrorHint } from "$lib/utils/providerError";
  import { copyText } from "$lib/utils/clipboard";
  import { isInstallingState, isRunningState } from "$lib/utils/status";
  import { onMount } from "svelte";
  import Button from "./Button.svelte";
  import StatusBadge from "./StatusBadge.svelte";
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

  const isRunning = $derived(isRunningState(tunnelState));

  const isInstalling = $derived(isInstallingState(tunnelState));

  const providerLabel = $derived(
    providerStore.providers.find((p) => p.id === tunnel.provider)?.name ??
      tunnel.provider,
  );

  const store = useTunnels();
  const stale = $derived(isRunning && !store.connected);

  const uptime = $derived(
    tunnel.startedAt && isRunning && !stale
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

  let copied = $state(false);
  let copiedTimer: ReturnType<typeof setTimeout> | null = null;

  async function copyUrl(url: string) {
    if (!(await copyText(url))) {
      toast.error(t("copy_failed"));
      return;
    }

    toast.info(t("copied"));

    copied = true;
    if (copiedTimer) clearTimeout(copiedTimer);
    copiedTimer = setTimeout(() => {
      copied = false;
    }, 3000);
  }

  const installPercent = $derived(
    Math.trunc((installProgress?.percent ?? 0) * 100) / 100,
  );
  const installStep = $derived(installProgress?.step ?? t("installing"));
  const errorHint = $derived(providerErrorHint(tunnel.errorMessage));
</script>

<div
  class={cn(
    "rounded-card border bg-card transition-colors duration-150",
    selected
      ? "border-primary bg-primary/5"
      : "border-border hover:border-primary/40",
  )}
>
  <div class="flex flex-col">
    <div
      class="flex flex-row items-start justify-between gap-2 p-2.5 sm:items-stretch"
    >
      <button
        type="button"
        aria-pressed={selected}
        onclick={() => onAction("select", tunnel.id)}
        class="flex min-w-0 flex-1 cursor-pointer gap-2.5 text-left"
      >
        <span
          aria-hidden="true"
          class="mt-0.5 flex h-6 w-6 shrink-0 items-center justify-center rounded-control bg-secondary-btn font-mono text-xs font-semibold text-secondary-btn-text uppercase"
        >
          {providerLabel.slice(0, 2)}
        </span>
        <span class="min-w-0 flex-1">
        <span
          class="mb-0.5 block overflow-hidden text-ellipsis whitespace-nowrap text-sm font-semibold"
        >
          {tunnel.name}
        </span>
        <span class="mb-1.5 block text-xs text-text-muted">
          {providerLabel} · <span class="font-mono">localhost:{tunnel.port}</span>
        </span>
        <span class="flex flex-wrap items-center gap-x-2 gap-y-1">
          <StatusBadge
            state={tunnelState}
            {stale}
            percent={tunnelState === "installing" && installProgress
              ? installPercent
              : null}
          />
          {#if stale}
            <span class="text-xs text-text-muted">{t("connection_lost")}</span>
          {/if}
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
        </span>
        </span>
      </button>
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

    {#if tunnelState === "installing" && installProgress}
      <div class="px-2.5 pb-2.5">
        <div class="h-1 w-full overflow-hidden rounded bg-border">
          <div
            class="h-full rounded bg-status-installing"
            style="width: {installPercent}%"
          ></div>
        </div>
        <div
          class="mt-1 overflow-hidden text-ellipsis whitespace-nowrap text-xs text-text-muted"
        >
          {installStep}
        </div>
      </div>
    {/if}

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
        <span class="flex h-4 w-4 shrink-0 text-text-muted">
          {#if copied}<Check size={14} />{:else}<Copy size={14} />{/if}
        </span>
        <span
          class="min-w-0 flex-1 overflow-hidden text-ellipsis whitespace-nowrap text-start font-mono text-xs text-primary"
          >{tunnel.publicUrl}</span
        >
        <span class="shrink-0 text-xs text-text-muted">
          {copied ? t("copied") : t("click_to_copy")}
        </span>
      </button>
    {/if}

    {#if tunnel.errorMessage}
      <div
        class="flex items-start gap-2 rounded-b-card border-t border-t-status-error/70 bg-status-error/15 px-2.5 py-2 text-status-error"
      >
        <span class="mt-0.5 h-4 w-4 shrink-0"><AlertCircle size={16} /></span>
        <span class="min-w-0 flex-1">
          {#if errorHint}
            <span class="block text-xs font-medium">{t(errorHint)}</span>
          {/if}
          <span class="block font-mono text-xs break-words opacity-80">
            {tunnel.errorMessage}
          </span>
        </span>
      </div>
    {/if}

  </div>
</div>
