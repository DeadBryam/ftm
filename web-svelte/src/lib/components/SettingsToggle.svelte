<script lang="ts">
  import { cn } from '$lib/utils/cn';
  import { t } from '$lib/stores/i18n.svelte';
  import ToggleTrack from './ToggleTrack.svelte';

  interface Props {
    checked?: boolean;
    disabled?: boolean;
    onchange?: (checked: boolean) => void;
  }

  let { checked = false, disabled = false, onchange }: Props = $props();

  function toggle() {
    if (disabled) return;
    onchange?.(!checked);
  }
</script>

<button
  type="button"
  role="switch"
  onclick={toggle}
  {disabled}
  aria-checked={checked}
  aria-label={checked ? t('disable') : t('enable')}
  class={cn(
    "inline-flex items-center",
    disabled ? "cursor-not-allowed" : "cursor-pointer",
  )}
>
  <ToggleTrack {checked} {disabled} />
</button>
