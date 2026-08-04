import { api } from '$lib/api/client';

export interface Translations {
  [key: string]: string;
}

interface I18nResponse {
  translations?: Translations;
  current?: string;
  available?: string[];
}

/** Follow the operating system language. Mirrors config.LanguageAuto in Go. */
export const LANGUAGE_AUTO = 'auto';

let translations = $state<Translations>({});
let language = $state('en');
let available = $state<string[]>(['en', 'es', 'pt']);
let ready = $state(false);

/**
 * Translate a key, substituting {name} placeholders.
 *
 * Reads module-level $state directly, so calling it in a component's markup
 * re-renders that markup when the language changes. The previous store-based
 * helper read a snapshot with get(), which never updated after a language
 * switch.
 *
 * An unknown key returns the key itself, which is what the Go side does too.
 */
export function t(key: string, params?: Record<string, string | number>): string {
  let text = translations[key] ?? key;

  if (params) {
    for (const [name, value] of Object.entries(params)) {
      text = text.replaceAll(`{${name}}`, String(value));
    }
  }

  return text;
}

async function fetchCatalogue(lang?: string): Promise<void> {
  const path = lang ? `i18n?lang=${encodeURIComponent(lang)}` : 'i18n';
  const data = await api.get(path).json<I18nResponse>();

  translations = data.translations ?? {};
  language = data.current ?? language;
  available = data.available ?? available;
}

/**
 * Load the catalogue before the app renders.
 *
 * `ready` is set even when the request fails: an unreachable API should leave
 * the UI showing untranslated keys, not an empty page forever.
 */
async function init(): Promise<void> {
  if (ready) return;

  try {
    await fetchCatalogue();
  } catch (error) {
    console.error('i18n: failed to load translations', error);
  } finally {
    ready = true;
  }
}

/**
 * Switch language. Pass LANGUAGE_AUTO to follow the system again.
 *
 * Persisting the preference is the settings store's job; this only refreshes
 * the strings, so the caller does both.
 */
async function setLanguage(lang: string): Promise<void> {
  await fetchCatalogue(lang === LANGUAGE_AUTO ? undefined : lang);
}

export function useI18n() {
  return {
    get language() {
      return language;
    },
    get available() {
      return available;
    },
    get ready() {
      return ready;
    },
    init,
    setLanguage,
  };
}
