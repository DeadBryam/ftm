<script lang="ts">
  import { cn } from "$lib/utils/cn";

  let {
    checked = false,
    disabled = false,
  }: {
    checked?: boolean;
    disabled?: boolean;
  } = $props();
</script>

<span
  aria-hidden="true"
  class={cn(
    "relative block h-5 w-9 shrink-0 rounded-full border transition-colors duration-200",
    checked ? "border-primary bg-primary" : "border-border bg-secondary",
    disabled && "opacity-50",
  )}
>
  <span
    class={cn(
      "absolute top-1/2 h-3.5 w-3.5 -translate-y-1/2 rounded-full shadow-sm transition-all duration-200",
      checked ? "left-[calc(100%-1rem)] bg-btn-text" : "left-0.5 bg-muted",
    )}
  ></span>
</span>
