<script lang="ts">
  import { page } from "$app/stores";
  import MiniPanel from "$lib/components/MiniPanel.svelte";
  import { useI18n } from "$lib/stores/i18n.svelte";

  const i18n = useI18n();

  const tunnelId = $derived($page.url.searchParams.get("id") ?? "");
</script>

{#if i18n.ready}
  <MiniPanel {tunnelId} />
{/if}
