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
    "relative block h-7 w-12 shrink-0 rounded-full border transition-all duration-200",
    checked ? "border-primary bg-primary" : "border-border bg-secondary",
    disabled && "opacity-50",
  )}
>
  <span
    class={cn(
      "absolute top-0.5 h-5 w-5 rounded-full shadow transition-all duration-200",
      checked ? "left-6 bg-btn-text" : "left-0.5 bg-muted",
    )}
  ></span>
</span>
