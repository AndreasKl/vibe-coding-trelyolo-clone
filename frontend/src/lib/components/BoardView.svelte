<script lang="ts">
	import type { Board } from '$lib/types';
	import * as api from '$lib/api';
	import BoardColumn from './BoardColumn.svelte';

	let { board: initialBoard }: { board: Board } = $props();

	let refreshedBoard = $state<Board | null>(null);
	let board = $derived(refreshedBoard ?? initialBoard);
	let addingColumn = $state(false);
	let newColumnName = $state('');

	// Column drag state
	let draggingColumnId = $state<string | null>(null);
	let dragOverColumnId = $state<string | null>(null);

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

	function setFocus(el: HTMLInputElement) {
		el.focus();
	}

	function handleColumnDragStart(columnId: string) {
		draggingColumnId = columnId;
	}

	function handleColumnDragEnd() {
		draggingColumnId = null;
		dragOverColumnId = null;
	}

	function handleColumnDragOver(columnId: string) {
		if (draggingColumnId && draggingColumnId !== columnId) {
			dragOverColumnId = columnId;
		}
	}

	async function handleColumnDrop(targetColumnId: string) {
		if (!draggingColumnId || draggingColumnId === targetColumnId) {
			draggingColumnId = null;
			dragOverColumnId = null;
			return;
		}

		const cols = board.columns ?? [];
		const srcIdx = cols.findIndex((c) => c.id === draggingColumnId);
		const tgtIdx = cols.findIndex((c) => c.id === targetColumnId);

		if (srcIdx === -1 || tgtIdx === -1) {
			draggingColumnId = null;
			dragOverColumnId = null;
			return;
		}

		const targetPosition = cols[tgtIdx].position;
		const columnToMove = cols[srcIdx].id;

		draggingColumnId = null;
		dragOverColumnId = null;

		await api.moveColumn(columnToMove, targetPosition);
		await refresh();
	}
</script>

<div class="board-header">
	<h2>{board.name}</h2>
</div>

<div class="board">
	{#each board.columns ?? [] as column (column.id)}
		<BoardColumn
			{column}
			onrefresh={refresh}
			isDragging={draggingColumnId === column.id}
			isDragTarget={dragOverColumnId === column.id}
			isAnyColumnDragging={draggingColumnId !== null}
			oncolumndragstart={() => handleColumnDragStart(column.id)}
			oncolumndragend={handleColumnDragEnd}
			oncolumndragover={() => handleColumnDragOver(column.id)}
			oncolumndrop={() => handleColumnDrop(column.id)}
		/>
	{/each}

	<div class="add-column">
		{#if addingColumn}
			<input
				use:setFocus
				bind:value={newColumnName}
				placeholder="Column name..."
				onkeydown={(e) => {
					if (e.key === 'Enter') addColumn();
					if (e.key === 'Escape') addingColumn = false;
				}}
			/>
			<div class="add-actions">
				<button class="add-btn" onclick={addColumn}>Add column</button>
				<button class="cancel-btn" onclick={() => (addingColumn = false)}>&times;</button>
			</div>
		{:else}
			<button class="add-col-btn" onclick={() => (addingColumn = true)}>+ Add column</button>
		{/if}
	</div>
</div>

<style>
	.board-header {
		padding: 1rem 1.5rem 0.4rem;
	}
	.board-header h2 {
		margin: 0;
		color: white;
		font-size: 1.2rem;
		font-weight: 700;
		text-shadow: 0 1px 3px rgba(0, 0, 0, 0.35);
		letter-spacing: -0.2px;
	}
	.board {
		display: flex;
		gap: 0.75rem;
		padding: 0.5rem 1.5rem 1.5rem;
		overflow-x: auto;
		flex: 1;
		align-items: flex-start;
	}
	.board::-webkit-scrollbar {
		height: 8px;
	}
	.board::-webkit-scrollbar-track {
		background: rgba(0, 0, 0, 0.12);
		border-radius: 4px;
		margin: 0 1.5rem;
	}
	.board::-webkit-scrollbar-thumb {
		background: rgba(255, 255, 255, 0.28);
		border-radius: 4px;
	}
	.board::-webkit-scrollbar-thumb:hover {
		background: rgba(255, 255, 255, 0.42);
	}
	.add-column {
		min-width: 264px;
		background: rgba(255, 255, 255, 0.14);
		backdrop-filter: blur(6px);
		border-radius: 10px;
		padding: 0.6rem;
		border: 1px solid rgba(255, 255, 255, 0.18);
		flex-shrink: 0;
	}
	.add-col-btn {
		background: none;
		border: none;
		color: rgba(255, 255, 255, 0.9);
		cursor: pointer;
		padding: 0.5rem 0.6rem;
		width: 100%;
		text-align: left;
		border-radius: 6px;
		font-size: 0.88rem;
		font-weight: 500;
	}
	.add-col-btn:hover {
		background: rgba(255, 255, 255, 0.18);
		color: white;
	}
	.add-column input {
		width: 100%;
		padding: 0.5rem 0.6rem;
		border: none;
		border-radius: 6px;
		font-size: 0.88rem;
		box-sizing: border-box;
		background: white;
		color: #172b4d;
		outline: none;
		font-family: inherit;
		box-shadow: 0 1px 4px rgba(0, 0, 0, 0.2);
	}
	.add-column input:focus {
		box-shadow: 0 0 0 3px rgba(255, 255, 255, 0.4);
	}
	.add-actions {
		display: flex;
		align-items: center;
		gap: 0.3rem;
		margin-top: 0.4rem;
	}
	.add-btn {
		background: white;
		color: #0052cc;
		border: none;
		padding: 0.4rem 0.9rem;
		border-radius: 6px;
		cursor: pointer;
		font-size: 0.82rem;
		font-weight: 700;
		transition: background 0.12s;
	}
	.add-btn:hover {
		background: #f0f4f8;
	}
	.cancel-btn {
		background: none;
		border: none;
		cursor: pointer;
		font-size: 1.3rem;
		color: rgba(255, 255, 255, 0.75);
		line-height: 1;
		padding: 0 0.3rem;
		transition: color 0.12s;
	}
	.cancel-btn:hover {
		color: white;
	}
</style>
