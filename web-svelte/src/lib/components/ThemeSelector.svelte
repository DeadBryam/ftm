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
  class="border-t border-border-light"
  role="button"
  tabindex="0"
  onclick={() => handleToggleAuto(!theme.isAuto)}
  onkeydown={(e) =>
    (e.key === 'Enter' || e.key === ' ') &&
    (e.preventDefault(), handleToggleAuto(!theme.isAuto))}
>
  <div class="grid grid-cols-[auto_1fr_auto] items-center gap-4 px-5 py-3.5">
    <div
      class={cn(
        'flex h-9 w-9 shrink-0 items-center justify-center rounded-control transition-colors',
        theme.isAuto ? 'bg-primary/15 text-primary' : 'bg-secondary text-text-muted',
      )}
    >
      <Monitor size={17} />
    </div>
    <div class="min-w-0">
      <p class="m-0 text-sm font-medium text-text-heading">
        {t('theme_match_system')}
      </p>
      <p class="m-0 truncate text-xs text-text-muted">
        {theme.isAuto ? t('theme_auto') : t('theme_manual')}
      </p>
    </div>
    <div role="presentation" onclick={(e) => e.stopPropagation()}>
      <SettingsToggle checked={theme.isAuto} onchange={handleToggleAuto} />
    </div>
  </div>
</div>

<div
  role="radiogroup"
  aria-label={t('theme_match_system')}
  class="grid grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-2 px-5 py-4"
>
  {#each families as family (family.id)}
    {@const active = activeFamily?.id === family.id}
    <div
      class={cn(
        'flex cursor-pointer flex-col items-stretch gap-2 rounded-control p-2.5 transition-colors',
        active
          ? 'bg-primary/10 ring-1 ring-primary'
          : 'bg-secondary/50 hover:bg-hover ring-1 ring-transparent',
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
    class="border-t border-border-light px-5 py-3.5"
    role="radiogroup"
    aria-label={t('theme_dark')}
  >
    <div class="flex flex-wrap gap-1.5">
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
  </div>
{/if}