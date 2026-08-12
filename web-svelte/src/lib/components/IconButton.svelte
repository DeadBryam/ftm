<script lang="ts">
  import { cn } from "$lib/utils/cn";
  import type { ComponentType, SvelteComponentTyped } from "svelte";
  import type { IconProps } from "lucide-svelte";

  type IconButtonVariant = "ghost" | "outline" | "danger";
  type IconButtonSize = "sm" | "md";

  interface IconButtonProps {
    icon: ComponentType<SvelteComponentTyped<IconProps>>;
    label: string;
    variant?: IconButtonVariant;
    size?: IconButtonSize;
    disabled?: boolean;
    type?: "button" | "submit" | "reset";
    class?: string;
    onclick?: (e: MouseEvent) => void;
  }

  const VARIANT_CLASSES: Record<IconButtonVariant, string> = {
    ghost:
      "border-transparent bg-transparent text-text-muted hover:bg-hover hover:text-text",
    outline:
      "border-border bg-input-bg text-text hover:border-primary/50 hover:bg-hover",
    danger:
      "border-status-error/40 bg-status-error/10 text-text-error hover:bg-status-error/20",
  };

  const SIZE_CLASSES: Record<IconButtonSize, string> = {
    sm: "h-7 w-7",
    md: "h-9 w-9",
  };

  let {
    icon: Icon,
    label,
    variant = "ghost",
    size = "md",
    disabled = false,
    type = "button",
    class: className = "",
    onclick,
  }: IconButtonProps = $props();

  const iconSize = $derived(size === "sm" ? 14 : 16);
</script>

<button
  {type}
  {disabled}
  {onclick}
  aria-label={label}
  title={label}
  class={cn(
    "inline-flex shrink-0 cursor-pointer items-center justify-center rounded-control border",
    "transition-all duration-150 active:scale-[0.98]",
    "disabled:cursor-not-allowed disabled:opacity-50",
    VARIANT_CLASSES[variant],
    SIZE_CLASSES[size],
    className,
  )}
>
  <Icon size={iconSize} />
</button>
