<script lang="ts">
  import SettingsToggle from "./SettingsToggle.svelte";
  import ThemePreview from "./ThemePreview.svelte";
  import { useTheme } from "$lib/stores/theme.svelte";
  import { t } from "$lib/stores/i18n.svelte";
  import type { ThemeFamily } from "$lib/data/themes";
  import { cn } from "$lib/utils/cn";
  import { rovingRadioKeydown } from "$lib/utils/roving";
  import { Check, Monitor, Moon, Sun } from "lucide-svelte";

  interface Props {
    families: ThemeFamily[];
  }

  let { families }: Props = $props();

  const theme = useTheme();

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

<div class="border-t border-border-light">
  <div class="grid grid-cols-[auto_1fr_auto] items-center gap-4 px-5 py-3.5">
    <div
      class={cn(
        "flex h-9 w-9 shrink-0 items-center justify-center rounded-control transition-colors",
        theme.isAuto
          ? "bg-primary/15 text-primary"
          : "bg-secondary text-text-muted",
      )}
    >
      <Monitor size={17} />
    </div>
    <div class="min-w-0">
      <p class="m-0 text-sm font-medium text-text-heading">
        {t("theme_match_system")}
      </p>
      <p class="m-0 truncate text-xs text-text-muted">
        {theme.isAuto ? t("theme_auto") : t("theme_manual")}
      </p>
    </div>
    <div class="flex items-center gap-3">
      {#if !theme.isAuto && activeFamily}
        <div
          role="radiogroup"
          aria-label={t("theme_manual")}
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
      <SettingsToggle checked={theme.isAuto} onchange={handleToggleAuto} />
    </div>
  </div>
</div>

<div
  role="radiogroup"
  aria-label={t("appearance_section")}
  onkeydown={rovingRadioKeydown}
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
