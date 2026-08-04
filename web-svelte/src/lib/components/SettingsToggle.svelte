<script lang="ts">
  import { cn } from '$lib/utils/cn';
  import { t } from '$lib/stores/i18n.svelte';

  interface Props {
    checked?: boolean;
    disabled?: boolean;
    onchange?: (checked: boolean) => void;
  }

  let { checked, disabled = false, onchange }: Props = $props();

  function toggle() {
    if (disabled) return;
    onchange?.(!checked);
  }
</script>

<button
  type="button"
  onclick={toggle}
  {disabled}
  aria-pressed={checked}
  aria-label={checked ? t('disable') : t('enable')}
  class={cn(
    "relative h-7 w-12 shrink-0 rounded-full border transition-all duration-200",
    checked ? "border-primary bg-primary" : "border-border bg-secondary",
    disabled ? "cursor-not-allowed opacity-50" : "cursor-pointer",
  )}
>
  <span
    class={cn(
      "absolute top-0.5 h-5 w-5 rounded-full shadow transition-all duration-200",
      checked ? "left-6 bg-btn-text" : "left-0.5 bg-muted",
    )}
  ></span>
</button>
