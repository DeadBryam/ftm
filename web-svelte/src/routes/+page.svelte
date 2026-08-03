<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { useTunnels } from "$lib/stores/tunnels.svelte";
  import { useToast } from "$lib/stores/toast.svelte";
  import { useProviders } from "$lib/stores/providers.svelte";
  import { useTheme } from "$lib/stores/theme.svelte";
  import Header from "$lib/components/Header.svelte";
  import UpdateBanner from "$lib/components/UpdateBanner.svelte";
  import DeleteModal from "$lib/components/DeleteModal.svelte";
  import NotificationPermission from "$lib/components/NotificationPermission.svelte";
  import ConnectionsPanel from "$lib/components/ConnectionsPanel.svelte";
  import ConnectionForm from "$lib/components/ConnectionForm.svelte";
  import type { Tunnel } from "$lib/types";
  import { cn } from "$lib/utils/cn";
  import { t } from "$lib/stores/i18n.svelte";
  import { ChevronDown } from "lucide-svelte";

  const store = useTunnels();
  const toast = useToast();
  const providerStore = useProviders();
  const theme = useTheme();

  let deleteTunnel: Tunnel | null = $state(null);
  let editingTunnelId: string | null = $state(null);
  let formOpen = $state(true);

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

  function handleEditCancel() {
    editingTunnelId = null;
  }

  function handleEditSaved() {
    editingTunnelId = null;
  }

  function handleCreateFirst() {
    formOpen = true;
    editingTunnelId = null;
    queueMicrotask(() => {
      const el = document.getElementById("conn-name");
      el?.focus();
      el?.scrollIntoView({ behavior: "smooth", block: "center" });
    });
  }
</script>

<div class="flex min-h-0 flex-1 flex-col">
  <Header />

  <UpdateBanner />

  <main
    class={cn(
      "grid min-h-0 flex-1 gap-app",
      "max-md:grid-cols-1 max-md:content-start max-md:overflow-y-auto",
      "md:grid-cols-[minmax(0,17.5rem)_1fr] lg:grid-cols-[minmax(0,19rem)_1fr]",
    )}
  >
    <div class="flex min-h-0 flex-1 max-md:order-1 md:order-2 md:h-full md:min-h-0">
      <ConnectionsPanel
        onAction={handleAction}
        onCreateFirst={handleCreateFirst}
      />
    </div>

    <div class="max-md:order-2 md:order-1 md:flex md:h-full md:min-h-0 md:flex-col">
      <button
        type="button"
        class="mb-2 flex w-full items-center justify-between rounded-panel border border-border bg-card px-3 py-2 text-left text-sm font-medium text-text-heading md:hidden"
        onclick={() => (formOpen = !formOpen)}
        aria-expanded={formOpen}
      >
        <span>
          {editingTunnelId ? t("edit_connection") : t("new_connection")}
        </span>
        <ChevronDown
          size={18}
          class={cn("transition-transform", formOpen && "rotate-180")}
        />
      </button>

      <div class={cn("flex min-h-0 flex-1 flex-col", !formOpen && "max-md:hidden")}>
        {#if editingTunnelId}
          <ConnectionForm
            mode="edit"
            tunnelId={editingTunnelId}
            onCancel={handleEditCancel}
            onSaved={handleEditSaved}
            class="flex-1 min-h-0"
          />
        {:else}
          <ConnectionForm mode="create" class="flex-1 min-h-0" />
        {/if}
      </div>
    </div>
  </main>
</div>

<DeleteModal
  show={deleteTunnel !== null}
  name={deleteTunnel?.name ?? ""}
  onConfirm={handleConfirmDelete}
  onCancel={handleCancelDelete}
/>

<NotificationPermission />
