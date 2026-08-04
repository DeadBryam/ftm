<script lang="ts">
  import { cn } from "$lib/utils/cn";

  interface Props {
    theme: string;
    class?: string;
  }

  let { theme, class: className }: Props = $props();
</script>

<div
  data-theme={theme}
  class={cn("flex flex-col justify-end gap-1 overflow-hidden p-1.5", className)}
  style="background: var(--color-bg)"
>
  <div
    class="flex flex-1 flex-col justify-center gap-1 rounded-sm border px-1.5"
    style="background: var(--color-card); border-color: var(--color-border)"
  >
    <span
      class="block h-1 w-3/5 rounded-full"
      style="background: var(--color-heading)"
    ></span>
    <span
      class="block h-1 w-2/5 rounded-full"
      style="background: var(--color-muted)"
    ></span>
  </div>
  <div class="flex items-center gap-1">
    <span
      class="block h-1.5 w-1/3 rounded-full"
      style="background: var(--color-primary)"
    ></span>
    <span
      class="block h-1.5 w-1/5 rounded-full"
      style="background: var(--color-status-running)"
    ></span>
  </div>
</div>
