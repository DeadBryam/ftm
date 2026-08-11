<script lang="ts">
  import { Plus, PlugZap, Radio, RefreshCw, Search } from "lucide-svelte";
  import { useTunnels } from "$lib/stores/tunnels.svelte";
  import { useToast } from "$lib/stores/toast.svelte";
  import { t } from "$lib/stores/i18n.svelte";
  import { cn } from "$lib/utils/cn";
  import TunnelCard from "./TunnelCard.svelte";
  import Button from "./Button.svelte";

  let {
    onAction,
    onCreateFirst,
    onCreate,
    selectedId = null,
  }: {
    onAction: (action: string, data: unknown) => void;
    onCreateFirst?: () => void;
    onCreate?: () => void;
    selectedId?: string | null;
  } = $props();

  const store = useTunnels();
  const toast = useToast();

  const SEARCH_THRESHOLD = 8;

  async function start(id: string) {
    try {
      await store.start(id);
    } catch (err) {
      toast.error((err as Error).message);
    }
  }

  async function stop(id: string) {
    try {
      await store.stop(id);
    } catch (err) {
      toast.error((err as Error).message);
    }
  }

  let query = $state("");

  const searchable = $derived(store.tunnels.length >= SEARCH_THRESHOLD);

  const STATE_ORDER: Record<string, number> = {
    error: 0,
    timeout: 0,
    online: 1,
    starting: 2,
    connecting: 2,
    stopping: 2,
    installing: 3,
    downloading: 3,
  };

  const visible = $derived.by(() => {
    const term = query.trim().toLowerCase();
    const matches =
      !term || !searchable
        ? store.tunnels
        : store.tunnels.filter(
            (tunnel) =>
              tunnel.name.toLowerCase().includes(term) ||
              tunnel.provider.toLowerCase().includes(term) ||
              String(tunnel.port).includes(term),
          );

    return [...matches].sort(
      (a, b) => (STATE_ORDER[a.state] ?? 9) - (STATE_ORDER[b.state] ?? 9),
    );
  });
</script>

<section
  class={cn(
    "ftm-enter flex h-full min-h-56 w-full flex-col overflow-hidden rounded-card border border-border bg-card",
    "max-md:h-auto",
  )}
>
  <div
    class="ftm-enter ftm-enter-delay-1 flex shrink-0 items-center justify-between gap-3 border-b border-border-light bg-url-bg px-3 py-2"
  >
    <h2
      class="m-0 flex min-w-0 items-center gap-2 text-sm font-semibold text-text-heading"
    >
      {t("connections")}
      <span
        class="rounded-control bg-primary px-2 py-0.5 text-2xs font-semibold text-btn-text"
      >
        {store.tunnels.length}
      </span>
    </h2>
    {#if onCreate}
      <Button
        variant="primary"
        size="sm"
        icon={Plus}
        onclick={onCreate}
        class="max-sm:aspect-square max-sm:px-0"
        ariaLabel={t("new_connection_action")}
      >
        <span class="max-sm:hidden">{t("new_connection_action")}</span>
      </Button>
    {/if}
  </div>

  {#if searchable}
    <div class="shrink-0 border-b border-border-light px-3 py-2">
      <div class="relative">
        <Search
          size={14}
          class="pointer-events-none absolute top-1/2 left-2.5 -translate-y-1/2 text-text-muted"
        />
        <input
          type="search"
          bind:value={query}
          placeholder={t("search_placeholder")}
          aria-label={t("search_placeholder")}
          class="h-8 w-full rounded-control border border-border bg-input-bg py-1.5 pr-2.5 pl-8 text-sm text-text transition-colors duration-150 hover:border-primary/50"
        />
      </div>
    </div>
  {/if}

  <div class="ftm-enter ftm-enter-delay-2 relative min-h-0 flex-1 overflow-hidden">
    <div class="panel-pattern" aria-hidden="true"></div>

    <div class="relative z-10 h-full min-h-0 overflow-y-auto p-2.5 max-md:h-auto max-md:overflow-visible">
      {#if store.loading}
        <div
          class="flex flex-col items-center justify-center gap-2 py-8 text-text-muted"
        >
          <div
            class="h-7 w-7 animate-spin rounded-full border-2 border-border border-t-primary"
          ></div>
          <span>{t("loading")}</span>
        </div>
      {:else if store.error}
        <div
          class="flex h-full min-h-40 flex-col items-center justify-center px-3 py-6 text-center text-text-muted"
        >
          <PlugZap class="mx-auto mb-2 h-8 w-8 text-status-error" size={32} />
          <h3 class="mt-0 mb-1 text-sm text-text-heading">
            {t("panel_unreachable")}
          </h3>
          <p class="m-0 mb-3 max-w-xs text-xs leading-relaxed">
            {t("panel_unreachable_desc")}
          </p>
          <p class="m-0 mb-3 max-w-xs font-mono text-xs break-words">
            {store.error}
          </p>
          <Button variant="primary" icon={RefreshCw} onclick={() => store.retry()}>
            {t("retry")}
          </Button>
        </div>
      {:else if store.tunnels.length === 0}
        <div
          class="flex h-full min-h-40 flex-col items-center justify-center px-3 py-6 text-center text-text-muted"
        >
          <Radio class="mx-auto mb-2 h-8 w-8" size={32} />
          <h3 class="mt-0 mb-1 text-sm text-text-heading">
            {t("no_tunnels")}
          </h3>
          <p class="m-0 mb-3 max-w-xs text-xs leading-relaxed">
            {t("tunnels_desc")}
          </p>
          {#if onCreateFirst}
            <Button variant="primary"  onclick={onCreateFirst}>
              {t("create_first")}
            </Button>
          {/if}
        </div>
      {:else if visible.length === 0}
        <p class="px-3 py-6 text-center text-xs text-text-muted">
          {t("search_empty")}
        </p>
      {:else}
        <div class="flex flex-col gap-1.5">
          {#each visible as tunnel, index (tunnel.id)}
            <TunnelCard
              {tunnel}
              {index}
              totalItems={visible.length}
              selected={tunnel.id === selectedId}
              onStart={start}
              onStop={stop}
              {onAction}
              installProgress={store.installProgress[
                tunnel.provider as keyof typeof store.installProgress
              ]}
            />
          {/each}
        </div>
      {/if}
    </div>
  </div>

  {#if store.tunnels.length > 0}
    <p
      class="shrink-0 border-t border-border-light px-3 py-1.5 text-center text-2xs text-text-muted max-md:hidden"
    >
      {t("list_hint_new")}
    </p>
  {/if}
</section>
