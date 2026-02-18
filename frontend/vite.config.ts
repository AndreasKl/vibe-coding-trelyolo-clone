/// <reference types="vitest" />
import { sveltekit } from '@sveltejs/kit/vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import { defineConfig } from 'vite';
import { fileURLToPath } from 'url';
import { dirname, resolve } from 'path';

const __dirname = dirname(fileURLToPath(import.meta.url));

export default defineConfig(({ mode }) => ({
	plugins: mode === 'test' ? [svelte()] : [sveltekit()],
	server: {
		proxy: {
			'/api': 'http://localhost:8080'
		}
	},
	resolve: {
		conditions: mode === 'test' ? ['browser', 'import', 'module', 'default'] : undefined,
		alias: mode === 'test' ? {
			$lib: resolve(__dirname, 'src/lib'),
			'$app/navigation': resolve(__dirname, 'src/__mocks__/app-navigation.ts'),
			'$app/paths': resolve(__dirname, 'src/__mocks__/app-paths.ts')
		} : {}
	},
	test: {
		environment: 'jsdom',
		setupFiles: ['src/test-setup.ts'],
		include: ['src/**/*.test.ts']
	}
}));
