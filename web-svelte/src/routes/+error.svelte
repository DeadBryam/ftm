<script lang="ts">
  import { page } from "$app/state";
  import { t } from "$lib/stores/i18n.svelte";
  import { Home, ArrowLeft } from "lucide-svelte";

  const status = $derived(page.status);
  const isNotFound = $derived(status === 404);
</script>

<main class="mx-auto flex w-full max-w-app min-h-0 flex-1 items-center justify-center px-app-x py-12">
  <div class="w-full max-w-md text-center">
    <p class="font-mono text-sm text-text-muted">{status}</p>
    <h1 class="mt-3 font-serif text-3xl font-bold tracking-tight text-text-heading sm:text-4xl">
      {isNotFound ? t("error_404_title") : t("error_generic_title")}
    </h1>
    <p class="mt-3 text-sm text-text-muted">
      {isNotFound ? t("error_404_desc") : t("error_generic_desc")}
    </p>
    <div class="mt-8 flex flex-wrap items-center justify-center gap-2">
      <a
        href="/"
        class="inline-flex cursor-pointer items-center gap-1.5 rounded-control bg-primary px-4 py-2 text-sm font-medium text-btn-text shadow-sm transition-opacity hover:opacity-90"
      >
        <Home size={14} />
        {t("go_home")}
      </a>
      <button
        type="button"
        onclick={() => history.back()}
        class="inline-flex cursor-pointer items-center gap-1.5 rounded-control border border-border bg-input-bg px-4 py-2 text-sm text-text transition-colors hover:border-primary/50 hover:bg-hover"
      >
        <ArrowLeft size={14} />
        {t("go_back")}
      </button>
    </div>
  </div>
</main>