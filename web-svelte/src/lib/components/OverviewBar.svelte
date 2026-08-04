<script lang="ts">
  import { onMount } from "svelte";
  import { Activity, Copy, LayoutDashboard, Radio } from "lucide-svelte";
  import { useTunnels } from "$lib/stores/tunnels.svelte";
  import { useToast } from "$lib/stores/toast.svelte";
  import { statusApi } from "$lib/api/endpoints/status";
  import { t } from "$lib/stores/i18n.svelte";
  import { cn } from "$lib/utils/cn";

  const store = useTunnels();
  const toast = useToast();

  const ACTIVE = ["online", "starting", "connecting"];

  let port = $state(0);
  let version = $state("");

  onMount(() => {
    statusApi
      .get()
      .then((status) => {
        port = status.port;
        version = status.version;
      })
      .catch(() => {});
  });

  const online = $derived(
    store.tunnels.filter((tunnel) => ACTIVE.includes(tunnel.state)).length,
  );

  const failing = $derived(
    store.tunnels.filter((tunnel) =>
      ["error", "timeout"].includes(tunnel.state),
    ).length,
  );

  const dashboard = $derived(port ? `${location.hostname}:${port}` : "");

  async function copyDashboard() {
    if (!dashboard) return;
    await navigator.clipboard.writeText(`http://${dashboard}`);
    toast.success(t("overview_dashboard_copied"));
  }
</script>

<section
  class="ftm-enter mb-3 flex shrink-0 flex-wrap items-stretch gap-2 text-sm"
>
  <div
    class="flex flex-1 basis-32 items-center gap-2.5 rounded-card border border-border bg-card px-3 py-2"
  >
    <span
      class={cn(
        "flex h-7 w-7 shrink-0 items-center justify-center rounded-control",
        online > 0
          ? "bg-status-running/15 text-status-running"
          : "bg-secondary text-text-muted",
      )}
    >
      <Radio size={15} />
    </span>
    <span class="min-w-0">
      <span class="block leading-tight font-semibold text-text-heading"
        >{online}</span
      >
      <span class="block truncate text-xs text-text-muted"
        >{t("overview_online")}</span
      >
    </span>
  </div>

  <div
    class="flex flex-1 basis-32 items-center gap-2.5 rounded-card border border-border bg-card px-3 py-2"
  >
    <span
      class={cn(
        "flex h-7 w-7 shrink-0 items-center justify-center rounded-control",
        failing > 0
          ? "bg-status-error/15 text-status-error"
          : "bg-secondary text-text-muted",
      )}
    >
      <Activity size={15} />
    </span>
    <span class="min-w-0">
      <span class="block leading-tight font-semibold text-text-heading">
        {failing > 0 ? failing : store.tunnels.length}
      </span>
      <span class="block truncate text-xs text-text-muted">
        {failing > 0 ? t("overview_failing") : t("overview_total")}
      </span>
    </span>
  </div>

  {#if dashboard}
    <button
      type="button"
      onclick={copyDashboard}
      title={t("overview_dashboard_hint")}
      class="flex flex-1 basis-56 cursor-pointer items-center gap-2.5 rounded-card border border-border bg-card px-3 py-2 text-left transition-colors hover:border-primary/40"
    >
      <span
        class="flex h-7 w-7 shrink-0 items-center justify-center rounded-control bg-primary/15 text-primary"
      >
        <LayoutDashboard size={15} />
      </span>
      <span class="min-w-0 flex-1">
        <span
          class="block truncate font-mono text-xs leading-tight text-text-heading"
          >{dashboard}</span
        >
        <span class="block truncate text-xs text-text-muted"
          >{t("overview_dashboard")}</span
        >
      </span>
      <Copy size={13} class="shrink-0 text-text-muted" />
    </button>
  {/if}

  {#if version}
    <div
      class="flex items-center rounded-card border border-border bg-card px-3 py-2 font-mono text-xs text-text-muted max-sm:hidden"
    >
      v{version}
    </div>
  {/if}
</section>

