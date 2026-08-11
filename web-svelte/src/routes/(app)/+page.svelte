<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { useTunnels } from "$lib/stores/tunnels.svelte";
  import { useToast } from "$lib/stores/toast.svelte";
  import { useProviders } from "$lib/stores/providers.svelte";
  import UpdateBanner from "$lib/components/UpdateBanner.svelte";
  import DeleteModal from "$lib/components/DeleteModal.svelte";
  import NotificationPermission from "$lib/components/NotificationPermission.svelte";
  import ConnectionsPanel from "$lib/components/ConnectionsPanel.svelte";
  import ConnectionModal from "$lib/components/ConnectionModal.svelte";
  import DetailPanel from "$lib/components/DetailPanel.svelte";
  import OverviewBar from "$lib/components/OverviewBar.svelte";
  import Welcome from "$lib/components/Welcome.svelte";
  import type { Tunnel } from "$lib/types";
  import { cn } from "$lib/utils/cn";
  import { t } from "$lib/stores/i18n.svelte";

  const store = useTunnels();
  const toast = useToast();
  const providerStore = useProviders();

  let deleteTunnel: Tunnel | null = $state(null);
  let editingTunnelId: string | null = $state(null);
  let formOpen = $state(false);
  let selectedId: string | null = $state(null);

  const selected = $derived(
    store.tunnels.find((tunnel) => tunnel.id === selectedId) ??
      store.tunnels[0] ??
      null,
  );

  const isEmpty = $derived(!store.loading && store.tunnels.length === 0);
  const showWelcome = $derived(isEmpty && !store.onboarded && !store.error);

  let announcement = $state("");
  let lastStates: Record<string, string> = {};

  $effect(() => {
    const changed = store.tunnels.filter(
      (tunnel) =>
        lastStates[tunnel.id] !== undefined &&
        lastStates[tunnel.id] !== tunnel.state,
    );

    lastStates = Object.fromEntries(
      store.tunnels.map((tunnel) => [tunnel.id, tunnel.state]),
    );

    if (changed.length) {
      announcement = changed
        .map((tunnel) => `${tunnel.name}: ${t(tunnel.state)}`)
        .join(". ");
    }
  });

  onMount(async () => {
    store.connect();
    providerStore.fetch();
  });

  onDestroy(() => {
    store.disconnect();
  });

  function handleAction(action: string, data: unknown) {
    switch (action) {
      case "select":
        selectedId = data as string;
        break;
      case "edit":
        editingTunnelId = data as string;
        formOpen = true;
        break;
      case "delete": {
        const tunnel = data as Tunnel;
        if (deleteTunnel?.id === tunnel.id) return;
        deleteTunnel = tunnel;
        break;
      }
    }
  }

  async function handleConfirmDelete() {
    if (!deleteTunnel) return;

    const { id, name } = deleteTunnel;

    try {
      await store.delete(id);
      toast.success(t("connection_deleted", { name }));
      deleteTunnel = null;
    } catch (err) {
      toast.error(
        t("connection_delete_failed", { error: (err as Error).message }),
      );
    }
  }

  function handleCancelDelete() {
    deleteTunnel = null;
  }

  function closeForm() {
    formOpen = false;
    editingTunnelId = null;
  }

  function openCreate() {
    editingTunnelId = null;
    formOpen = true;
  }

  function handleKeydown(e: KeyboardEvent) {
    if (formOpen || deleteTunnel) return;
    const target = e.target as HTMLElement | null;
    if (target && /^(INPUT|TEXTAREA|SELECT)$/.test(target.tagName)) return;
    if (e.key === "n") {
      e.preventDefault();
      openCreate();
    }
  }
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="flex min-h-0 flex-1 flex-col">
  <div role="status" aria-live="polite" class="sr-only">{announcement}</div>

  <UpdateBanner />
  {#if !isEmpty}
    <OverviewBar />
  {/if}

  {#if showWelcome}
    <Welcome onCreate={openCreate} />
  {:else}
    <main
      class={cn(
        "grid min-h-0 flex-1 gap-app",
        "max-md:grid-cols-1 max-md:content-start",
        "md:grid-cols-3",
      )}
    >
      <div class="flex min-h-0 flex-1 max-md:flex-none md:col-span-2 md:h-full md:min-h-0">
        <ConnectionsPanel
          onAction={handleAction}
          onCreateFirst={openCreate}
          onCreate={openCreate}
          selectedId={selected?.id ?? null}
        />
      </div>

      <div class="min-h-0 max-md:min-h-64 md:h-full">
        <DetailPanel tunnel={selected} onAction={handleAction} />
      </div>
    </main>
  {/if}
</div>

<ConnectionModal
  show={formOpen}
  tunnelId={editingTunnelId}
  onClose={closeForm}
/>

<DeleteModal
  show={deleteTunnel !== null}
  name={deleteTunnel?.name ?? ""}
  onConfirm={handleConfirmDelete}
  onCancel={handleCancelDelete}
/>

{#if !showWelcome}
  <NotificationPermission />
{/if}
