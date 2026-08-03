<script lang="ts">
  import ThemeButton from "./ThemeButton.svelte";
  import { useTheme } from "$lib/stores/theme.svelte";
  import { t } from "$lib/stores/i18n.svelte";

  interface ThemeGroup {
    name: string;
    themes: { id: string; color: string }[];
  }

  interface Props {
    groups: ThemeGroup[];
  }

  let { groups }: Props = $props();

  const theme = useTheme();
  let showAll = $state(false);

  const FEATURED = new Set([
    "dracula",
    "nord",
    "tokyo-night",
    "catppuccin-mocha",
    "gruvbox",
    "rose-pine",
    "nord-light",
    "catppuccin-latte",
    "rose-pine-dawn",
    "gruvbox-light",
  ]);

  const themeNames: Record<string, string> = {
    dracula: "Dracula",
    nord: "Nord",
    "nord-light": "Nord Light",
    "tokyo-night": "Tokyo Night",
    "tokyo-night-storm": "Tokyo Storm",
    "tokyo-night-light": "Tokyo Light",
    "catppuccin-mocha": "Catppuccin",
    "catppuccin-latte": "Catppuccin Latte",
    "one-dark": "One Dark",
    gruvbox: "Gruvbox",
    "gruvbox-light": "Gruvbox Light",
    "solarized-dark": "Solarized",
    "solarized-light": "Solarized Light",
    "rose-pine": "Rose Pine",
    "rose-pine-dawn": "Rose Pine Dawn",
    red: "Red",
    "red-light": "Red Light",
    blue: "Blue",
    "blue-light": "Blue Light",
    purple: "Purple",
    "purple-light": "Purple Light",
    rust: "Rust",
    "rust-dark": "Rust Dark",
  };

  function getName(id: string): string {
    return themeNames[id] || id;
  }

  function getCurrentColor(): string {
    for (const group of groups) {
      const found = group.themes.find((item) => item.id === theme.current);
      if (found) return found.color;
    }
    return "#bd93f9";
  }

  const visibleGroups = $derived(
    groups.map((group) => ({
      ...group,
      themes: showAll
        ? group.themes
        : group.themes.filter(
            (item) => FEATURED.has(item.id) || item.id === theme.current,
          ),
    })).filter((group) => group.themes.length > 0),
  );

  const hiddenCount = $derived(
    groups.reduce(
      (n, g) => n + g.themes.filter((item) => !FEATURED.has(item.id)).length,
      0,
    ),
  );
</script>

{#each visibleGroups as group}
  <div class="mb-4 last:mb-0">
    <h3 class="mb-2 text-xs font-medium text-text-muted">{group.name}</h3>
    <div class="flex flex-wrap gap-2">
      {#each group.themes as item}
        <ThemeButton
          id={item.id}
          color={item.color}
          selected={theme.current === item.id}
          label={getName(item.id)}
          onclick={() => theme.set(item.id)}
        />
      {/each}
    </div>
  </div>
{/each}

{#if hiddenCount > 0}
  <button
    type="button"
    class="mb-2 cursor-pointer border-none bg-transparent p-0 text-xs font-medium text-primary hover:underline"
    onclick={() => (showAll = !showAll)}
  >
    {showAll ? t("fewer_themes") : t("more_themes")}
  </button>
{/if}

<div class="mt-4 flex items-center gap-3 border-t border-border pt-4">
  <div
    class="h-8 w-8 shrink-0 rounded-full shadow-md"
    style="background: {getCurrentColor()};"
  ></div>
  <div>
    <p class="text-sm font-medium">{getName(theme.current)}</p>
    <p class="text-xs text-text-muted">{t("current_theme")}</p>
  </div>
</div>
