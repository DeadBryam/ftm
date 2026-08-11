<script lang="ts">
  import { X, Trash2 } from "lucide-svelte";
  import { fade, scale } from "svelte/transition";
  import { cn } from "$lib/utils/cn";
  import { revealDuration } from "$lib/utils/motion";
  import { t } from "$lib/stores/i18n.svelte";
  import Button from "./Button.svelte";
  import IconButton from "./IconButton.svelte";

  let {
    show,
    name,
    title,
    body,
    onConfirm,
    onCancel,
  }: {
    show: boolean;
    name?: string;
    title?: string;
    body?: string;
    onConfirm: () => void;
    onCancel: () => void;
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
    if (e.key === "Escape") onCancel();
  }
</script>

<svelte:window onkeydown={show ? handleKeydown : undefined} />

{#if show}
  <div
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm"
    role="presentation"
    onclick={onCancel}
    in:fade={{ duration: revealDuration(160) }}
    out:fade={{ duration: revealDuration(110) }}
  >
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_interactive_supports_focus -->
    <div
      bind:this={dialogEl}
      role="dialog"
      aria-modal="true"
      aria-labelledby="modal-title"
      tabindex="-1"
      class={cn("w-[90%] max-w-md rounded-panel bg-card p-4 shadow-xl outline-none")}
      onclick={(e) => e.stopPropagation()}
      in:scale={{ duration: revealDuration(200), start: 0.96, opacity: 0 }}
      out:scale={{ duration: revealDuration(130), start: 0.98, opacity: 0 }}
    >
      <div class="mb-4 flex items-start justify-between gap-3">
        <h2
          id="modal-title"
          class="m-0 flex items-center gap-2 text-lg font-semibold text-text-heading"
        >
          <Trash2 size={20} class="text-status-error" />
          {title ?? t("confirm_delete_title")}
        </h2>
        <IconButton
          icon={X}
          label={t("close")}
          size="sm"
          onclick={onCancel}
        />
      </div>
      <p class="mb-6 text-sm text-text-muted">
        {body ?? t("confirm_delete_body", { 0: name ?? "" })}
      </p>
      <div class="flex justify-end gap-3">
        <Button variant="default" onclick={onCancel}>{t("cancel")}</Button>
        <Button variant="error" onclick={onConfirm}>{t("delete")}</Button>
      </div>
    </div>
  </div>
{/if}
