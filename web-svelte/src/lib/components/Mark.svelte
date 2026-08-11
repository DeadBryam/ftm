<script lang="ts">
  let { class: className = "" }: { class?: string } = $props();
</script>

<svg
  viewBox="0 0 40 40"
  class={className}
  aria-hidden="true"
  xmlns="http://www.w3.org/2000/svg"
>
  <rect width="40" height="40" rx="9" class="fill-primary" />
  <g transform="translate(8 8)">
    <path
      d="M3.34 7 L12 2 L20.66 7 L20.66 17 L12 22 L3.34 17 Z"
      class="fill-btn-text"
    />
    <g
      fill="none"
      class="stroke-primary"
      stroke-width="1.15"
      stroke-linejoin="round"
      stroke-linecap="round"
    >
      <path d="M12 5.82 L17.352 15.09 L6.648 15.09 Z" />
      <path
        d="M12 2 L12 5.82 M3.34 7 L12 5.82 M20.66 7 L12 5.82 M3.34 7 L6.648 15.09 M20.66 7 L17.352 15.09 M3.34 17 L6.648 15.09 M20.66 17 L17.352 15.09 M12 22 L6.648 15.09 M12 22 L17.352 15.09"
      />
    </g>
  </g>
</svg>
