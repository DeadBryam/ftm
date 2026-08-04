<script lang="ts">
  import SettingsToggle from "./SettingsToggle.svelte";
  import { useTheme } from "$lib/stores/theme.svelte";
  import { t } from "$lib/stores/i18n.svelte";
  import type { ThemeFamily } from "$lib/data/themes";
  import { cn } from "$lib/utils/cn";
  import { Monitor } from "lucide-svelte";

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

  function familyOfActive(): ThemeFamily | null {
    return findCurrent()?.family ?? null;
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
  const activeFamily = $derived(familyOfActive());
</script>

<div
  class={cn(
    "mb-4 flex items-center justify-between gap-3 rounded-control border px-3 py-2",
    theme.isAuto
      ? "border-primary/30 bg-primary/5"
      : "border-border-light bg-bg/40",
  )}
>
  <div class="flex min-w-0 items-center gap-2">
    <Monitor
      size={15}
      class={cn(
        "shrink-0",
        theme.isAuto ? "text-primary" : "text-text-muted",
      )}
    />
    <span class="truncate text-sm font-medium text-text-heading">
      {t("theme_match_system")}
    </span>
  </div>
  <SettingsToggle checked={theme.isAuto} onchange={handleToggleAuto} />
</div>

{#if theme.isAuto}
  <div
    role="radiogroup"
    aria-label={t("theme_match_system")}
    class="grid grid-cols-3 gap-2 sm:grid-cols-4"
  >
    {#each families as family (family.id)}
      {@const active = activeFamily?.id === family.id}
      <button
        type="button"
        role="radio"
        aria-checked={active}
        onclick={() => theme.setFamily(family.id)}
        class={cn(
          "group flex cursor-pointer flex-col items-stretch gap-1.5 rounded-control border p-2 text-left transition-all",
          active
            ? "border-primary bg-primary/5 shadow-sm"
            : "border-border-light bg-bg/40 hover:border-primary/40 hover:bg-hover",
        )}
      >
        <span
          class="relative h-6 w-full overflow-hidden rounded-control transition-transform duration-200 group-hover:scale-[1.02]"
          style={active
            ? `box-shadow: 0 0 0 2px var(--color-bg), 0 0 0 4px var(--color-primary);`
            : `box-shadow: 0 1px 4px ${family.dark.color}50;`}
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
        <span
          class={cn(
            "truncate text-center text-xs font-medium",
            active ? "font-semibold text-text-heading" : "text-text-heading",
          )}
        >
          {family.name}
        </span>
      </button>
    {/each}
  </div>
{:else}
  <div
    role="radiogroup"
    aria-label={t("theme_dark")}
    class="grid grid-cols-3 gap-2 sm:grid-cols-4"
  >
    {#each families as family (family.id)}
      {@const darkSelected = theme.current === family.dark.id}
      {@const lightSelected = theme.current === family.light.id}
      <div
        class={cn(
          "flex flex-col items-stretch gap-1.5 rounded-control border p-2 transition-colors",
          darkSelected || lightSelected
            ? "border-primary/40 bg-primary/5"
            : "border-border-light bg-bg/40 hover:border-primary/30",
        )}
      >
        <div class="relative h-7 overflow-hidden rounded-control shadow-sm">
          <span class="absolute inset-y-0 left-0 w-1/2" style="background: {family.dark.color}"></span>
          <span class="absolute inset-y-0 right-0 w-1/2" style="background: {family.light.color}"></span>
          <span class="absolute inset-y-0 left-1/2 w-px bg-bg/60 mix-blend-overlay"></span>
          <button
            type="button"
            onclick={() => theme.setManual(family.dark.id)}
            aria-label={`${family.name} ${t("theme_dark")}`}
            aria-pressed={darkSelected}
            class={cn(
              "absolute inset-y-0 left-0 w-1/2 cursor-pointer transition-all",
              darkSelected ? "ring-2 ring-primary ring-inset" : "hover:brightness-110",
            )}
          ></button>
          <button
            type="button"
            onclick={() => theme.setManual(family.light.id)}
            aria-label={`${family.name} ${t("theme_light")}`}
            aria-pressed={lightSelected}
            class={cn(
              "absolute inset-y-0 right-0 w-1/2 cursor-pointer transition-all",
              lightSelected ? "ring-2 ring-primary ring-inset" : "hover:brightness-110",
            )}
          ></button>
        </div>
        <span class="truncate text-center text-xs font-medium text-text-heading">
          {family.name}
        </span>
      </div>
    {/each}
  </div>
{/if}