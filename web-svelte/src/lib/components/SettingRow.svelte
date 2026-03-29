<script lang="ts">
  import type { IconProps } from "lucide-svelte";
  import SettingsToggle from "./SettingsToggle.svelte";
  import type { ComponentType, SvelteComponentTyped } from "svelte";
  import { cn } from "$lib/utils/cn";

  interface Props {
    icon?: ComponentType<SvelteComponentTyped<IconProps>>;
    iconActive?: ComponentType<SvelteComponentTyped<IconProps>>;
    active: boolean;
    label: string;
    disabled?: boolean;
    onchange?: (checked: boolean) => void;
  }

  let {
    icon: Icon,
    iconActive: IconActive,
    active,
    label,
    disabled = false,
    onchange,
  }: Props = $props();
</script>

<div class="flex items-center justify-between gap-4">
  <div class="flex items-center gap-3">
    <div
      class={cn(
        "w-10 h-10 rounded-lg flex items-center justify-center transition-colors",
        active ? "bg-primary/20" : "bg-secondary",
      )}
    >
      {#if active && IconActive}
        <IconActive size={20} class="text-primary" />
      {:else}
        <Icon size={20} class="text-text-muted" />
      {/if}
    </div>
    <span class="font-medium">{label}</span>
  </div>
  <SettingsToggle
    {disabled}
    checked={active}
    onchange={onchange ? (v) => onchange(v) : undefined}
  />
</div>
