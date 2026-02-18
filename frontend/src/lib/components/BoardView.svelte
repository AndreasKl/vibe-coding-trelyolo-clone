<script lang="ts">
	import type { Board } from '$lib/types';
	import * as api from '$lib/api';
	import BoardColumn from './BoardColumn.svelte';

	let { board: initialBoard }: { board: Board } = $props();

	let refreshedBoard = $state<Board | null>(null);
	let board = $derived(refreshedBoard ?? initialBoard);
	let addingColumn = $state(false);
	let newColumnName = $state('');

	async function refresh() {
		refreshedBoard = await api.getBoard(board.id);
	}

	async function addColumn() {
		if (!newColumnName.trim()) return;
		await api.createColumn(board.id, newColumnName.trim());
		newColumnName = '';
		addingColumn = false;
		await refresh();
	}

	function setFocus(el: HTMLInputElement){
		el.focus();
	}
</script>

<div class="board-header">
	<h2>{board.name}</h2>
</div>

<div class="board">
	{#each board.columns ?? [] as column (column.id)}
		<BoardColumn {column} onrefresh={refresh} />
	{/each}

	<div class="add-column">
		{#if addingColumn}
			<input
				use:setFocus
				bind:value={newColumnName}
				placeholder="Column name..."
				onkeydown={(e) => e.key === 'Enter' && addColumn()}
			/>
			<div class="add-actions">
				<button class="add-btn" onclick={addColumn}>Add</button>
				<button class="cancel-btn" onclick={() => (addingColumn = false)}>&times;</button>
			</div>
		{:else}
			<button class="add-col-btn" onclick={() => (addingColumn = true)}>+ Add column</button>
		{/if}
	</div>
</div>

<style>
	.board-header {
		padding: 0.75rem 1.5rem;
	}
	.board-header h2 {
		margin: 0;
		color: white;
		font-size: 1.2rem;
		text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);
	}
	.board {
		display: flex;
		gap: 0.75rem;
		padding: 0 1.5rem 1.5rem;
		overflow-x: auto;
		flex: 1;
		align-items: flex-start;
	}
	.add-column {
		min-width: 272px;
		background: rgba(255, 255, 255, 0.24);
		border-radius: 6px;
		padding: 0.5rem;
	}
	.add-col-btn {
		background: none;
		border: none;
		color: white;
		cursor: pointer;
		padding: 0.4rem;
		width: 100%;
		text-align: left;
		border-radius: 4px;
		font-size: 0.85rem;
	}
	.add-col-btn:hover {
		background: rgba(255, 255, 255, 0.15);
	}
	.add-column input {
		width: 100%;
		padding: 0.4rem;
		border: 1px solid #ccc;
		border-radius: 4px;
		font-size: 0.85rem;
		box-sizing: border-box;
	}
	.add-actions {
		display: flex;
		align-items: center;
		gap: 0.25rem;
		margin-top: 0.3rem;
	}
	.add-btn {
		background: #5aac44;
		color: white;
		border: none;
		padding: 0.35rem 0.8rem;
		border-radius: 4px;
		cursor: pointer;
		font-size: 0.85rem;
	}
	.cancel-btn {
		background: none;
		border: none;
		cursor: pointer;
		font-size: 1.2rem;
		color: white;
	}
</style>
