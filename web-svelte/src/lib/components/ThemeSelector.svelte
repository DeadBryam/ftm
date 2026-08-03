<script lang="ts">
  import ThemeButton from "./ThemeButton.svelte";
  import { useTheme } from "$lib/stores/theme.svelte";
  import { t } from "$lib/stores/i18n.svelte";
  import type { ThemeFamily, ThemeVariant } from "$lib/data/themes";

  interface Props {
    families: ThemeFamily[];
  }

  let { families }: Props = $props();

  const theme = useTheme();

  function findCurrent(): { family: ThemeFamily; variant: ThemeVariant; mode: "dark" | "light" } | null {
    for (const family of families) {
      if (family.dark.id === theme.current) {
        return { family, variant: family.dark, mode: "dark" };
      }
      if (family.light.id === theme.current) {
        return { family, variant: family.light, mode: "light" };
      }
    }
    return null;
  }

  const current = $derived(findCurrent());
  const modeLabel = $derived(
    current?.mode === "dark" ? t("theme_dark") : t("theme_light"),
  );
</script>

<div class="flex items-center justify-end gap-6 pb-1 text-[10px] font-semibold uppercase tracking-wider text-text-muted">
  <span class="w-6 text-center">◐</span>
  <span>{t("theme_dark")}</span>
  <span>{t("theme_light")}</span>
</div>

<ul class="flex flex-col gap-1.5">
  {#each families as family (family.id)}
    {@const darkSelected = theme.current === family.dark.id}
    {@const lightSelected = theme.current === family.light.id}
    <li
      class="grid grid-cols-[1fr_auto_auto] items-center gap-3 rounded-control px-2 py-1.5 transition-colors hover:bg-hover"
    >
      <span class="truncate text-sm font-medium text-text-heading">
        {family.name}
      </span>
      <ThemeButton
        id={family.dark.id}
        color={family.dark.color}
        selected={darkSelected}
        label={`${family.name} ${t("theme_dark")}`}
        onclick={() => theme.set(family.dark.id)}
      />
      <ThemeButton
        id={family.light.id}
        color={family.light.color}
        selected={lightSelected}
        label={`${family.name} ${t("theme_light")}`}
        onclick={() => theme.set(family.light.id)}
      />
    </li>
  {/each}
</ul>

{#if current}
  <div class="mt-4 flex items-center gap-3 border-t border-border-light pt-3">
    <div
      class="h-9 w-9 shrink-0 rounded-full shadow-md ring-2 ring-bg ring-offset-2 ring-offset-card"
      style="background: {current.variant.color};"
    ></div>
    <div class="min-w-0">
      <p class="truncate text-sm font-semibold text-text-heading">
        {current.family.name}
      </p>
      <p class="text-xs text-text-muted">
        {modeLabel} · {t("current_theme")}
      </p>
    </div>
  </div>
{/if}