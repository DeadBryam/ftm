<script lang="ts">
  import Header from "$lib/components/Header.svelte";
  import { useI18n } from "$lib/stores/i18n.svelte";

  let { children } = $props();

  const i18n = useI18n();
</script>

<div
  class="mx-auto flex min-h-dvh w-full max-w-app flex-col px-app-x py-app-y sm:px-5 md:h-dvh md:min-h-0 md:overflow-hidden"
>
  {#if i18n.ready}
    <Header />
    <div class="flex min-h-0 flex-1 flex-col">
      {@render children()}
    </div>
  {:else}
    <div class="flex min-h-0 flex-1 flex-col gap-3" aria-busy="true" aria-live="polite">
      <div class="flex items-center gap-3 border-b border-border pb-3">
        <div class="h-9 w-9 animate-pulse rounded-control bg-border"></div>
        <div class="flex flex-1 flex-col gap-1.5">
          <div class="h-5 w-40 max-w-full animate-pulse rounded-sm bg-border"></div>
          <div class="h-3 w-56 max-w-full animate-pulse rounded-sm bg-border/70"></div>
        </div>
      </div>
      <div class="grid min-h-0 flex-1 gap-app md:grid-cols-[minmax(0,18rem)_1fr]">
        <div class="min-h-0 rounded-panel border border-border bg-card"></div>
        <div class="min-h-0 rounded-card border border-border bg-card"></div>
      </div>
    </div>
  {/if}
</div>
