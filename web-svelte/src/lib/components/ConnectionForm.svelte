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
        onSaved?.();
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
    "ftm-enter relative flex min-h-0 flex-col overflow-hidden rounded-panel border border-border bg-card",
    className,
  )}
>
  <div class="panel-pattern pointer-events-none absolute inset-0 opacity-30" aria-hidden="true"></div>

  <div class="ftm-enter ftm-enter-delay-1 relative z-10 flex shrink-0 items-center justify-between border-b border-border-light bg-card px-3.5 py-2.5">
    <h2 class="m-0 text-sm font-semibold text-text-heading">
      {isEdit ? t("edit_connection") : t("new_connection")}
    </h2>
    {#if onCancel}
      <button
        type="button"
        onclick={onCancel}
        class="cursor-pointer rounded-control border-none bg-transparent p-1 text-text-muted transition-colors hover:bg-hover hover:text-text aspect-square size-6 flex items-center justify-center"
        aria-label={t("cancel")}
      >
        ×
      </button>
    {/if}
  </div>

  <div class="ftm-enter ftm-enter-delay-2 relative z-10 flex min-h-0 flex-1 flex-col overflow-y-auto p-3.5">
    {#if providerStore.error}
      <p class="mb-2 rounded-control border border-status-error/40 bg-status-error/10 px-2.5 py-1.5 text-xs text-status-error">
        {providerStore.error}
      </p>
    {/if}

    <form onsubmit={handleSubmit}>
      <div class="mb-2.5">
        <label for="conn-name" class="mb-1 block text-xs font-medium text-text-muted">
          {isEdit ? t("connection_name_label") : t("tunnel_name")}
        </label>
        <input
          type="text"
          id="conn-name"
          bind:value={formData.name}
          placeholder={isEdit ? t("name_placeholder") : t("tunnel_name_hint")}
          required
          autocomplete="off"
          class="h-8 w-full rounded-control border border-border bg-input-bg px-2.5 py-1.5 text-sm text-text transition-all duration-150 focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-1"
        />
      </div>

      <div class="grid grid-cols-2 gap-2.5">
        <div class="mb-2.5">
          <label for="conn-port" class="mb-1 block text-xs font-medium text-text-muted">
            {isEdit ? t("port") : t("local_port")}
          </label>
          <input
            type="number"
            id="conn-port"
            bind:value={formData.localPort}
            min="1"
            max="65535"
            required
            class="h-8 w-full rounded-control border border-border bg-input-bg px-2.5 py-1.5 text-sm text-text transition-all duration-150 focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-1"
          />
        </div>
        <div class="mb-2.5">
          <label for="conn-provider" class="mb-1 block text-xs font-medium text-text-muted">
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

      <div class={cn("mt-3", isEdit && "flex gap-2")}>
        {#if isEdit}
          <Button
            variant="default"
            size="md"
            type="button"
            onclick={onCancel}
            class="flex-1"
          >
            {t("cancel")}
          </Button>
          <Button variant="primary" size="md" type="submit" class="flex-1" disabled={submitting}>
            {t("save")}
          </Button>
        {:else}
          <Button
            variant="primary"
            size="md"
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

    {#if !isEdit}
      <div class="mt-4 rounded-control border border-border-light bg-bg/40 p-3 text-xs leading-relaxed text-text-muted">
        <p class="m-0 mb-1.5 font-semibold text-text-heading">
          {t("form_tips_title")}
        </p>
        <ul class="m-0 list-none space-y-1 p-0">
          <li class="flex gap-2">
            <span class="text-primary">·</span>
            <span>{t("form_tips_port")}</span>
          </li>
          <li class="flex gap-2">
            <span class="text-primary">·</span>
            <span>{t("form_tips_provider")}</span>
          </li>
          <li class="flex gap-2">
            <span class="text-primary">·</span>
            <span>{t("form_tips_name")}</span>
          </li>
        </ul>
      </div>
    {/if}
  </div>
</section>
