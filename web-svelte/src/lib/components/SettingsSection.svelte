<script lang="ts">
  interface Props {
    title: string;
    children: import('svelte').Snippet;
  }

  let { title, children }: Props = $props();
</script>

<section class="space-y-3 rounded-panel border border-border bg-card p-4">
  <h2 class="text-xs font-medium uppercase tracking-wider text-text-muted">{title}</h2>
  {@render children()}
</section>
