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
    "relative h-7 w-12 shrink-0 rounded-full transition-all duration-200",
    checked ? "bg-primary" : "bg-secondary",
    disabled ? "cursor-not-allowed opacity-50" : "cursor-pointer",
  )}
>
  <span
    class={cn(
      "absolute top-0.5 w-6 h-6 bg-white rounded-full shadow transition-all duration-200",
      checked ? "left-5" : "left-0.5"
    )}
  ></span>
</button>
