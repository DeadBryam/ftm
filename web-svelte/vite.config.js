import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vite';


const goBackend =
  (typeof process !== 'undefined' && process.env?.FTM_DEV_API) ||
  'http://127.0.0.1:40500';

export default defineConfig({
  plugins: [tailwindcss(), sveltekit()],
  server: {
    port: 3000,
    strictPort: false,
    proxy: {
      '/api': {
        target: goBackend,
        changeOrigin: true,
      },
      '/ws': {
        target: goBackend,
        changeOrigin: true,
        ws: true,
      },
    },
  },
});
