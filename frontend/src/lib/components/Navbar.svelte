<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import * as api from '$lib/api';
	import type { User } from '$lib/types';

	let { user }: { user: User | null } = $props();

	async function handleLogout() {
		await api.logout();
		goto(resolve('/login'));
	}
</script>

<nav>
	<a href={resolve('/')} class="logo">Trello Clone</a>
	{#if user}
		<div class="nav-right">
			<span class="user-name">{user.name || user.email}</span>
			<button onclick={handleLogout}>Logout</button>
		</div>
	{/if}
</nav>

<style>
	nav {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 0 1.5rem;
		height: 3rem;
		background: #026aa7;
		color: white;
	}
	.logo {
		font-weight: 700;
		font-size: 1.2rem;
		color: white;
		text-decoration: none;
	}
	.nav-right {
		display: flex;
		align-items: center;
		gap: 1rem;
	}
	.user-name {
		font-size: 0.9rem;
		opacity: 0.9;
	}
	button {
		background: rgba(255, 255, 255, 0.2);
		color: white;
		border: none;
		padding: 0.4rem 0.8rem;
		border-radius: 4px;
		cursor: pointer;
		font-size: 0.85rem;
	}
	button:hover {
		background: rgba(255, 255, 255, 0.3);
	}
</style>
