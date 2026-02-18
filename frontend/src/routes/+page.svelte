<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import * as api from '$lib/api';
	import type { Board } from '$lib/types';

	let boards = $state<Board[]>([]);
	let newName = $state('');
	let loading = $state(true);
	let creating = $state(false);

	const gradients = [
		'linear-gradient(135deg, #0079bf 0%, #5067c5 100%)',
		'linear-gradient(135deg, #c87941 0%, #d4a017 100%)',
		'linear-gradient(135deg, #519839 0%, #2e7527 100%)',
		'linear-gradient(135deg, #ae4132 0%, #c0392b 100%)',
		'linear-gradient(135deg, #89609e 0%, #6c3483 100%)',
		'linear-gradient(135deg, #00aecc 0%, #0089a8 100%)',
		'linear-gradient(135deg, #c14f87 0%, #9b2c66 100%)',
		'linear-gradient(135deg, #4b9a6f 0%, #2e7a53 100%)'
	];

	function getBoardGradient(id: string): string {
		const hash = id.split('').reduce((acc, c) => acc + c.charCodeAt(0), 0);
		return gradients[hash % gradients.length];
	}

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
		creating = false;
		goto(resolve('/boards/[id]', { id: b.id }));
	}

	async function handleDelete(id: string) {
		await api.deleteBoard(id);
		boards = boards.filter((b) => b.id !== id);
	}

	function focusInput(el: HTMLInputElement) {
		el.focus();
	}

	load();
</script>

<div class="dashboard">
	<div class="page-header">
		<h1>My Boards</h1>
	</div>

	{#if loading}
		<div class="loading">
			<div class="spinner"></div>
		</div>
	{:else}
		<div class="boards-grid">
			{#each boards as board (board.id)}
				<a
					href={resolve('/boards/[id]', { id: board.id })}
					class="board-tile"
					style="background: {getBoardGradient(board.id)}"
				>
					<span class="board-name">{board.name}</span>
					<button
						class="delete-btn"
						onclick={(e) => {
							e.preventDefault();
							handleDelete(board.id);
						}}
						title="Delete board">&times;</button
					>
				</a>
			{/each}

			{#if creating}
				<div class="create-tile">
					<input
						use:focusInput
						bind:value={newName}
						placeholder="Board name..."
						onkeydown={(e) => {
							if (e.key === 'Enter') createBoard();
							if (e.key === 'Escape') {
								creating = false;
								newName = '';
							}
						}}
					/>
					<div class="create-actions">
						<button class="create-btn" onclick={createBoard}>Create</button>
						<button
							class="cancel-btn"
							onclick={() => {
								creating = false;
								newName = '';
							}}>Cancel</button
						>
					</div>
				</div>
			{:else}
				<button class="new-board-tile" onclick={() => (creating = true)}>
					<span class="plus-icon">+</span>
					<span>Create board</span>
				</button>
			{/if}
		</div>
	{/if}
</div>

<style>
	.dashboard {
		padding: 2.5rem 3rem;
	}
	.page-header {
		margin-bottom: 2rem;
	}
	h1 {
		font-size: 1.45rem;
		font-weight: 700;
		color: #172b4d;
		margin: 0;
		letter-spacing: -0.3px;
	}
	.loading {
		display: flex;
		justify-content: center;
		padding: 4rem;
	}
	.spinner {
		width: 2rem;
		height: 2rem;
		border: 3px solid #e2e8f0;
		border-top-color: #0052cc;
		border-radius: 50%;
		animation: spin 0.65s linear infinite;
	}
	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}
	.boards-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(200px, 220px));
		gap: 1.25rem;
	}
	.board-tile {
		border-radius: 10px;
		height: 115px;
		padding: 1rem;
		text-decoration: none;
		position: relative;
		display: flex;
		align-items: flex-start;
		box-shadow:
			0 2px 8px rgba(0, 0, 0, 0.18),
			0 1px 2px rgba(0, 0, 0, 0.1);
		transition:
			transform 0.15s,
			box-shadow 0.15s;
		overflow: hidden;
	}
	.board-tile::after {
		content: '';
		position: absolute;
		inset: 0;
		background: linear-gradient(135deg, rgba(255, 255, 255, 0.12) 0%, transparent 60%);
		pointer-events: none;
	}
	.board-tile:hover {
		transform: translateY(-2px);
		box-shadow:
			0 8px 24px rgba(0, 0, 0, 0.22),
			0 2px 6px rgba(0, 0, 0, 0.12);
	}
	.board-name {
		color: white;
		font-weight: 700;
		font-size: 0.95rem;
		line-height: 1.35;
		text-shadow: 0 1px 2px rgba(0, 0, 0, 0.25);
		flex: 1;
	}
	.delete-btn {
		position: absolute;
		top: 0.5rem;
		right: 0.5rem;
		background: rgba(0, 0, 0, 0.18);
		border: none;
		color: rgba(255, 255, 255, 0.65);
		font-size: 1rem;
		cursor: pointer;
		padding: 0.15rem 0.4rem;
		border-radius: 4px;
		line-height: 1;
		opacity: 0;
		transition:
			opacity 0.15s,
			background 0.15s,
			color 0.15s;
	}
	.board-tile:hover .delete-btn {
		opacity: 1;
	}
	.delete-btn:hover {
		background: rgba(0, 0, 0, 0.38);
		color: white;
	}
	.new-board-tile {
		border-radius: 10px;
		height: 115px;
		padding: 1rem;
		background: #f0f4f8;
		border: 2px dashed #c4d0dc;
		cursor: pointer;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 0.4rem;
		color: #5e7a90;
		font-size: 0.88rem;
		font-weight: 500;
		transition:
			background 0.15s,
			border-color 0.15s,
			color 0.15s;
	}
	.new-board-tile:hover {
		background: #e4edf6;
		border-color: #0052cc;
		color: #0052cc;
	}
	.plus-icon {
		font-size: 1.6rem;
		line-height: 1;
		font-weight: 300;
	}
	.create-tile {
		border-radius: 10px;
		height: 115px;
		padding: 0.85rem;
		background: white;
		border: 2px solid #0052cc;
		box-shadow: 0 0 0 3px rgba(0, 82, 204, 0.15);
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		box-sizing: border-box;
	}
	.create-tile input {
		width: 100%;
		padding: 0.45rem 0.6rem;
		border: 1.5px solid #c4d0dc;
		border-radius: 6px;
		font-size: 0.88rem;
		color: #172b4d;
		outline: none;
		background: #f8fafc;
		box-sizing: border-box;
		font-family: inherit;
		transition:
			border-color 0.12s,
			box-shadow 0.12s;
	}
	.create-tile input:focus {
		border-color: #0052cc;
		background: white;
		box-shadow: 0 0 0 2px rgba(0, 82, 204, 0.12);
	}
	.create-actions {
		display: flex;
		gap: 0.4rem;
		align-items: center;
	}
	.create-btn {
		background: #0052cc;
		color: white;
		border: none;
		padding: 0.4rem 0.9rem;
		border-radius: 6px;
		cursor: pointer;
		font-size: 0.82rem;
		font-weight: 600;
		transition: background 0.12s;
	}
	.create-btn:hover {
		background: #0041a3;
	}
	.cancel-btn {
		background: none;
		border: none;
		cursor: pointer;
		font-size: 0.82rem;
		color: #5e6c84;
		padding: 0.4rem 0.5rem;
		border-radius: 5px;
		transition:
			background 0.12s,
			color 0.12s;
	}
	.cancel-btn:hover {
		background: #f0f4f8;
		color: #172b4d;
	}
</style>
