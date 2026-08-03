<script lang="ts">
  interface Props {
    title: string;
    description?: string;
    children: import("svelte").Snippet;
  }

  let { title, description, children }: Props = $props();
</script>

<header class="mb-2.5 flex items-baseline gap-3 px-1">
  <h2
    class="font-serif text-sm font-semibold tracking-tight text-text-heading"
  >
    {title}
  </h2>
  {#if description}
    <span class="truncate text-xs text-text-muted">{description}</span>
  {/if}
</header>

<div>{@render children()}</div>