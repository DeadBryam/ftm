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
  <UpdateBanner />
  <OverviewBar />

  <main
    class={cn(
      "grid min-h-0 flex-1 gap-app",
      "max-md:grid-cols-1 max-md:content-start max-md:overflow-y-auto",
      "md:grid-cols-[1fr_minmax(0,17.5rem)] lg:grid-cols-[1fr_minmax(0,19rem)]",
    )}
  >
    <div class="flex min-h-0 flex-1 md:h-full md:min-h-0">
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

<NotificationPermission />
