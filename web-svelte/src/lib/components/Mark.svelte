<script lang="ts">
  let { class: className = "" }: { class?: string } = $props();

  const gradientId = $props.id();
</script>

<svg
  viewBox="0 0 40 40"
  class={className}
  aria-hidden="true"
  xmlns="http://www.w3.org/2000/svg"
>
  <defs>
    <linearGradient id={gradientId} x1="0" y1="0" x2="1" y2="1">
      <stop offset="0" stop-color="#b5862a" />
      <stop offset="1" stop-color="#8a6d1f" />
    </linearGradient>
  </defs>
  <rect width="40" height="40" rx="9" fill="url(#{gradientId})" />
  <g
    transform="translate(8 8)"
    fill="none"
    stroke="#fff8e1"
    stroke-width="1.6"
    stroke-linejoin="round"
    stroke-linecap="round"
  >
    <path d="M12 2 20 6.9v10.2L12 22l-8-4.9V6.9z" />
    <path d="M12 7.2 17.6 16.6H6.4z" />
  </g>
</svg>
