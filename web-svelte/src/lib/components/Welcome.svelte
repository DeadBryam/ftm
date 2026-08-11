<script lang="ts">
  import { Plus } from "lucide-svelte";
  import Button from "./Button.svelte";
  import Mark from "./Mark.svelte";
  import { t } from "$lib/stores/i18n.svelte";

  let { onCreate }: { onCreate: () => void } = $props();

  const steps = [
    ["welcome_step_provider_title", "welcome_step_provider_body"],
    ["welcome_step_port_title", "welcome_step_port_body"],
    ["welcome_step_share_title", "welcome_step_share_body"],
  ];
</script>

<main class="flex min-h-0 flex-1 items-center justify-center py-6">
  <section
    class="ftm-enter relative w-full max-w-[34rem] overflow-hidden rounded-card border border-border bg-card px-5 py-6 text-center sm:px-8 sm:py-8"
  >
    <div class="panel-pattern"></div>

    <div class="relative flex flex-col items-center gap-3">
      <Mark class="h-12 w-12" />

      <h2
        class="m-0 font-serif text-2xl font-semibold tracking-tight text-text-heading"
      >
        {t("welcome_title")}
      </h2>

      <p class="m-0 max-w-[34ch] text-sm leading-relaxed text-text-muted">
        {t("welcome_lead")}
      </p>
    </div>

    <ol
      class="relative mx-auto mt-6 flex max-w-[26rem] list-none flex-col gap-4 p-0 text-left"
    >
      {#each steps as [title, body], index (title)}
        <li class="flex items-start gap-3">
          <span
            class="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-primary text-2xs font-semibold text-btn-text"
            aria-hidden="true"
          >
            {index + 1}
          </span>
          <div class="min-w-0">
            <p class="m-0 text-sm font-medium text-text-heading">{t(title)}</p>
            <p class="m-0 text-xs leading-relaxed text-text-muted">{t(body)}</p>
          </div>
        </li>
      {/each}
    </ol>

    <div class="relative mt-7 flex flex-col items-center gap-2">
      <Button variant="primary" size="lg" icon={Plus} onclick={onCreate}>
        {t("create_first")}
      </Button>
      <p class="m-0 text-2xs text-text-muted">{t("list_hint_new")}</p>
    </div>
  </section>
</main>
