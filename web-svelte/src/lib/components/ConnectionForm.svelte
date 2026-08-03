<script lang="ts">
  import { Plus } from "lucide-svelte";
  import { onMount } from "svelte";
  import { useProviders, detectPort } from "$lib/stores/providers.svelte";
  import { useToast } from "$lib/stores/toast.svelte";
  import { useTunnels } from "$lib/stores/tunnels.svelte";
  import { t } from "$lib/stores/i18n.svelte";
  import Button from "./Button.svelte";
  import Dropdown from "./Dropdown.svelte";
  import type { DropdownOption } from "$lib/types";
  import { cn } from "$lib/utils/cn";

  type Props = {
    mode?: "create" | "edit";
    tunnelId?: string | null;
    onCancel?: () => void;
    onSaved?: () => void;
    class?: string;
  };

  let {
    mode = "create",
    tunnelId = null,
    onCancel,
    onSaved,
    class: className = "",
  }: Props = $props();

  const store = useTunnels();
  const toast = useToast();
  const providerStore = useProviders();

  let formData = $state({
    name: "",
    provider: "cloudflared",
    localPort: 30000,
  });
  let loadedEditId = $state("");
  let submitting = $state(false);

  const isEdit = $derived(mode === "edit" && !!tunnelId);

  const providerOptions: DropdownOption[] = $derived(
    providerStore.providers.map((p) => ({
      label: p.name,
      value: p.id,
    })),
  );

  const selectedProvider = $derived(
    providerOptions.find((p) => p.value === formData.provider),
  );

  const defaultProvider = $derived(
    providerStore.providers[0]?.id ?? "cloudflared",
  );

  onMount(async () => {
    if (!isEdit) {
      const port = await detectPort();
      formData.localPort = port;
      formData.provider = defaultProvider;
    }
  });

  $effect(() => {
    if (!isEdit || !tunnelId) return;
    if (loadedEditId === tunnelId) return;

    const tunnel = store.getById(tunnelId);
    if (!tunnel) return;

    loadedEditId = tunnelId;
    formData = {
      name: tunnel.name || "",
      provider: tunnel.provider || defaultProvider,
      localPort: tunnel.port || 30000,
    };

    if (!tunnel.port) {
      detectPort().then((port) => {
        if (loadedEditId === tunnelId) {
          formData = { ...formData, localPort: port };
        }
      });
    }
  });

  function selectProvider(option: DropdownOption) {
    if (option.value) formData.provider = option.value;
  }

  async function handleSubmit(e: Event) {
    e.preventDefault();
    if (submitting) return;
    submitting = true;

    const name = formData.name;

    try {
      if (isEdit && tunnelId) {
        await store.update(tunnelId, {
          name: formData.name,
          provider: formData.provider,
          localPort: formData.localPort,
        });
        toast.success(t("connection_updated", { name }));
        onSaved?.();
      } else {
        await store.create({ ...formData });
        const detectedPort = await detectPort({ forceRefresh: true });
        formData = {
          name: "",
          provider: defaultProvider,
          localPort: detectedPort,
        };
        toast.success(t("connection_created", { name }));
      }
    } catch (err) {
      const key = isEdit ? "connection_update_failed" : "connection_create_failed";
      toast.error(t(key, { error: (err as Error).message }));
    } finally {
      submitting = false;
    }
  }
</script>

<section
  class={cn(
    "ftm-enter rounded-[var(--radius-panel)] border border-border bg-card p-5",
    className,
  )}
>
  <div class="ftm-enter ftm-enter-delay-1 mb-5 flex items-center justify-between">
    <h2 class="m-0 flex items-center gap-2 text-base font-semibold text-text-heading">
      {isEdit ? t("edit_connection") : t("new_connection")}
    </h2>
    {#if isEdit && onCancel}
      <button
        type="button"
        onclick={onCancel}
        class="cursor-pointer rounded-xl border-none bg-transparent p-1 text-text-muted transition-colors hover:bg-hover hover:text-text"
        aria-label={t("cancel")}
      >
        ×
      </button>
    {/if}
  </div>

  <div class="ftm-enter ftm-enter-delay-2">
    {#if providerStore.error}
      <p class="mb-3 rounded-lg border border-status-error/40 bg-status-error/10 px-3 py-2 text-xs text-status-error">
        {providerStore.error}
      </p>
    {/if}

    <form onsubmit={handleSubmit}>
      <div class="mb-4">
        <label for="conn-name" class="mb-1.5 block text-xs font-medium text-text-muted">
          {isEdit ? t("connection_name_label") : t("tunnel_name")}
        </label>
        <input
          type="text"
          id="conn-name"
          bind:value={formData.name}
          placeholder={isEdit ? t("name_placeholder") : t("tunnel_name_hint")}
          required
          autocomplete="off"
          class="h-9 w-full rounded-xl border border-border bg-input-bg px-3 py-2 text-sm text-text transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-1"
        />
      </div>

      <div class="grid grid-cols-2 gap-4">
        <div class="mb-4">
          <label for="conn-port" class="mb-1.5 block text-xs font-medium text-text-muted">
            {isEdit ? t("port") : t("local_port")}
          </label>
          <input
            type="number"
            id="conn-port"
            bind:value={formData.localPort}
            min="1"
            max="65535"
            required
            class="h-9 w-full rounded-xl border border-border bg-input-bg px-3 py-2 text-sm text-text transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-1"
          />
        </div>
        <div class="mb-4">
          <label for="conn-provider" class="mb-1.5 block text-xs font-medium text-text-muted">
            {isEdit ? t("provider_label") : t("select_provider")}
          </label>
          <Dropdown
            id="conn-provider"
            class="w-full"
            options={providerOptions}
            onSelect={selectProvider}
            align="left"
            ariaLabel={t("select_provider")}
            label={selectedProvider?.label || t("select")}
          />
        </div>
      </div>

      <div class={cn("mt-5", isEdit && "flex gap-3")}>
        {#if isEdit}
          <Button
            variant="default"
            size="lg"
            type="button"
            onclick={onCancel}
            class="flex-1"
          >
            {t("cancel")}
          </Button>
          <Button variant="primary" size="lg" type="submit" class="flex-1" disabled={submitting}>
            {t("save")}
          </Button>
        {:else}
          <Button
            variant="primary"
            size="lg"
            type="submit"
            class="w-full"
            icon={Plus}
            disabled={submitting}
          >
            {t("create_connection")}
          </Button>
        {/if}
      </div>
    </form>
  </div>
</section>
