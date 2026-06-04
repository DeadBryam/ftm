<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { page } from "$app/stores";
  import { useTunnels } from "$lib/stores/tunnels.svelte";
  import { useToast } from "$lib/stores/toast.svelte";
  import { useProviders } from "$lib/stores/providers.svelte";
  import { useTheme } from "$lib/stores/theme.svelte";
  import Header from "$lib/components/Header.svelte";
  import UpdateBanner from "$lib/components/UpdateBanner.svelte";
  import DeleteModal from "$lib/components/DeleteModal.svelte";
  import NotificationPermission from "$lib/components/NotificationPermission.svelte";
  import ConnectionsPanel from "$lib/components/ConnectionsPanel.svelte";
  import NewConnection from "$lib/components/NewConnection.svelte";
  import EditConnection from "$lib/components/EditConnection.svelte";
  import type { Tunnel } from "$lib/types";

  import { cn } from "$lib/utils/cn";
  import { translate } from "$lib/i18n";

  const store = useTunnels();
  const toast = useToast();
  const providerStore = useProviders();
  const theme = useTheme();

  let t = $derived($translate);
  let embedded = $derived($page.url.searchParams.get("embedded") === "1");

  let deleteTunnel: Tunnel | null = $state(null);
  let editingTunnelId: string | null = $state(null);

  onMount(async () => {
    theme.init();
    store.connect();
    providerStore.fetch();
  });

  onDestroy(() => {
    store.disconnect();
  });

  function handleAction(action: string, data: unknown) {
    switch (action) {
      case "edit":
        editingTunnelId = data as string;
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

  function handleEditCancel() {
    editingTunnelId = null;
  }

  function handleEditSaved() {
    editingTunnelId = null;
  }
</script>

<svelte:head>
  <link rel="preconnect" href="https://fonts.googleapis.com" />
  <link
    rel="preconnect"
    href="https://fonts.gstatic.com"
    crossorigin="anonymous"
  />
  <link
    href="https://fonts.googleapis.com/css2?family=Crimson+Pro:wght@400;500;600;700&family=Inter:wght@300;400;500;600&display=swap"
    rel="stylesheet"
  />
</svelte:head>

<div class={cn("flex-1 flex flex-col box-border", embedded ? "p-4" : "max-w-[1200px] mx-auto")}>
  {#if !embedded}
    <Header />
  {/if}

  <UpdateBanner />

  <main
    class={cn(
      "flex-1 min-h-0",
      embedded
        ? "grid grid-cols-1 gap-3"
        : "grid grid-cols-[360px_1fr] gap-5 md:grid-cols-[320px_1fr] md:gap-4 lg:grid-cols-[360px_1fr] lg:gap-5 max-md:grid-cols-1 max-md:overflow-y-auto max-md:gap-4 max-h-[calc(100dvh-11rem)]",
    )}
  >
    {#if !embedded}
      {#if editingTunnelId}
        <EditConnection
          tunnelId={editingTunnelId}
          onCancel={handleEditCancel}
          onSaved={handleEditSaved}
        />
      {:else}
        <NewConnection />
      {/if}
    {/if}

    <ConnectionsPanel onAction={handleAction} />
  </main>
</div>

<DeleteModal
  show={deleteTunnel !== null}
  name={deleteTunnel?.name ?? ""}
  onConfirm={handleConfirmDelete}
  onCancel={handleCancelDelete}
/>

<NotificationPermission />
