<script lang="ts">
	import type { Column } from '$lib/types';
	import * as api from '$lib/api';
	import BoardCard from './BoardCard.svelte';

	let {
		column,
		onrefresh,
		isDragging = false,
		isDragTarget = false,
		isAnyColumnDragging = false,
		oncolumndragstart,
		oncolumndragend,
		oncolumndragover,
		oncolumndrop
	}: {
		column: Column;
		onrefresh: () => void;
		isDragging?: boolean;
		isDragTarget?: boolean;
		isAnyColumnDragging?: boolean;
		oncolumndragstart?: () => void;
		oncolumndragend?: () => void;
		oncolumndragover?: () => void;
		oncolumndrop?: () => void;
	} = $props();

	let adding = $state(false);
	let newTitle = $state('');
	let cardDragOver = $state(false);

	async function addCard() {
		if (!newTitle.trim()) return;
		await api.createCard(column.id, newTitle.trim());
		newTitle = '';
		adding = false;
		onrefresh();
	}

	async function handleCardDelete(id: string) {
		await api.deleteCard(id);
		onrefresh();
	}

	async function handleDeleteColumn() {
		await api.deleteColumn(column.id);
		onrefresh();
	}

	function handleDragOver(e: DragEvent) {
		e.preventDefault();
		e.dataTransfer!.dropEffect = 'move';
		if (isAnyColumnDragging) {
			if (!isDragging) {
				oncolumndragover?.();
			}
		} else {
			cardDragOver = true;
		}
	}

	function handleDragLeave(e: DragEvent) {
		const rel = e.relatedTarget as HTMLElement | null;
		if (rel && (e.currentTarget as HTMLElement).contains(rel)) return;
		cardDragOver = false;
	}

	async function handleDrop(e: DragEvent) {
		e.preventDefault();
		cardDragOver = false;

		if (isAnyColumnDragging) {
			oncolumndrop?.();
			return;
		}

		const raw = e.dataTransfer!.getData('text/plain');
		if (!raw) return;
		const data = JSON.parse(raw);
		const position = column.cards.length;
		await api.moveCard(data.id, column.id, position);
		onrefresh();
	}

	function handleHeaderDragStart(e: DragEvent) {
		e.dataTransfer!.effectAllowed = 'move';
		e.dataTransfer!.setData('text/plain', JSON.stringify({ type: 'column', id: column.id }));
		oncolumndragstart?.();
	}

	function handleHeaderDragEnd() {
		oncolumndragend?.();
	}

	function setFocus(el: HTMLInputElement) {
		el.focus();
	}
</script>

<div
	class="column"
	class:dragging={isDragging}
	class:drag-target={isDragTarget}
	class:card-drag-over={cardDragOver && !isAnyColumnDragging}
	ondragover={handleDragOver}
	ondragleave={handleDragLeave}
	ondrop={handleDrop}
	role="list"
