<script lang="ts">
  import { useTheme } from "$lib/stores/theme.svelte";
  import { useSound } from "$lib/stores/sound.svelte";
  import Dropdown from './Dropdown.svelte';
  import type { DropdownOption } from '$lib/types';

  const theme = useTheme();
  const sound = useSound();

  const themeLabels: Record<string, string> = {
    nord: "Nord",
    "nord-light": "Nord Light",
    "rose-pine": "Rose Pine",
    "rose-pine-dawn": "Rose Pine Dawn",
    "tokyo-night": "Tokyo Night",
    "tokyo-night-storm": "Tokyo Night Storm",
    "tokyo-night-light": "Tokyo Night Light",
    "catppuccin-mocha": "Catppuccin Mocha",
    "catppuccin-latte": "Catppuccin Latte",
    "one-dark": "One Dark",
    gruvbox: "Gruvbox",
    "gruvbox-light": "Gruvbox Light",
    "solarized-dark": "Solarized Dark",
    "solarized-light": "Solarized Light",
    dracula: "Dracula",
    red: "Red",
    'red-light': "Red Light",
    blue: "Blue",
    'blue-light': "Blue Light",
    purple: "Purple",
    'purple-light': "Purple Light",
  };

  const themeOptions = $derived(theme.themes.map(t => ({
    label: themeLabels[t] || t,
    value: t
  })));

  const selectedTheme = $derived(themeOptions.find(t => t.value === theme.current));

  function selectTheme(option: DropdownOption) {
    if (option.value) {
      theme.set(option.value);
    }
  }
</script>

<div class="flex items-center gap-2 mt-[10px]" role="group" aria-label="Theme and sound controls">
  <Dropdown 
    options={themeOptions} 
    onSelect={selectTheme}
    align="left"
    class="min-w-38"
    ariaLabel="Select theme"
    label={selectedTheme?.label || 'Theme'}
  />

  <button
    class="w-9 h-9 flex items-center justify-center rounded-md border cursor-pointer text-base transition-colors duration-150 bg-card border-border text-text"
    onclick={() => sound.toggle()}
    aria-pressed={sound.enabled}
    title={sound.enabled ? "Sound on" : "Sound off"}
  >
    {sound.enabled ? "🔊" : "🔇"}
  </button>
</div>
