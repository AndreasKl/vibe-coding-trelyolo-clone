<script lang="ts">
	import type { Card } from '$lib/types';
	import * as api from '$lib/api';

	let {
		card,
		onupdate,
		ondelete
	}: {
		card: Card;
		onupdate: () => void;
		ondelete: (id: string) => void;
	} = $props();

	let editing = $state(false);
	let editTitle = $state('');
	let editDesc = $state('');

	$effect(() => {
		editTitle = card.title;
		editDesc = card.description;
	});

	function handleDragStart(e: DragEvent) {
		e.dataTransfer!.effectAllowed = 'move';
		e.dataTransfer!.setData('text/plain', JSON.stringify({ id: card.id, columnId: card.column_id }));
		(e.target as HTMLElement).classList.add('dragging');
	}

	function handleDragEnd(e: DragEvent) {
		(e.target as HTMLElement).classList.remove('dragging');
	}

	async function save() {
		await api.updateCard(card.id, { title: editTitle, description: editDesc });
		editing = false;
		onupdate();
	}

	function cancel() {
		editTitle = card.title;
		editDesc = card.description;
		editing = false;
	}
</script>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
	class="card"
	draggable="true"
	ondragstart={handleDragStart}
	ondragend={handleDragEnd}
>
	{#if editing}
		<input bind:value={editTitle} class="edit-input" />
		<textarea bind:value={editDesc} class="edit-textarea" rows="2"></textarea>
		<div class="edit-actions">
			<button class="save-btn" onclick={save}>Save</button>
			<button class="cancel-btn" onclick={cancel}>Cancel</button>
		</div>
	{:else}
		<div class="card-title">{card.title}</div>
		{#if card.description}
			<div class="card-desc">{card.description}</div>
		{/if}
		<div class="card-actions">
			<button onclick={() => (editing = true)}>Edit</button>
			<button onclick={() => ondelete(card.id)}>Delete</button>
		</div>
	{/if}
</div>

<style>
	.card {
		background: white;
		border-radius: 4px;
		padding: 0.5rem;
		box-shadow: 0 1px 2px rgba(0, 0, 0, 0.15);
		cursor: grab;
		transition: opacity 0.2s;
	}
	.card:active {
		cursor: grabbing;
	}
	.card :global(.dragging) {
		opacity: 0.4;
	}
	.card-title {
		font-size: 0.9rem;
		font-weight: 500;
	}
	.card-desc {
		font-size: 0.8rem;
		color: #666;
		margin-top: 0.25rem;
	}
	.card-actions {
		display: flex;
		gap: 0.25rem;
		margin-top: 0.4rem;
		opacity: 0;
		transition: opacity 0.15s;
	}
	.card:hover .card-actions {
		opacity: 1;
	}
	.card-actions button {
		font-size: 0.7rem;
		padding: 0.15rem 0.4rem;
		background: #f0f0f0;
		border: none;
		border-radius: 3px;
		cursor: pointer;
		color: #666;
	}
	.card-actions button:hover {
		background: #e0e0e0;
	}
	.edit-input,
	.edit-textarea {
		width: 100%;
		font-size: 0.85rem;
		padding: 0.3rem;
		border: 1px solid #ccc;
		border-radius: 3px;
		box-sizing: border-box;
	}
	.edit-textarea {
		margin-top: 0.3rem;
		resize: vertical;
	}
	.edit-actions {
		display: flex;
		gap: 0.25rem;
		margin-top: 0.3rem;
	}
	.save-btn {
		background: #5aac44;
		color: white;
		border: none;
		padding: 0.3rem 0.6rem;
		border-radius: 3px;
		cursor: pointer;
		font-size: 0.8rem;
	}
	.cancel-btn {
		background: none;
		border: none;
		cursor: pointer;
		font-size: 0.8rem;
		color: #666;
	}
</style>
