<script lang="ts">
	import { goto, invalidateAll } from '$app/navigation';
	import { resolve } from '$app/paths';
	import * as api from '$lib/api';

	let email = $state('');
	let password = $state('');
	let error = $state('');

	async function handleLogin() {
		error = '';
		try {
			await api.login(email, password);
			await invalidateAll();
			goto(resolve('/'));
		} catch (e: unknown) {
			error = e instanceof Error ? e.message : 'Login failed';
		}
	}
</script>

<div class="auth-page">
	<div class="auth-card">
		<h1>Log in</h1>
		{#if error}
			<div class="error">{error}</div>
		{/if}
		<form onsubmit={(e) => { e.preventDefault(); handleLogin(); }}>
			<input type="email" bind:value={email} placeholder="Email" required />
			<input type="password" bind:value={password} placeholder="Password" required />
			<button type="submit">Log in</button>
		</form>

		<div class="oauth-divider">or</div>
		<!-- eslint-disable svelte/no-navigation-without-resolve -- external API routes, not SvelteKit pages -->
		<a href="/api/auth/oauth/google" class="oauth-btn google">Sign in with Google</a>
		<a href="/api/auth/oauth/microsoft" class="oauth-btn microsoft">Sign in with Microsoft</a>
		<!-- eslint-enable svelte/no-navigation-without-resolve -->

		<p class="switch">Don't have an account? <a href={resolve('/signup')}>Sign up</a></p>
	</div>
</div>

<style>
	.auth-page {
		display: flex;
		align-items: center;
		justify-content: center;
		flex: 1;
		padding: 2rem;
	}
	.auth-card {
		background: white;
		padding: 2rem;
		border-radius: 8px;
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
		width: 100%;
		max-width: 380px;
	}
	h1 {
		text-align: center;
		margin: 0 0 1.5rem;
		font-size: 1.3rem;
	}
	.error {
		background: #fef2f2;
		color: #dc2626;
		padding: 0.5rem;
		border-radius: 4px;
		font-size: 0.85rem;
		margin-bottom: 1rem;
	}
	form {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}
	input {
		padding: 0.6rem;
		border: 1px solid #ccc;
		border-radius: 4px;
		font-size: 0.9rem;
	}
	button {
		background: #5aac44;
		color: white;
		border: none;
		padding: 0.6rem;
		border-radius: 4px;
		cursor: pointer;
		font-size: 0.9rem;
		font-weight: 600;
	}
	button:hover {
		background: #4e9c3b;
	}
	.oauth-divider {
		text-align: center;
		color: #999;
		margin: 1rem 0;
		font-size: 0.85rem;
	}
	.oauth-btn {
		display: block;
		text-align: center;
		padding: 0.6rem;
		border-radius: 4px;
		text-decoration: none;
		font-size: 0.9rem;
		margin-bottom: 0.5rem;
		border: 1px solid #ccc;
		color: #333;
	}
	.oauth-btn:hover {
		background: #f5f5f5;
	}
	.switch {
		text-align: center;
		font-size: 0.85rem;
		margin-top: 1rem;
		color: #666;
	}
	.switch a {
		color: #026aa7;
	}
</style>
