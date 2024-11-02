// frontend/svelte.config.js
import preprocess from 'svelte-preprocess';
import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';

export default defineConfig({
  plugins: [
    svelte({
      preprocess: preprocess({
        typescript: true, // Включаем поддержку TypeScript
        // Вы можете настроить дополнительные опции здесь
      }),
    }),
  ],
});