>
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div
		class="column-header"
		draggable="true"
		ondragstart={handleHeaderDragStart}
		ondragend={handleHeaderDragEnd}
	>
		<div class="drag-handle" title="Drag to reorder">
			<span></span><span></span>
			<span></span><span></span>
			<span></span><span></span>
		</div>
		<h3>{column.name}</h3>
		<div class="header-right">
			<span class="card-count">{column.cards.length}</span>
			<button class="delete-col" onclick={handleDeleteColumn} title="Delete column"
				>&times;</button
			>
		</div>
	</div>

	<div class="cards">
		{#each column.cards as card (card.id)}
			<BoardCard {card} onupdate={onrefresh} ondelete={handleCardDelete} />
		{/each}
	</div>

	{#if adding}
		<div class="add-form">
			<input
				use:setFocus
				bind:value={newTitle}
				placeholder="Card title..."
				onkeydown={(e) => {
					if (e.key === 'Enter') addCard();
					if (e.key === 'Escape') adding = false;
				}}
			/>
			<div class="add-actions">
				<button class="add-btn" onclick={addCard}>Add card</button>
				<button class="cancel-btn" onclick={() => (adding = false)}>&times;</button>
			</div>
		</div>
	{:else}
		<button class="add-card-btn" onclick={() => (adding = true)}>+ Add a card</button>
	{/if}
</div>

<style>
	.column {
		background: #f8f9fb;
		border-radius: 10px;
		min-width: 272px;
		max-width: 272px;
		max-height: calc(100vh - 8rem);
		display: flex;
		flex-direction: column;
		box-shadow:
			0 2px 8px rgba(0, 0, 0, 0.14),
			0 1px 2px rgba(0, 0, 0, 0.08);
		border: 2px solid transparent;
		transition:
			opacity 0.18s,
			box-shadow 0.18s,
			border-color 0.18s,
			transform 0.18s;
		overflow: hidden;
		flex-shrink: 0;
	}
	.column.dragging {
		opacity: 0.4;
		transform: scale(0.97);
		box-shadow: 0 8px 28px rgba(0, 0, 0, 0.28);
	}
	.column.drag-target {
		border-color: rgba(255, 255, 255, 0.7);
		box-shadow:
			0 0 0 2px rgba(255, 255, 255, 0.5),
			0 6px 20px rgba(0, 0, 0, 0.2);
	}
	.column.card-drag-over {
		background: #eef2f8;
		box-shadow:
			0 0 0 2px #4a90e2,
			0 4px 16px rgba(74, 144, 226, 0.18);
	}
	.column-header {
		display: flex;
		align-items: center;
		gap: 0.4rem;
		padding: 0.65rem 0.75rem 0.5rem;
		cursor: grab;
		background: #f0f2f6;
		border-bottom: 1px solid rgba(0, 0, 0, 0.06);
		user-select: none;
	}
	.column-header:active {
		cursor: grabbing;
	}
	.drag-handle {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 2.5px 3px;
		width: 12px;
		flex-shrink: 0;
		opacity: 0;
		transition: opacity 0.15s;
	}
	.column-header:hover .drag-handle {
		opacity: 1;
	}
	.drag-handle span {
		display: block;
		width: 3px;
		height: 3px;
		background: #8898aa;
		border-radius: 50%;
	}
	h3 {
		margin: 0;
		font-size: 0.86rem;
		font-weight: 700;
		color: #1a2537;
		flex: 1;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
		letter-spacing: 0.1px;
	}
	.header-right {
		display: flex;
		align-items: center;
		gap: 0.3rem;
	}
	.card-count {
		background: rgba(0, 0, 0, 0.09);
		color: #5e6c84;
		font-size: 0.7rem;
		font-weight: 700;
		padding: 0.1rem 0.42rem;
		border-radius: 10px;
		min-width: 1.4rem;
		text-align: center;
	}
	.delete-col {
		background: none;
		border: none;
		cursor: pointer;
		font-size: 1.1rem;
		color: #8898aa;
		padding: 0;
		line-height: 1;
		width: 1.5rem;
		height: 1.5rem;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 4px;
		transition:
			background 0.12s,
			color 0.12s;
	}
	.delete-col:hover {
		background: #fde8e8;
		color: #d73535;
	}
	.cards {
		display: flex;
		flex-direction: column;
		gap: 0.4rem;
		overflow-y: auto;
		flex: 1;
		min-height: 0.5rem;
		padding: 0.55rem 0.6rem;
	}
	.cards::-webkit-scrollbar {
		width: 4px;
	}
	.cards::-webkit-scrollbar-track {
		background: transparent;
	}
	.cards::-webkit-scrollbar-thumb {
		background: rgba(0, 0, 0, 0.12);
		border-radius: 4px;
	}
	.add-card-btn {
		background: none;
		border: none;
		text-align: left;
		padding: 0.45rem 0.75rem;
		color: #5e6c84;
		cursor: pointer;
		font-size: 0.83rem;
		margin: 0.2rem 0.4rem 0.5rem;
		border-radius: 6px;
		transition:
			background 0.12s,
			color 0.12s;
	}
	.add-card-btn:hover {
		background: rgba(0, 0, 0, 0.06);
		color: #1a2537;
	}
	.add-form {
		padding: 0.4rem 0.6rem 0.6rem;
	}
	.add-form input {
		width: 100%;
		padding: 0.45rem 0.6rem;
		border: 1.5px solid #c4d0dc;
		border-radius: 6px;
		font-size: 0.86rem;
		box-sizing: border-box;
		color: #172b4d;
		background: white;
		outline: none;
		font-family: inherit;
		transition:
			border-color 0.12s,
			box-shadow 0.12s;
	}
	.add-form input:focus {
		border-color: #0052cc;
		box-shadow: 0 0 0 2px rgba(0, 82, 204, 0.14);
	}
	.add-actions {
		display: flex;
		align-items: center;
		gap: 0.3rem;
		margin-top: 0.4rem;
	}
	.add-btn {
		background: #0052cc;
		color: white;
		border: none;
		padding: 0.38rem 0.8rem;
		border-radius: 6px;
		cursor: pointer;
		font-size: 0.8rem;
		font-weight: 700;
		transition: background 0.12s;
	}
	.add-btn:hover {
		background: #0041a3;
	}
	.cancel-btn {
		background: none;
		border: none;
		cursor: pointer;
		font-size: 1.3rem;
		color: #8898aa;
		line-height: 1;
		padding: 0 0.2rem;
		border-radius: 4px;
		transition: color 0.12s;
	}
	.cancel-btn:hover {
		color: #3d4f5c;
	}
</style>
