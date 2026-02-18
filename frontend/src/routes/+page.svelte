<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import * as api from '$lib/api';
	import type { Board } from '$lib/types';

	let boards = $state<Board[]>([]);
	let newName = $state('');
	let loading = $state(true);

	async function load() {
		try {
			boards = await api.listBoards();
		} catch {
			// not authenticated, layout will redirect
		}
		loading = false;
	}

	async function createBoard() {
		if (!newName.trim()) return;
		const b = await api.createBoard(newName.trim());
		newName = '';
		goto(resolve('/boards/[id]', { id: b.id }));
	}

	async function handleDelete(id: string) {
		await api.deleteBoard(id);
		boards = boards.filter((b) => b.id !== id);
	}

	load();
</script>

<div class="dashboard">
	<h1>Your Boards</h1>

	{#if loading}
		<p>Loading...</p>
	{:else}
		<div class="create-board">
			<input
				bind:value={newName}
				placeholder="New board name..."
				onkeydown={(e) => e.key === 'Enter' && createBoard()}
			/>
			<button onclick={createBoard}>Create Board</button>
		</div>

		<div class="boards-grid">
			{#each boards as board (board.id)}
				<div class="board-tile">
					<a href={resolve('/boards/[id]', { id: board.id })} class="board-link">{board.name}</a>
					<button class="delete-btn" onclick={() => handleDelete(board.id)}>&times;</button>
				</div>
			{:else}
				<p class="empty">No boards yet. Create one above!</p>
			{/each}
		</div>
	{/if}
</div>

<style>
	.dashboard {
		max-width: 800px;
		margin: 2rem auto;
		padding: 0 1.5rem;
	}
	h1 {
		font-size: 1.4rem;
		margin-bottom: 1.5rem;
	}
	.create-board {
		display: flex;
		gap: 0.5rem;
		margin-bottom: 2rem;
	}
	.create-board input {
		flex: 1;
		padding: 0.5rem;
		border: 1px solid #ccc;
		border-radius: 4px;
		font-size: 0.9rem;
	}
	.create-board button {
		background: #026aa7;
		color: white;
		border: none;
		padding: 0.5rem 1rem;
		border-radius: 4px;
		cursor: pointer;
		font-size: 0.9rem;
	}
	.create-board button:hover {
		background: #015d93;
	}
	.boards-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
		gap: 1rem;
	}
	.board-tile {
		background: #026aa7;
		border-radius: 6px;
		padding: 1rem;
		position: relative;
		height: 100px;
		display: flex;
		align-items: flex-start;
	}
	.board-link {
		color: white;
		text-decoration: none;
		font-weight: 600;
		font-size: 1rem;
		flex: 1;
	}
	.board-link:hover {
		text-decoration: underline;
	}
	.delete-btn {
		position: absolute;
		top: 0.5rem;
		right: 0.5rem;
		background: none;
		border: none;
		color: rgba(255, 255, 255, 0.6);
		font-size: 1.2rem;
		cursor: pointer;
		padding: 0;
		line-height: 1;
	}
	.delete-btn:hover {
		color: white;
	}
	.empty {
		color: #999;
		grid-column: 1 / -1;
	}
</style>
