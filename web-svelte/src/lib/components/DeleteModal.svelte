<script lang="ts">
  import { X, Trash2 } from "lucide-svelte";
  import { cn } from "$lib/utils/cn";
  import { t } from "$lib/stores/i18n.svelte";
  import Button from "./Button.svelte";

  let {
    show,
    name,
    onConfirm,
    onCancel,
  }: {
    show: boolean;
    name: string;
    onConfirm: () => void;
    onCancel: () => void;
  } = $props();

  let dialogEl: HTMLDivElement | undefined = $state();

  $effect(() => {
    if (!show || !dialogEl) return;
    const previous = document.activeElement as HTMLElement | null;
    dialogEl.focus();
    return () => previous?.focus?.();
  });

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === "Escape") onCancel();
  }
</script>

<svelte:window onkeydown={show ? handleKeydown : undefined} />

{#if show}
  <div
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm ftm-enter"
    role="presentation"
    onclick={onCancel}
  >
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_interactive_supports_focus -->
    <div
      bind:this={dialogEl}
      role="dialog"
      aria-modal="true"
      aria-labelledby="modal-title"
      tabindex="-1"
      class={cn(
        "w-[90%] max-w-md rounded-panel bg-card p-4 shadow-xl outline-none ftm-enter ftm-enter-delay-1",
      )}
      onclick={(e) => e.stopPropagation()}
    >
      <div class="mb-4 flex items-start justify-between gap-3">
        <h2
          id="modal-title"
          class="m-0 flex items-center gap-2 text-lg font-semibold text-text-heading"
        >
          <Trash2 size={20} class="text-status-error" />
          {t("confirm_delete_title")}
        </h2>
        <button
          type="button"
          class="cursor-pointer rounded-lg border-none bg-transparent p-1 text-text-muted hover:bg-hover"
          onclick={onCancel}
          aria-label={t("close")}
        >
          <X size={18} />
        </button>
      </div>
      <p class="mb-6 text-sm text-text-muted">
        {t("confirm_delete_body", { 0: name })}
      </p>
      <div class="flex justify-end gap-3">
        <Button variant="default" onclick={onCancel}>{t("cancel")}</Button>
        <Button variant="error" onclick={onConfirm}>{t("delete")}</Button>
      </div>
    </div>
  </div>
{/if}
