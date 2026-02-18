<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import * as api from '$lib/api';
	import type { User } from '$lib/types';

	let { user }: { user: User | null } = $props();

	async function handleLogout() {
		await api.logout();
		goto(resolve('/login'), { invalidateAll: true });
	}

	function getInitials(u: User): string {
		const name = u.name || u.email;
		return name.slice(0, 2).toUpperCase();
	}
</script>

<nav>
	<a href={resolve('/')} class="logo">
		<svg width="22" height="22" viewBox="0 0 22 22" fill="none" aria-hidden="true">
			<rect x="2" y="2" width="8" height="16" rx="2.5" fill="currentColor" />
			<rect x="12" y="2" width="8" height="10" rx="2.5" fill="currentColor" opacity="0.75" />
		</svg>
		<span>FlowBoard</span>
	</a>
	{#if user}
		<div class="nav-right">
			<div class="avatar" title={user.name || user.email}>{getInitials(user)}</div>
			<button onclick={handleLogout}>Sign out</button>
		</div>
	{/if}
</nav>

<style>
	nav {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 0 1.5rem;
		height: 3.25rem;
		background: linear-gradient(90deg, #0d1b2e 0%, #1a3050 100%);
		color: white;
		box-shadow: 0 1px 0 rgba(255, 255, 255, 0.07), 0 2px 10px rgba(0, 0, 0, 0.35);
		position: relative;
		z-index: 10;
	}
	.logo {
		display: flex;
		align-items: center;
		gap: 0.55rem;
		font-weight: 700;
		font-size: 1.1rem;
		color: white;
		text-decoration: none;
		letter-spacing: -0.2px;
	}
	.logo:hover {
		opacity: 0.9;
	}
	.nav-right {
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}
	.avatar {
		width: 2rem;
		height: 2rem;
		border-radius: 50%;
		background: linear-gradient(135deg, #4a90e2, #7b52c4);
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 0.68rem;
		font-weight: 700;
		color: white;
		letter-spacing: 0.5px;
		flex-shrink: 0;
	}
	button {
		background: rgba(255, 255, 255, 0.1);
		color: rgba(255, 255, 255, 0.85);
		border: 1px solid rgba(255, 255, 255, 0.15);
		padding: 0.35rem 0.9rem;
		border-radius: 6px;
		cursor: pointer;
		font-size: 0.82rem;
		transition: background 0.15s, color 0.15s;
	}
	button:hover {
		background: rgba(255, 255, 255, 0.2);
		color: white;
	}
</style>
