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

  const themeNames: Record<string, string> = {
    "foundry-dark": "Foundry Dark",
    "foundry-light": "Foundry Light",
    nord: "Nord",
    "catppuccin-mocha": "Catppuccin",
    gruvbox: "Gruvbox",
    "rose-pine": "Rose Pine",
  };

  function getName(id: string): string {
    return themeNames[id] || id;
  }

  function getCurrentColor(): string {
    for (const group of groups) {
      const found = group.themes.find((item) => item.id === theme.current);
      if (found) return found.color;
    }
    return "#c9a227";
  }
</script>

{#each groups as group}
  <div class="mb-3 last:mb-0">
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

<div class="mt-4 flex items-center gap-3 border-t border-border pt-3">
  <div
    class="h-8 w-8 shrink-0 rounded-full shadow-md"
    style="background: {getCurrentColor()};"
  ></div>
  <div>
    <p class="text-sm font-medium">{getName(theme.current)}</p>
    <p class="text-xs text-text-muted">{t("current_theme")}</p>
  </div>
</div>
