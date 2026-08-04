<script lang="ts">
  import SettingsToggle from "./SettingsToggle.svelte";
  import { useTheme } from "$lib/stores/theme.svelte";
  import { t } from "$lib/stores/i18n.svelte";
  import type { ThemeFamily } from "$lib/data/themes";
  import { cn } from "$lib/utils/cn";
  import { Monitor, Moon, Sun } from "lucide-svelte";

  interface Props {
    families: ThemeFamily[];
  }

  let { families }: Props = $props();

  const theme = useTheme();

  function familyOfTheme(themeId: string): ThemeFamily | null {
    return families.find(
      (f) => f.dark.id === themeId || f.light.id === themeId,
    ) ?? null;
  }

  const activeFamily = $derived(familyOfTheme(theme.current));
  const activeMode = $derived<'dark' | 'light' | null>(
    activeFamily
      ? activeFamily.dark.id === theme.current
        ? 'dark'
        : 'light'
      : null,
  );

  function handleToggleAuto(checked: boolean) {
    if (checked) {
      const familyId = activeFamily?.id ?? families[0]?.id ?? 'foundry';
      theme.setFamily(familyId);
    } else {
      const fallback = activeFamily
        ? activeMode === 'dark'
          ? activeFamily.dark.id
          : activeFamily.light.id
        : (families[0]?.dark.id ?? 'foundry-dark');
      theme.setManual(fallback);
    }
  }

  function pickFamily(family: ThemeFamily) {
    if (theme.isAuto) {
      theme.setFamily(family.id);
    } else {
      const nextId = activeMode === 'dark' ? family.dark.id : family.light.id;
      theme.setManual(nextId);
    }
  }
</script>

<div
  class={cn(
    'mb-3 flex items-center justify-between gap-3 rounded-control border px-3 py-2',
    theme.isAuto
      ? 'border-primary/30 bg-primary/5'
      : 'border-border-light bg-bg/40',
  )}
>
  <div class="flex min-w-0 items-center gap-2">
    <Monitor
      size={15}
      class={cn(
        'shrink-0',
        theme.isAuto ? 'text-primary' : 'text-text-muted',
      )}
    />
    <span class="truncate text-sm font-medium text-text-heading">
      {t('theme_match_system')}
    </span>
  </div>
  <SettingsToggle checked={theme.isAuto} onchange={handleToggleAuto} />
</div>

<div
  role="radiogroup"
  aria-label={t('theme_match_system')}
  class="grid grid-cols-3 gap-2"
>
  {#each families as family (family.id)}
    {@const active = activeFamily?.id === family.id}
    <div
      class={cn(
        'flex cursor-pointer flex-col items-stretch gap-2 rounded-control border p-2.5 transition-colors',
        active
          ? 'border-primary bg-primary/5'
          : 'border-border-light bg-bg/40 hover:border-primary/40 hover:bg-hover',
      )}
      role="radio"
      aria-checked={active}
      tabindex="0"
      onclick={() => pickFamily(family)}
      onkeydown={(e) =>
        (e.key === 'Enter' || e.key === ' ') &&
        (e.preventDefault(), pickFamily(family))}
    >
      <div class="relative h-9 overflow-hidden rounded-sm">
        <span class="absolute inset-y-0 left-0 w-1/2" style="background: {family.dark.color}"></span>
        <span class="absolute inset-y-0 right-0 w-1/2" style="background: {family.light.color}"></span>
        <span class="absolute inset-y-0 left-1/2 w-px bg-bg/60 mix-blend-overlay"></span>
      </div>
      <span
        class={cn(
          'truncate text-center text-xs',
          active ? 'font-semibold text-text-heading' : 'text-text-heading',
        )}
      >
        {family.name}
      </span>
    </div>
  {/each}
</div>

{#if !theme.isAuto && activeFamily}
  <div
    class="mt-3 flex flex-wrap gap-1.5"
    role="radiogroup"
    aria-label={t('theme_dark')}
  >
    <button
      type="button"
      role="radio"
      aria-checked={activeMode === 'dark'}
      onclick={() => theme.setManual(activeFamily.dark.id)}
      class={cn(
        'inline-flex cursor-pointer items-center gap-1.5 rounded-control border px-3 py-1.5 text-sm transition-all',
        activeMode === 'dark'
          ? 'border-primary bg-primary text-btn-text shadow-sm'
          : 'border-border bg-input-bg text-text hover:border-primary/50 hover:bg-hover',
      )}
    >
      <Moon size={13} />
      {t('theme_dark')}
    </button>
    <button
      type="button"
      role="radio"
      aria-checked={activeMode === 'light'}
      onclick={() => theme.setManual(activeFamily.light.id)}
      class={cn(
        'inline-flex cursor-pointer items-center gap-1.5 rounded-control border px-3 py-1.5 text-sm transition-all',
        activeMode === 'light'
          ? 'border-primary bg-primary text-btn-text shadow-sm'
          : 'border-border bg-input-bg text-text hover:border-primary/50 hover:bg-hover',
      )}
    >
      <Sun size={13} />
      {t('theme_light')}
    </button>
  </div>
{/if}