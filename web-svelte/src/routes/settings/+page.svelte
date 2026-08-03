<script lang="ts">
  import { onMount } from "svelte";
  import { useSettings } from "$lib/stores/settings.svelte";
  import { useNotifications } from "$lib/stores/notification.svelte";
  import { useTheme } from "$lib/stores/theme.svelte";
  import { t, useI18n, LANGUAGE_AUTO } from "$lib/stores/i18n.svelte";
  import {
    Bell,
    BellOff,
    Volume2,
    VolumeX,
    ChevronLeft,
    Check,
  } from "lucide-svelte";
  import SettingsSection from "$lib/components/SettingsSection.svelte";
  import SettingRow from "$lib/components/SettingRow.svelte";
  import ThemeSelector from "$lib/components/ThemeSelector.svelte";
  import { themeFamilies } from "$lib/data/themes";
  import { cn } from "$lib/utils/cn";

  const settingsStore = useSettings();
  const notifications = useNotifications();
  const theme = useTheme();
  const i18n = useI18n();

  let saving = $state(false);

  const languageOptions = $derived([LANGUAGE_AUTO, ...i18n.available]);

  onMount(async () => {
    theme.init();
    settingsStore.load();
    await i18n.init();
  });

  async function toggleNotifications() {
    saving = true;
    try {
      if (settingsStore.settings.notifications_enabled === "granted") {
        await settingsStore.update({ notifications_enabled: "rejected" });
        return;
      }

      await notifications.requestPermission();
      await settingsStore.load();
    } finally {
      saving = false;
    }
  }

  async function toggleSound() {
    saving = true;
    try {
      await settingsStore.update({
        notification_sound: !settingsStore.settings.notification_sound,
      });
    } finally {
      saving = false;
    }
  }

  async function changeLanguage(lang: string) {
    if (settingsStore.settings.language === lang || saving) return;
    saving = true;
    try {
      await settingsStore.update({ language: lang });
      await i18n.setLanguage(lang);
    } finally {
      saving = false;
    }
  }
</script>

<div class="flex min-h-0 w-full flex-1 flex-col overflow-y-auto">
  <div class="mb-4 flex items-center gap-2">
    <a
      href="/"
      class="cursor-pointer rounded-control p-1.5 transition-colors hover:bg-secondary"
      aria-label={t("go_back")}
    >
      <ChevronLeft size={18} />
    </a>
    <h1 class="text-lg font-semibold">{t("web_settings_title")}</h1>
    {#if saving}
      <div
        class="ml-auto h-4 w-4 animate-spin rounded-full border-2 border-primary border-t-transparent"
      ></div>
    {/if}
  </div>

  {#if !settingsStore.loaded || !i18n.ready}
    <div class="flex justify-center py-8">
      <div
        class="h-7 w-7 animate-spin rounded-full border-2 border-primary border-t-transparent"
      ></div>
    </div>
  {:else}
    <div class="grid grid-cols-1 gap-app lg:grid-cols-2 lg:items-start">
      <div class="flex flex-col gap-app">
        <SettingsSection title={t("notifications_section")}>
          {#snippet children()}
            <div class="space-y-1">
              <SettingRow
                icon={BellOff}
                iconActive={Bell}
                active={settingsStore.settings.notifications_enabled ===
                  "granted"}
                label={t("enable_notifications_web")}
                disabled={saving}
                onchange={toggleNotifications}
              />
              <SettingRow
                icon={VolumeX}
                iconActive={Volume2}
                active={settingsStore.settings.notification_sound}
                label={t("sound_effects")}
                disabled={saving}
                onchange={toggleSound}
              />
            </div>
          {/snippet}
        </SettingsSection>

        <SettingsSection title={t("language_section")}>
          {#snippet children()}
            <div
              role="radiogroup"
              aria-label={t("language_section")}
              class="flex flex-col gap-1.5"
            >
              {#each languageOptions as lang}
                {@const selected = settingsStore.settings.language === lang}
                <button
                  type="button"
                  role="radio"
                  aria-checked={selected}
                  disabled={saving}
                  onclick={() => changeLanguage(lang)}
                  class={cn(
                    "flex w-full cursor-pointer items-center gap-3 rounded-control border px-3 py-2.5 text-left transition-colors",
                    selected
                      ? "border-primary bg-primary/10 text-primary"
                      : "border-border text-text hover:border-primary/40 hover:bg-hover",
                    saving && "cursor-not-allowed opacity-60",
                  )}
                >
                  <span
                    class={cn(
                      "flex h-7 w-7 shrink-0 items-center justify-center rounded-control",
                      selected ? "bg-primary/20" : "bg-secondary",
                    )}
                  >
                    {#if selected}
                      <Check size={14} class="text-primary" />
                    {/if}
                  </span>
                  <span class="text-sm font-medium">
                    {t(`lang_${lang}`)}
                  </span>
                </button>
              {/each}
            </div>
          {/snippet}
        </SettingsSection>
      </div>

      <SettingsSection title={t("appearance_section")}>
        {#snippet children()}
          <ThemeSelector families={themeFamilies} />
        {/snippet}
      </SettingsSection>
    </div>
  {/if}
</div>
