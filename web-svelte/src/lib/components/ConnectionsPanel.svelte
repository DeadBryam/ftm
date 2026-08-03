<script lang="ts">
  import { Radio } from "lucide-svelte";
  import { useTunnels } from "$lib/stores/tunnels.svelte";
  import { t } from "$lib/stores/i18n.svelte";
  import { cn } from "$lib/utils/cn";
  import TunnelCard from "./TunnelCard.svelte";
  import Button from "./Button.svelte";

  let {
    onAction,
    onCreateFirst,
  }: {
    onAction: (action: string, data: unknown) => void;
    onCreateFirst?: () => void;
  } = $props();

  const store = useTunnels();
</script>

<section
  class={cn(
    "ftm-enter flex h-full min-h-[14rem] w-full flex-col overflow-hidden rounded-card border border-border bg-card",
    "max-md:max-h-[min(55dvh,28rem)]",
  )}
>
  <div
    class="ftm-enter ftm-enter-delay-1 flex shrink-0 items-center justify-between border-b border-border-light bg-url-bg px-3 py-2"
  >
    <h2
      class="m-0 flex items-center gap-2 text-sm font-semibold text-text-heading"
    >
      {t("connections")}
    </h2>
    <span
      class="rounded-control bg-primary px-2 py-0.5 text-xs font-semibold text-btn-text"
    >
      {store.tunnels.length}
    </span>
  </div>

  <div class="ftm-enter ftm-enter-delay-2 min-h-0 flex-1 overflow-y-auto p-2.5">
    {#if store.loading}
      <div
        class="flex flex-col items-center justify-center gap-2 py-8 text-text-muted"
      >
        <div
          class="h-7 w-7 animate-spin rounded-full border-2 border-border border-t-primary"
        ></div>
        <span>{t("loading")}</span>
      </div>
    {:else if store.tunnels.length === 0}
      <div class="px-2 py-6 text-center text-text-muted">
        <Radio class="mx-auto mb-2 h-8 w-8" size={32} />
        <h3 class="mb-1 mt-0 text-sm text-text-heading">
          {t("no_tunnels")}
        </h3>
        <p class="m-0 mb-3 text-xs leading-relaxed">
          {t("tunnels_desc")}
        </p>
        {#if onCreateFirst}
          <Button variant="primary" size="md" onclick={onCreateFirst}>
            {t("create_first")}
          </Button>
        {/if}
      </div>
    {:else}
      <div class="flex flex-col gap-1.5">
        {#each store.tunnels as tunnel, index (tunnel.id)}
          <TunnelCard
            {tunnel}
            {index}
            totalItems={store.tunnels.length}
            onStart={store.start}
            onStop={store.stop}
            {onAction}
            installProgress={store.installProgress[
              tunnel.provider as keyof typeof store.installProgress
            ]}
          />
        {/each}
      </div>
    {/if}
  </div>
</section>
