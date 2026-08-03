<script lang="ts">
  import { useNotifications } from "$lib/stores/notification.svelte.js";
  import { cn } from "$lib/utils/cn";
  import { t } from "$lib/stores/i18n.svelte";

  const notifications = useNotifications();

  let show = $derived(notifications.status === "pending");

  async function request() {
    await notifications.requestPermission();
  }

  function later() {
    notifications.reject();
  }
</script>

{#if show}
  <div
    class={cn(
      "ftm-enter fixed bottom-4 right-4 z-50 max-w-[320px] rounded-xl border border-border bg-card p-5 shadow-lg",
    )}
  >
    <div class="flex flex-col gap-2">
      <h3 class="m-0 text-[1.1rem] font-semibold text-text-heading">
        {t("enable_notifications_web")}
      </h3>
      <p class="m-0 text-sm leading-relaxed text-text-muted">
        {t("notifications_prompt")}
      </p>
      <div class="mt-1 flex gap-2">
        <button
          type="button"
          onclick={request}
          class={cn(
            "inline-flex flex-1 cursor-pointer items-center justify-center gap-2 rounded-lg border-none px-3 py-2 text-sm font-medium",
            "bg-primary text-heading transition-all duration-150 hover:-translate-y-px",
          )}
        >
          {t("enable_notifications")}
        </button>
        <button
          type="button"
          onclick={later}
          class={cn(
            "inline-flex flex-1 cursor-pointer items-center justify-center gap-2 rounded-lg border border-border px-3 py-2 text-sm font-medium",
            "bg-secondary-btn text-secondary-btn-text transition-all duration-150 hover:-translate-y-px",
          )}
        >
          {t("not_now")}
        </button>
      </div>
    </div>
  </div>
{/if}
