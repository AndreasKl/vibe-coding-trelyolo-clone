<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import Navbar from '$lib/components/Navbar.svelte';
	import type { Snippet } from 'svelte';

	let { children }: { children: Snippet } = $props();

	const data = $derived(page.data);

	$effect(() => {
		if (data.needsAuth) {
			goto(resolve('/login'));
		}
	});
</script>

<div class="app">
	<Navbar user={data.user} />
	<main>
		{@render children()}
	</main>
</div>

<style>
	:global(body) {
		margin: 0;
		font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial,
			sans-serif;
		background: #fafbfc;
		color: #333;
	}
	:global(*) {
		box-sizing: border-box;
	}
	.app {
		min-height: 100vh;
		display: flex;
		flex-direction: column;
	}
	main {
		flex: 1;
		display: flex;
		flex-direction: column;
	}
</style>
