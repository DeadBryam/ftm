<script lang="ts">
  import ThemeButton from "./ThemeButton.svelte";
  import SettingsToggle from "./SettingsToggle.svelte";
  import { useTheme } from "$lib/stores/theme.svelte";
  import { t } from "$lib/stores/i18n.svelte";
  import type { ThemeFamily } from "$lib/data/themes";
  import { cn } from "$lib/utils/cn";

  interface Props {
    families: ThemeFamily[];
  }

  let { families }: Props = $props();

  const theme = useTheme();

  function findCurrent(): { family: ThemeFamily; mode: "dark" | "light" } | null {
    for (const family of families) {
      if (family.dark.id === theme.current) {
        return { family, mode: "dark" };
      }
      if (family.light.id === theme.current) {
        return { family, mode: "light" };
      }
    }
    return null;
  }

  function isFamilyActive(family: ThemeFamily): boolean {
    return findCurrent()?.family.id === family.id;
  }

  function handleToggleAuto(checked: boolean) {
    if (checked) {
      const current = findCurrent();
      const familyId = current?.family.id ?? families[0]?.id ?? "foundry";
      theme.setFamily(familyId);
    } else {
      const current = findCurrent();
      const fallback = current
        ? current.mode === "dark"
          ? current.family.dark.id
          : current.family.light.id
        : families[0]?.dark.id ?? "foundry-dark";
      theme.setManual(fallback);
    }
  }

  const current = $derived(findCurrent());
  const systemModeLabel = $derived(
    theme.systemScheme === "dark" ? t("theme_dark") : t("theme_light"),
  );
  const modeLabel = $derived(theme.isAuto ? t("theme_auto") : t("theme_manual"));
</script>

<div
  class={cn(
    "mb-3 flex items-center justify-between gap-3 rounded-control border px-3 py-2",
    theme.isAuto
      ? "border-primary/30 bg-primary/5"
      : "border-border-light bg-bg/40",
  )}
>
  <div class="flex min-w-0 items-center gap-2">
    <span class="text-sm font-medium text-text-heading">
      {t("theme_match_system")}
    </span>
    <span class="truncate text-xs text-text-muted">
      {theme.isAuto
        ? `· ${systemModeLabel}`
        : `· ${t("theme_manual")}`}
    </span>
  </div>
  <SettingsToggle checked={theme.isAuto} onchange={handleToggleAuto} />
</div>

{#if theme.isAuto}
  <ul class="flex flex-col gap-1">
    {#each families as family (family.id)}
      {@const active = isFamilyActive(family)}
      <li>
        <button
          type="button"
          onclick={() => theme.setFamily(family.id)}
          class={cn(
            "group flex w-full cursor-pointer items-center gap-3 rounded-control px-2 py-1.5 text-left transition-colors",
            active ? "bg-hover" : "hover:bg-hover",
          )}
        >
          <span
            class={cn(
              "flex-1 truncate text-sm",
              active ? "font-semibold text-text-heading" : "font-medium text-text-heading",
            )}
          >
            {family.name}
          </span>
          <span
            class="relative h-8 w-16 shrink-0 overflow-hidden rounded-full transition-transform duration-200 group-hover:scale-105"
            style={active
              ? `box-shadow: 0 0 0 3px var(--color-bg), 0 0 0 5px var(--color-primary);`
              : `box-shadow: 0 2px 8px ${family.dark.color}40;`}
          >
            <span
              class="absolute inset-y-0 left-0 w-1/2"
              style="background: {family.dark.color}"
            ></span>
            <span
              class="absolute inset-y-0 right-0 w-1/2"
              style="background: {family.light.color}"
            ></span>
          </span>
        </button>
      </li>
    {/each}
  </ul>
{:else}
  <div class="flex items-center justify-end gap-6 pb-1 text-[10px] font-semibold uppercase tracking-wider text-text-muted">
    <span>{t("theme_dark")}</span>
    <span>{t("theme_light")}</span>
  </div>
  <ul class="flex flex-col gap-1">
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
          onclick={() => theme.setManual(family.dark.id)}
        />
        <ThemeButton
          id={family.light.id}
          color={family.light.color}
          selected={lightSelected}
          label={`${family.name} ${t("theme_light")}`}
          onclick={() => theme.setManual(family.light.id)}
        />
      </li>
    {/each}
  </ul>
{/if}

{#if current}
  <div class="mt-4 flex items-center gap-3 border-t border-border-light pt-3">
    <div
      class="h-9 w-9 shrink-0 rounded-full shadow-md"
      style="background: {current.mode === 'dark' ? current.family.dark.color : current.family.light.color};"
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