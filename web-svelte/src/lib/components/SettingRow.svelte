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

<div
  role="presentation"
  onclick={() => {
    if (disabled) return;
    onchange?.(!active);
  }}
  class={cn(
    "flex w-full items-center justify-between gap-4 rounded-control",
    disabled ? "cursor-not-allowed opacity-60" : "cursor-pointer",
  )}
>
  <div class="flex min-w-0 items-center gap-3">
    <div
      class={cn(
        "flex h-9 w-9 shrink-0 items-center justify-center rounded-control transition-colors",
        active ? "bg-primary/20" : "bg-secondary",
      )}
    >
      {#if active && IconActive}
        <IconActive size={18} class="text-primary" />
      {:else if Icon}
        <Icon size={18} class="text-text-muted" />
      {/if}
    </div>
    <span class="text-sm font-medium">{label}</span>
  </div>
  <div
    role="presentation"
    onclick={(e) => e.stopPropagation()}
  >
    <SettingsToggle
      {disabled}
      checked={active}
      onchange={onchange ? (v) => onchange(v) : undefined}
    />
  </div>
</div>
