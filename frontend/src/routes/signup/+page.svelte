<script lang="ts">
	import { goto, invalidateAll } from '$app/navigation';
	import { resolve } from '$app/paths';
	import * as api from '$lib/api';

	let name = $state('');
	let email = $state('');
	let password = $state('');
	let error = $state('');

	async function handleSignup() {
		error = '';
		try {
			await api.signup(email, password, name);
			await invalidateAll();
			goto(resolve('/'));
		} catch (e: unknown) {
			error = e instanceof Error ? e.message : 'Signup failed';
		}
	}
</script>

<div class="auth-page">
	<div class="auth-card">
		<h1>Sign up</h1>
		{#if error}
			<div class="error">{error}</div>
		{/if}
		<form onsubmit={(e) => { e.preventDefault(); handleSignup(); }}>
			<input type="text" bind:value={name} placeholder="Name" />
			<input type="email" bind:value={email} placeholder="Email" required />
			<input type="password" bind:value={password} placeholder="Password (min 8 chars)" required minlength="8" />
			<button type="submit">Sign up</button>
		</form>
		<p class="switch">Already have an account? <a href={resolve('/login')}>Log in</a></p>
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
