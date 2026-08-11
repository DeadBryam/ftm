<script lang="ts">
  import { fly } from "svelte/transition";
  import { useNotifications } from "$lib/stores/notification.svelte.js";
  import { cn } from "$lib/utils/cn";
  import { revealDuration } from "$lib/utils/motion";
  import { t } from "$lib/stores/i18n.svelte";
  import Button from "./Button.svelte";

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
      "fixed bottom-3 right-3 z-50 max-w-[300px] rounded-panel border border-border bg-card p-3.5 shadow-lg",
    )}
    in:fly={{ y: 12, duration: revealDuration(220) }}
    out:fly={{ y: 8, duration: revealDuration(140) }}
  >
    <div class="flex flex-col gap-2">
      <h3 class="m-0 text-base font-semibold text-text-heading">
        {t("enable_notifications_web")}
      </h3>
      <p class="m-0 text-sm leading-relaxed text-text-muted">
        {t("notifications_prompt")}
      </p>
      <div class="mt-1 flex gap-2">
        <Button variant="primary" size="md" onclick={request} class="flex-1">
          {t("enable_notifications")}
        </Button>
        <Button variant="default" size="md" onclick={later} class="flex-1">
          {t("not_now")}
        </Button>
      </div>
    </div>
  </div>
{/if}
