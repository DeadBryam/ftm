const THEME_IDS = [
  'foundry-dark',
  'foundry-light',
  'nord',
  'nord-light',
  'catppuccin-mocha',
  'catppuccin-latte',
  'gruvbox',
  'gruvbox-light',
  'tokyo-night',
  'tokyo-night-day',
  'dracula',
  'alucard',
  'kanagawa',
  'kanagawa-lotus',
  'one-dark',
  'one-light',
  'ayu',
  'ayu-light',
  'everforest',
  'everforest-light',
  'solarized',
  'solarized-light',
  'rose-pine',
  'rose-pine-dawn',
] as const;

export type Theme = (typeof THEME_IDS)[number];

const FAMILY_THEMES: Record<string, { dark: Theme; light: Theme }> = {
  foundry: { dark: 'foundry-dark', light: 'foundry-light' },
  nord: { dark: 'nord', light: 'nord-light' },
  catppuccin: { dark: 'catppuccin-mocha', light: 'catppuccin-latte' },
  'rose-pine': { dark: 'rose-pine', light: 'rose-pine-dawn' },
  gruvbox: { dark: 'gruvbox', light: 'gruvbox-light' },
  'tokyo-night': { dark: 'tokyo-night', light: 'tokyo-night-day' },
  dracula: { dark: 'dracula', light: 'alucard' },
  kanagawa: { dark: 'kanagawa', light: 'kanagawa-lotus' },
  one: { dark: 'one-dark', light: 'one-light' },
  ayu: { dark: 'ayu', light: 'ayu-light' },
  everforest: { dark: 'everforest', light: 'everforest-light' },
  solarized: { dark: 'solarized', light: 'solarized-light' },
};

const STORAGE_KEY = 'ftm-theme';
const DEFAULT_THEME: Theme = 'foundry-dark';
const DEFAULT_FAMILY = 'foundry';

type Preference =
  | { kind: 'auto'; family: string }
  | { kind: 'manual'; theme: Theme };

let preference = $state<Preference>({ kind: 'auto', family: DEFAULT_FAMILY });
let systemScheme = $state<'dark' | 'light'>('dark');
let currentTheme = $state<Theme>(DEFAULT_THEME);
let initialized = false;

function isTheme(value: string): value is Theme {
  return (THEME_IDS as readonly string[]).includes(value);
}

function isFamily(value: string): boolean {
  return value in FAMILY_THEMES;
}

function getSystemScheme(): 'dark' | 'light' {
  if (typeof window === 'undefined') return 'dark';
  return window.matchMedia('(prefers-color-scheme: light)').matches
    ? 'light'
    : 'dark';
}

function resolveTheme(pref: Preference, scheme: 'dark' | 'light'): Theme {
  if (pref.kind === 'manual') return pref.theme;
  const family = FAMILY_THEMES[pref.family];
  if (!family) return DEFAULT_THEME;
  return scheme === 'light' ? family.light : family.dark;
}

function loadPreference(): Preference {
  if (typeof window === 'undefined') {
    return { kind: 'auto', family: DEFAULT_FAMILY };
  }

  const raw = localStorage.getItem(STORAGE_KEY);
  if (!raw) return { kind: 'auto', family: DEFAULT_FAMILY };

  if (raw.startsWith('{')) {
    try {
      const parsed = JSON.parse(raw) as { mode?: string; family?: string; theme?: string };
      if (parsed.mode === 'auto' && parsed.family && isFamily(parsed.family)) {
        return { kind: 'auto', family: parsed.family };
      }
      if (parsed.mode === 'manual' && parsed.theme && isTheme(parsed.theme)) {
        return { kind: 'manual', theme: parsed.theme };
      }
    } catch {
      
    }
  }

  if (isTheme(raw)) {
    return { kind: 'manual', theme: raw };
  }

  return { kind: 'auto', family: DEFAULT_FAMILY };
}

function savePreference(pref: Preference): void {
  if (typeof window === 'undefined') return;

  const payload =
    pref.kind === 'auto'
      ? { mode: 'auto', family: pref.family }
      : { mode: 'manual', theme: pref.theme };

  localStorage.setItem(STORAGE_KEY, JSON.stringify(payload));
}

function applyTheme(theme: Theme): void {
  if (typeof document === 'undefined') return;
  document.documentElement.setAttribute('data-theme', theme);
}

let mediaQuery: MediaQueryList | null = null;

export function useTheme() {
  return {
    get current() {
      return currentTheme;
    },
    get preference() {
      return preference;
    },
    get systemScheme() {
      return systemScheme;
    },
    get isAuto() {
      return preference.kind === 'auto';
    },
    get isLight() {
      return Object.values(FAMILY_THEMES).some((f) => f.light === currentTheme);
    },
    get themes() {
      return [...THEME_IDS];
    },

    init() {
      if (initialized) return;
      initialized = true;

      systemScheme = getSystemScheme();
      preference = loadPreference();
      currentTheme = resolveTheme(preference, systemScheme);

      mediaQuery = window.matchMedia('(prefers-color-scheme: light)');
      mediaQuery.addEventListener('change', handleSystemChange);
      applyTheme(currentTheme);
    },

    setFamily(family: string) {
      if (!isFamily(family)) return;
      preference = { kind: 'auto', family };
      currentTheme = resolveTheme(preference, systemScheme);
      savePreference(preference);
      applyTheme(currentTheme);
    },

    setManual(theme: Theme) {
      if (!isTheme(theme)) return;
      preference = { kind: 'manual', theme };
      currentTheme = theme;
      savePreference(preference);
      applyTheme(currentTheme);
    },

    set(theme: string) {
      if (!isTheme(theme)) return;
      preference = { kind: 'manual', theme };
      currentTheme = theme;
      savePreference(preference);
      applyTheme(currentTheme);
    },

    toggle() {
      const idx = THEME_IDS.indexOf(currentTheme);
      const next = THEME_IDS[(idx + 1) % THEME_IDS.length];
      this.setManual(next);
    },
  };
}

function handleSystemChange(event: MediaQueryListEvent): void {
  systemScheme = event.matches ? 'light' : 'dark';
  if (preference.kind === 'auto') {
    currentTheme = resolveTheme(preference, systemScheme);
    applyTheme(currentTheme);
  }
}