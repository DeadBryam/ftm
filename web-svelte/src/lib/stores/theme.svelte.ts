const THEMES = [
  'foundry-dark',
  'foundry-light',
  'nord',
  'catppuccin-mocha',
  'catppuccin-latte',
  'gruvbox',
  'rose-pine',
] as const;

export type Theme = (typeof THEMES)[number];

const STORAGE_KEY = 'ftm-theme';
const DEFAULT_THEME: Theme = 'foundry-dark';

let currentTheme = $state<Theme>(DEFAULT_THEME);

function getInitialTheme(): Theme {
  if (typeof window === 'undefined') return DEFAULT_THEME;

  const saved = localStorage.getItem(STORAGE_KEY);
  if (saved && THEMES.includes(saved as Theme)) return saved as Theme;

  return DEFAULT_THEME;
}

export function useTheme() {
  return {
    get current() {
      return currentTheme;
    },
    get themes() {
      return [...THEMES];
    },

    init() {
      currentTheme = getInitialTheme();
      document.documentElement.setAttribute('data-theme', currentTheme);
    },

    set(theme: string) {
      if (!THEMES.includes(theme as Theme)) return;
      currentTheme = theme as Theme;
      document.documentElement.setAttribute('data-theme', theme);
      localStorage.setItem(STORAGE_KEY, theme);
    },

    toggle() {
      const idx = THEMES.indexOf(currentTheme);
      const next = THEMES[(idx + 1) % THEMES.length];
      this.set(next);
    },
  };
}
