<script lang="ts">
  import { onMount } from "svelte";
  import { useSettings } from "$lib/stores/settings.svelte";
  import { useNotifications } from "$lib/stores/notification.svelte";
  import { useTheme } from "$lib/stores/theme.svelte";
  import { t, useI18n, LANGUAGE_AUTO } from "$lib/stores/i18n.svelte";
  import { statusApi } from "$lib/api/endpoints/status";
  import {
    Bell,
    BellOff,
    Volume2,
    VolumeX,
    Languages,
  } from "lucide-svelte";
  import SettingsToggle from "$lib/components/SettingsToggle.svelte";
  import ThemeSelector from "$lib/components/ThemeSelector.svelte";
  import { themeFamilies } from "$lib/data/themes";
  import { cn } from "$lib/utils/cn";

  const settingsStore = useSettings();
  const notifications = useNotifications();
  const theme = useTheme();
  const i18n = useI18n();

  let saving = $state(false);
  let version = $state("");

  const languageOptions = $derived([
    { id: LANGUAGE_AUTO, label: t("lang_auto"), native: t("lang_auto_native") },
    ...i18n.available.filter((l) => l !== LANGUAGE_AUTO).map((l) => ({
      id: l,
      label: t(`lang_${l}`),
      native: l === "en" ? "English" : l === "es" ? "Español" : l,
    })),
  ]);

  onMount(async () => {
    settingsStore.load();
    await i18n.init();
    statusApi
      .get()
      .then((s) => (version = s.version))
      .catch(() => {});
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

  const notifActive = $derived(
    settingsStore.settings.notifications_enabled === "granted",
  );
  const soundActive = $derived(!!settingsStore.settings.notification_sound);
</script>

<div class="mx-auto flex w-full max-w-app min-h-0 flex-1 flex-col">
  <div class="mb-4 flex items-center justify-between gap-3 px-1">
    <h1 class="font-serif text-xl font-bold tracking-tight text-text-heading">
      {t("web_settings_title")}
    </h1>
    {#if saving}
      <div
        class="h-4 w-4 animate-spin rounded-full border-2 border-primary border-t-transparent"
        aria-hidden="true"
      ></div>
    {/if}
  </div>

  <div class="min-h-0 flex-1 overflow-y-auto pb-2">
    {#if !settingsStore.loaded || !i18n.ready}
      <div class="flex justify-center py-8">
        <div
          class="h-7 w-7 animate-spin rounded-full border-2 border-primary border-t-transparent"
        ></div>
      </div>
    {:else}
      <section
        class="overflow-hidden rounded-card border border-border bg-card"
      >
        <header class="flex items-baseline justify-between gap-3 px-5 pt-5 pb-3">
          <h2
            class="font-serif text-base font-semibold tracking-tight text-text-heading"
          >
            {t("preferences_section")}
          </h2>
          <span class="text-xs text-text-muted">{t("preferences_hint")}</span>
        </header>

        <div
          class="border-t border-border-light cursor-pointer transition-colors hover:bg-hover"
          role="button"
          tabindex="0"
          onclick={toggleNotifications}
          onkeydown={(e) =>
            (e.key === "Enter" || e.key === " ") && (e.preventDefault(), toggleNotifications())}
        >
          <div
            class="grid grid-cols-[auto_1fr_auto] items-center gap-4 px-5 py-3.5"
          >
            <div
              class={cn(
                "flex h-9 w-9 shrink-0 items-center justify-center rounded-control transition-colors",
                notifActive ? "bg-primary/15 text-primary" : "bg-secondary text-text-muted",
              )}
            >
              {#if notifActive}
                <Bell size={17} />
              {:else}
                <BellOff size={17} />
              {/if}
            </div>
            <div class="min-w-0">
              <p class="m-0 text-sm font-medium text-text-heading">
                {t("enable_notifications_web")}
              </p>
              <p class="m-0 truncate text-xs text-text-muted">
                {t("settings_notifications_desc")}
              </p>
            </div>
            <div
              role="presentation"
              onclick={(e) => e.stopPropagation()}
            >
              <SettingsToggle
                disabled={saving}
                checked={notifActive}
                onchange={toggleNotifications}
              />
            </div>
          </div>
        </div>

        <div
          class="border-t border-border-light cursor-pointer transition-colors hover:bg-hover"
          role="button"
          tabindex="0"
          onclick={toggleSound}
          onkeydown={(e) =>
            (e.key === "Enter" || e.key === " ") && (e.preventDefault(), toggleSound())}
        >
          <div
            class="grid grid-cols-[auto_1fr_auto] items-center gap-4 px-5 py-3.5"
          >
            <div
              class={cn(
                "flex h-9 w-9 shrink-0 items-center justify-center rounded-control transition-colors",
                soundActive ? "bg-primary/15 text-primary" : "bg-secondary text-text-muted",
              )}
            >
              {#if soundActive}
                <Volume2 size={17} />
              {:else}
                <VolumeX size={17} />
              {/if}
            </div>
            <div class="min-w-0">
              <p class="m-0 text-sm font-medium text-text-heading">
                {t("sound_effects")}
              </p>
              <p class="m-0 truncate text-xs text-text-muted">
                {t("settings_sound_desc")}
              </p>
            </div>
            <div
              role="presentation"
              onclick={(e) => e.stopPropagation()}
            >
              <SettingsToggle
                disabled={saving}
                checked={soundActive}
                onchange={toggleSound}
              />
            </div>
          </div>
        </div>

        <div class="border-t border-border-light">
            <div class="grid grid-cols-[auto_1fr] items-center gap-4 px-5 pt-3.5 pb-2">
              <div
                class="flex h-9 w-9 shrink-0 items-center justify-center rounded-control bg-secondary text-text-muted"
              >
                <Languages size={17} />
              </div>
              <div class="min-w-0">
                <p class="m-0 text-sm font-medium text-text-heading">
                  {t("language_section")}
                </p>
                <p class="m-0 truncate text-xs text-text-muted">
                  {t("settings_language_desc")}
                </p>
              </div>
            </div>
            <div
              role="radiogroup"
              aria-label={t("language_section")}
              class="flex flex-wrap gap-1.5 px-5 pb-4"
            >
              {#each languageOptions as opt (opt.id)}
                {@const selected = settingsStore.settings.language === opt.id}
                <button
                  type="button"
                  role="radio"
                  aria-checked={selected}
                  disabled={saving}
                  onclick={() => changeLanguage(opt.id)}
                  class={cn(
                    "cursor-pointer rounded-control border px-3 py-1.5 text-sm transition-all",
                    selected
                      ? "border-primary bg-primary text-btn-text shadow-sm"
                      : "border-border bg-input-bg text-text hover:border-primary/50 hover:bg-hover",
                    saving && "cursor-not-allowed opacity-60",
                  )}
                >
                  {opt.native}
                </button>
              {/each}
            </div>
          </div>
      </section>

      <section
        class="mt-4 overflow-hidden rounded-card border border-border bg-card"
      >
        <header class="flex items-baseline justify-between gap-3 px-5 pt-5 pb-3">
          <h2
            class="font-serif text-base font-semibold tracking-tight text-text-heading"
          >
            {t("appearance_section")}
          </h2>
          <span class="text-xs text-text-muted">{t("themes_hint")}</span>
        </header>

        <ThemeSelector families={themeFamilies} />
      </section>

      <footer
        class="mt-4 flex items-center justify-between gap-3 px-1 text-xs text-text-muted"
      >
        <div class="flex min-w-0 items-center gap-2">
          <img
            src="/favicon.png"
            alt={t("app_name")}
            class="h-5 w-5 shrink-0 rounded-control object-cover opacity-80"
          />
          <span class="truncate font-medium text-text">{t("app_name")}</span>
          <span class="hidden sm:inline">·</span>
          <span class="hidden truncate sm:inline">{t("app_tagline")}</span>
        </div>
        {#if version}
          <span class="shrink-0 font-mono opacity-70">v{version}</span>
        {/if}
      </footer>
    {/if}
  </div>
</div>