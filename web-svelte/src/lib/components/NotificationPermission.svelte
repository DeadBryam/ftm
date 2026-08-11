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
      "fixed right-3 bottom-3 left-3 z-50 rounded-panel border border-border bg-card p-3 shadow-lg sm:left-auto sm:w-[19rem]",
    )}
    in:fly={{ y: 12, duration: revealDuration(220) }}
    out:fly={{ y: 8, duration: revealDuration(140) }}
  >
    <div class="flex flex-col gap-1.5">
      <h3 class="m-0 text-sm font-semibold text-text-heading">
        {t("enable_notifications_web")}
      </h3>
      <p class="m-0 text-xs text-text-muted">
        {t("notifications_prompt")}
      </p>
      <div class="mt-1 flex gap-2">
        <Button variant="primary" onclick={request} class="flex-1">
          {t("enable")}
        </Button>
        <Button variant="default" onclick={later} class="shrink-0">
          {t("not_now")}
        </Button>
      </div>
    </div>
  </div>
{/if}
