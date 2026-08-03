export interface ThemeGroup {
  name: string;
  themes: ThemeGroupItem[];
}

export interface ThemeGroupItem {
  id: string;
  color: string;
}

export const themeGroups: ThemeGroup[] = [
  {
    name: 'Dark',
    themes: [
      { id: 'foundry-dark', color: '#c9a227' },
      { id: 'nord', color: '#88c0d0' },
      { id: 'catppuccin-mocha', color: '#cba6f7' },
      { id: 'gruvbox', color: '#fabd2f' },
      { id: 'rose-pine', color: '#ebbcba' },
    ],
  },
  {
    name: 'Light',
    themes: [{ id: 'foundry-light', color: '#8a6d1f' }],
  },
];
