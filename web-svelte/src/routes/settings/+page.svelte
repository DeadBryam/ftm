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
    Languages,
  } from "lucide-svelte";
  import SettingsSection from "$lib/components/SettingsSection.svelte";
  import SettingRow from "$lib/components/SettingRow.svelte";
  import ThemeSelector from "$lib/components/ThemeSelector.svelte";
  import { themeGroups } from "$lib/data/themes";
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
    <div class="columns-1 gap-app md:columns-2">
      <div class="mb-app break-inside-avoid">
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
      </div>

      <div class="mb-app break-inside-avoid">
        <SettingsSection title={t("appearance_section")}>
          {#snippet children()}
            <ThemeSelector groups={themeGroups} />
          {/snippet}
        </SettingsSection>
      </div>

      <div class="mb-app break-inside-avoid">
        <SettingsSection title={t("language_section")}>
          {#snippet children()}
            <div
              role="radiogroup"
              aria-label={t("language_section")}
              class="grid grid-cols-1 gap-2 sm:grid-cols-3"
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
                    "flex cursor-pointer items-center gap-2.5 rounded-control border px-3 py-2.5 text-left transition-colors",
                    selected
                      ? "border-primary bg-primary/10 text-primary"
                      : "border-border text-text hover:border-primary/40 hover:bg-hover",
                    saving && "cursor-not-allowed opacity-60",
                  )}
                >
                  <span
                    class={cn(
                      "flex h-8 w-8 shrink-0 items-center justify-center rounded-control",
                      selected ? "bg-primary/20" : "bg-secondary",
                    )}
                  >
                    {#if selected}
                      <Check size={16} class="text-primary" />
                    {:else}
                      <Languages size={16} class="text-text-muted" />
                    {/if}
                  </span>
                  <span class="min-w-0 flex-1 text-sm font-medium leading-tight">
                    {t(`lang_${lang}`)}
                  </span>
                </button>
              {/each}
            </div>
          {/snippet}
        </SettingsSection>
      </div>
    </div>
  {/if}
</div>
