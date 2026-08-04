<script lang="ts">
  import ConnectionForm from "./ConnectionForm.svelte";

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
    return () => previous?.focus?.();
  });

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === "Escape") onClose();
  }
</script>

<svelte:window onkeydown={show ? handleKeydown : undefined} />

{#if show}
  <div
    class="ftm-enter fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4 backdrop-blur-sm"
    role="presentation"
    onclick={onClose}
  >
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_interactive_supports_focus -->
    <div
      bind:this={dialogEl}
      role="dialog"
      aria-modal="true"
      tabindex="-1"
      class="ftm-enter ftm-enter-delay-1 w-full max-w-md outline-none"
      onclick={(e) => e.stopPropagation()}
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
