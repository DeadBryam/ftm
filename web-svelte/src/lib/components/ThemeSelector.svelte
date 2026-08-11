<script lang="ts">
  import ToggleTrack from "./ToggleTrack.svelte";
  import ThemePreview from "./ThemePreview.svelte";
  import { useTheme } from "$lib/stores/theme.svelte";
  import { t } from "$lib/stores/i18n.svelte";
  import type { ThemeFamily } from "$lib/data/themes";
  import { cn } from "$lib/utils/cn";
  import { rovingRadioKeydown } from "$lib/utils/roving";
  import { revealDuration } from "$lib/utils/motion";
  import { slide } from "svelte/transition";
  import { Check, ChevronDown, Monitor, Moon, Sun } from "lucide-svelte";

  interface Props {
    families: ThemeFamily[];
  }

  let { families }: Props = $props();

  const theme = useTheme();

  let gridOpen = $state(false);

  function familyOfTheme(themeId: string): ThemeFamily | null {
    return (
      families.find((f) => f.dark.id === themeId || f.light.id === themeId) ??
      null
    );
  }

  const activeFamily = $derived(familyOfTheme(theme.current));
  const activeMode = $derived<"dark" | "light">(
    activeFamily && activeFamily.light.id === theme.current ? "light" : "dark",
  );

  function handleToggleAuto(checked: boolean) {
    if (checked) {
      theme.setFamily(activeFamily?.id ?? families[0]?.id ?? "foundry");
      return;
    }
    const fallback = activeFamily
      ? activeFamily[activeMode].id
      : (families[0]?.dark.id ?? "foundry-dark");
    theme.setManual(fallback);
  }

  function pickFamily(family: ThemeFamily) {
    if (theme.isAuto) {
      theme.setFamily(family.id);
      return;
    }
    theme.setManual(family[activeMode].id);
  }

  const modes = [
    { id: "dark", icon: Moon, labelKey: "theme_dark" },
    { id: "light", icon: Sun, labelKey: "theme_light" },
  ] as const;
</script>

<button
  type="button"
  role="switch"
  aria-checked={theme.isAuto}
  onclick={() => handleToggleAuto(!theme.isAuto)}
  class="grid w-full cursor-pointer grid-cols-[auto_1fr_auto] items-center gap-4 border-t border-border-light px-5 py-3.5 text-left transition-colors hover:bg-hover"
>
  <span
    class={cn(
      "flex h-9 w-9 shrink-0 items-center justify-center rounded-control transition-colors",
      theme.isAuto ? "bg-primary/15 text-primary" : "bg-secondary text-text-muted",
    )}
  >
    <Monitor size={17} />
  </span>
  <span class="min-w-0">
    <span class="block text-sm font-medium text-text-heading">
      {t("theme_match_system")}
    </span>
    <span class="block truncate text-xs text-text-muted">
      {theme.isAuto ? t("theme_auto") : t("theme_manual")}
    </span>
  </span>
  <ToggleTrack checked={theme.isAuto} />
</button>

<div
  class="flex items-center justify-between gap-3 border-t border-border-light px-5 py-3"
>
  <span class="min-w-0 truncate text-sm font-medium text-text-heading">
    {activeFamily?.name ?? families[0]?.name}
  </span>
  <div class="flex shrink-0 items-center gap-3">
    {#if !theme.isAuto && activeFamily}
      <div
        role="radiogroup"
        aria-label={t("theme_manual")}
        tabindex={-1}
        onkeydown={rovingRadioKeydown}
        class="flex rounded-control border border-border bg-input-bg p-0.5"
      >
        {#each modes as mode (mode.id)}
          {@const selected = activeMode === mode.id}
          <button
            type="button"
            role="radio"
            aria-checked={selected}
            tabindex={selected ? 0 : -1}
            onclick={() => theme.setManual(activeFamily[mode.id].id)}
            class={cn(
              "inline-flex cursor-pointer items-center gap-1.5 rounded-sm px-2.5 py-1 text-xs transition-colors",
              selected
                ? "bg-primary font-medium text-btn-text"
                : "text-text-muted hover:text-text",
            )}
          >
            <mode.icon size={13} />
            {t(mode.labelKey)}
          </button>
        {/each}
      </div>
    {/if}
    <button
      type="button"
      onclick={() => (gridOpen = !gridOpen)}
      aria-expanded={gridOpen}
      class="inline-flex cursor-pointer items-center gap-1 text-xs text-text-muted transition-colors hover:text-text"
    >
      {gridOpen ? t("fewer_themes") : t("more_themes")}
      <ChevronDown
        size={14}
        class={cn("transition-transform", !gridOpen && "-rotate-90")}
      />
    </button>
  </div>
</div>

{#if gridOpen}
  <div
    role="radiogroup"
    aria-label={t("appearance_section")}
    tabindex={-1}
    onkeydown={rovingRadioKeydown}
    transition:slide={{ duration: revealDuration() }}
    class="grid grid-cols-2 gap-2 px-5 pb-5 sm:grid-cols-3 lg:grid-cols-4"
  >
    {#each families as family (family.id)}
      {@const active = activeFamily?.id === family.id}
      <button
        type="button"
        role="radio"
        aria-checked={active}
        aria-label={family.name}
        tabindex={active ? 0 : -1}
        onclick={() => pickFamily(family)}
        class={cn(
          "flex cursor-pointer flex-col gap-2 rounded-control border p-2 text-left transition-colors",
          active
            ? "border-primary bg-primary/10"
            : "border-border bg-input-bg hover:border-primary/40",
        )}
      >
        <div class="flex h-14 gap-px overflow-hidden rounded-sm">
          {#if theme.isAuto}
            <ThemePreview theme={family.dark.id} class="flex-1" />
            <ThemePreview theme={family.light.id} class="flex-1" />
          {:else}
            <ThemePreview theme={family[activeMode].id} class="flex-1" />
          {/if}
        </div>
        <span class="flex items-center gap-1.5 px-0.5">
          <span
            class={cn(
              "min-w-0 flex-1 truncate text-xs",
              active ? "font-semibold text-text-heading" : "text-text",
            )}
          >
            {family.name}
          </span>
          {#if active}
            <Check size={13} class="shrink-0 text-primary" />
          {/if}
        </span>
      </button>
    {/each}
  </div>
{/if}
