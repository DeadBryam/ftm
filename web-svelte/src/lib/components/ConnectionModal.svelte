<script lang="ts">
  import { fade, scale } from "svelte/transition";
  import ConnectionForm from "./ConnectionForm.svelte";
  import { revealDuration } from "$lib/utils/motion";

  let {
    show,
    tunnelId = null,
    onClose,
  }: {
    show: boolean;
    tunnelId?: string | null;
    onClose: () => void;
  } = $props();

  let dialogEl: HTMLDivElement | undefined = $state();

  $effect(() => {
    if (!show || !dialogEl) return;
    const previous = document.activeElement as HTMLElement | null;
    dialogEl.focus();
    document.body.style.overflow = "hidden";
    return () => {
      document.body.style.overflow = "";
      previous?.focus?.();
    };
  });

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === "Escape") onClose();
  }
</script>

<svelte:window onkeydown={show ? handleKeydown : undefined} />

{#if show}
  <div
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4 backdrop-blur-sm"
    role="presentation"
    onclick={onClose}
    in:fade={{ duration: revealDuration(160) }}
    out:fade={{ duration: revealDuration(110) }}
  >
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_interactive_supports_focus -->
    <div
      bind:this={dialogEl}
      role="dialog"
      aria-modal="true"
      tabindex="-1"
      class="w-full max-w-md outline-none"
      onclick={(e) => e.stopPropagation()}
      in:scale={{ duration: revealDuration(200), start: 0.96, opacity: 0 }}
      out:scale={{ duration: revealDuration(130), start: 0.98, opacity: 0 }}
    >
      {#if tunnelId}
        <ConnectionForm
          mode="edit"
          {tunnelId}
          onCancel={onClose}
          onSaved={onClose}
          class="shadow-xl"
        />
      {:else}
        <ConnectionForm
          mode="create"
          onCancel={onClose}
          onSaved={onClose}
          class="shadow-xl"
        />
      {/if}
    </div>
  </div>
{/if}
