<script>
  import ThemeSwitcher from './ThemeSwitcher.svelte';
</script>

<header class="flex flex-col md:flex-row md:justify-between md:items-center mb-6 pb-5 border-b flex-shrink-0 z-10 border-border">
  <div class="flex items-center gap-4">
    <img src="/favicon.png" alt="Foundry Tunnel Manager" class="w-10 h-10 sm:w-12 sm:h-12 rounded-xl object-cover" />
    <div class="text">
      <h1 class="font-serif text-2xl sm:text-3xl lg:text-4xl font-bold m-0 mb-1 tracking-tight text-text-heading">Foundry Tunnel Manager</h1>
      <p class="text-xs sm:text-sm m-0 font-medium text-text-muted">Share your world with players everywhere</p>
    </div>
  </div>
  <ThemeSwitcher />
</header>
