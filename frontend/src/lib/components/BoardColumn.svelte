<script lang="ts">
	import type { Column } from '$lib/types';
	import * as api from '$lib/api';
	import BoardCard from './BoardCard.svelte';

	let {
		column,
		onrefresh
	}: {
		column: Column;
		onrefresh: () => void;
	} = $props();

	let adding = $state(false);
	let newTitle = $state('');
	let dragOver = $state(false);

	async function addCard() {
		if (!newTitle.trim()) return;
		await api.createCard(column.id, newTitle.trim());
		newTitle = '';
		adding = false;
		onrefresh();
	}

	async function handleCardUpdate() {
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
		dragOver = true;
	}

	function handleDragLeave() {
		dragOver = false;
	}

	async function handleDrop(e: DragEvent) {
		e.preventDefault();
		dragOver = false;
		const data = JSON.parse(e.dataTransfer!.getData('text/plain'));
		const position = column.cards.length;
		await api.moveCard(data.id, column.id, position);
		onrefresh();
	}

	function setFocus(el: HTMLInputElement){
		el.focus();
	}
</script>

<div
	class="column"
	class:drag-over={dragOver}
	ondragover={handleDragOver}
	ondragleave={handleDragLeave}
	ondrop={handleDrop}
	role="list"
>
	<div class="column-header">
		<h3>{column.name}</h3>
		<button class="delete-col" onclick={handleDeleteColumn} title="Delete column">&times;</button>
	</div>

	<div class="cards">
		{#each column.cards as card (card.id)}
			<BoardCard {card} onupdate={handleCardUpdate} ondelete={handleCardDelete} />
		{/each}
	</div>

	{#if adding}
		<div class="add-form">
			<input
				use:setFocus
				bind:value={newTitle}
				placeholder="Card title..."
				onkeydown={(e) => e.key === 'Enter' && addCard()}
			/>
			<div class="add-actions">
				<button class="add-btn" onclick={addCard}>Add</button>
				<button class="cancel-btn" onclick={() => (adding = false)}>&times;</button>
			</div>
		</div>
	{:else}
		<button class="add-card-btn" onclick={() => (adding = true)}>+ Add a card</button>
	{/if}
</div>

<style>
	.column {
		background: #ebecf0;
		border-radius: 6px;
		padding: 0.5rem;
		min-width: 272px;
		max-width: 272px;
		max-height: calc(100vh - 7rem);
		display: flex;
		flex-direction: column;
		transition: background 0.15s;
	}
	.column.drag-over {
		background: #d4d8e0;
	}
	.column-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 0 0.25rem 0.5rem;
	}
	.column-header h3 {
		margin: 0;
		font-size: 0.9rem;
		font-weight: 600;
	}
	.delete-col {
		background: none;
		border: none;
		cursor: pointer;
		font-size: 1.2rem;
		color: #999;
		padding: 0 0.2rem;
		line-height: 1;
	}
	.delete-col:hover {
		color: #e44;
	}
	.cards {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		overflow-y: auto;
		flex: 1;
		min-height: 2rem;
	}
	.add-card-btn {
		background: none;
		border: none;
		text-align: left;
		padding: 0.4rem;
		color: #666;
		cursor: pointer;
		border-radius: 4px;
		font-size: 0.85rem;
		margin-top: 0.25rem;
	}
	.add-card-btn:hover {
		background: rgba(0, 0, 0, 0.05);
		color: #333;
	}
	.add-form {
		margin-top: 0.25rem;
	}
	.add-form input {
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
		color: #999;
	}
</style>
