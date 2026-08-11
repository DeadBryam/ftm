<script lang="ts">
  import { cn } from "$lib/utils/cn";
  import { ChevronDown } from "lucide-svelte";
  import { t } from "$lib/stores/i18n.svelte";
  import type { DropdownOption } from "$lib/types";
  import type { Snippet } from "svelte";

  interface DropdownProps {
    options?: DropdownOption[];
    onSelect?: (option: DropdownOption) => void;
    align?: "left" | "right" | "top-left" | "top-right";
    ariaLabel?: string;
    label?: string;
    class?: string;
    id?: string;
    children?: Snippet;
  }

  const POSITION_MAP: Record<NonNullable<DropdownProps["align"]>, string> = {
    left: "left-auto right-0",
    right: "right-auto left-0",
    "top-left": "bottom-full mb-1.5 left-auto right-0",
    "top-right": "bottom-full mb-1.5 right-auto left-0",
  };

  let {
    options = [],
    onSelect,
    align = "left",
    ariaLabel = t("options"),
    label = t("options"),
    class: className = "",
    id = "",
    children,
  }: DropdownProps = $props();

  let isOpen = $state(false);

  const menuPosition = $derived.by(() => {
    const vert = align.startsWith("top") ? "" : "top-full mt-1.5";
    return `${POSITION_MAP[align]} ${vert}`;
  });

  function open() {
    isOpen = true;
  }

  function close() {
    isOpen = false;
  }

  function toggle(e: MouseEvent) {
    e?.stopPropagation();
    isOpen ? close() : open();
  }

  function handleOutsideClick(e: MouseEvent) {
    if (!isOpen) return;
    const target = e.target as HTMLElement;
    if (target.closest(".dropdown-trigger") || target.closest(".dropdown-menu"))
      return;
    close();
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === "Escape") close();
  }

  $effect(() => {
    if (!isOpen) return;
    document.addEventListener("click", handleOutsideClick, true);
    document.addEventListener("keydown", handleKeydown);
    return () => {
      document.removeEventListener("click", handleOutsideClick, true);
      document.removeEventListener("keydown", handleKeydown);
    };
  });
</script>

<div class={cn("dropdown-container relative flex h-fit", className)}>
  <button
    type="button"
    {id}
    onclick={toggle}
    aria-label={ariaLabel}
    aria-expanded={isOpen}
    aria-haspopup="true"
    class={cn(
      "dropdown-trigger flex h-8 min-h-8 flex-1 cursor-pointer items-center gap-1.5 rounded-control border px-2.5 py-1.5 text-xs",
      "border-border bg-card text-text hover:bg-hover",
    )}
  >
    {#if children}
      {@render children()}
    {:else}
      <span class="flex-1 text-left text-sm">{label}</span>
      <ChevronDown
        size={14}
        class={cn("transition-transform duration-150", isOpen && "rotate-180")}
      />
    {/if}
  </button>

  <div
    id={id ? `${id}-menu` : undefined}
    role="menu"
    inert={!isOpen}
    aria-hidden={!isOpen}
    aria-orientation="vertical"
    class={cn(
      "dropdown-menu absolute z-[9999] max-h-[280px] min-w-[140px] overflow-y-auto rounded-panel border border-border bg-card p-0.5",
      "origin-top transition-[opacity,transform] duration-150 ease-out",
      menuPosition,
      isOpen
        ? "pointer-events-auto translate-y-0 scale-100 opacity-100"
        : "pointer-events-none -translate-y-1 scale-95 opacity-0",
    )}
  >
    {#each options as option}
      {#if option.label === "separator"}
        <div class="mx-1.5 my-0.5 h-px bg-border"></div>
      {:else}
        <button
          type="button"
          role="menuitem"
          disabled={option.disabled}
          onclick={() => {
            close();
            onSelect?.(option);
          }}
          class={cn(
            "flex w-full cursor-pointer items-center gap-2 rounded-control border-none bg-transparent px-2.5 py-1.5 text-left text-xs text-text",
            "hover:bg-hover disabled:cursor-not-allowed disabled:opacity-50",
            option.danger && "text-red-500 hover:bg-red-500/10",
          )}
        >
          {#if option.icon}
            {@const IconComponent = option.icon as import("svelte").Component<{
              size?: number;
            }>}
            <IconComponent size={16} />
          {/if}
          <span class="flex-1">{option.label}</span>
          {#if option.hint}
            <span class="shrink-0 text-xs text-text-muted">{option.hint}</span>
          {/if}
        </button>
      {/if}
    {/each}
  </div>
</div>
