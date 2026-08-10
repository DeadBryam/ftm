<script lang="ts">
  import { Toaster } from 'svelte-sonner';
  import { useTheme } from '$lib/stores/theme.svelte';

  const theme = useTheme();
  const mode = $derived(theme.isLight ? 'light' : 'dark');
</script>

<Toaster
  position="top-right"
  richColors
  closeButton
  theme={mode}
/>
