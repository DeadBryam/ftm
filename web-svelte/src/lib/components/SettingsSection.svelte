<script lang="ts">
  interface Props {
    title: string;
    children: import('svelte').Snippet;
  }

  let { title, children }: Props = $props();
</script>

<section class="p-6 rounded-2xl bg-card border border-border space-y-4">
  <h2 class="text-sm uppercase tracking-wider text-text-muted font-medium">{title}</h2>
  {@render children()}
</section>
