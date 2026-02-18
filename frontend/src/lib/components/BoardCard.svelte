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
	role="listitem"
>
	{#if editing}
		<input bind:value={editTitle} class="edit-input" />
		<textarea
			bind:value={editDesc}
			class="edit-textarea"
			rows="3"
			placeholder="Add description..."></textarea>
		<div class="edit-actions">
			<button class="save-btn" onclick={save}>Save</button>
			<button class="cancel-link" onclick={cancel}>Cancel</button>
		</div>
	{:else}
		<div class="card-body">
			<div class="card-title">{card.title}</div>
			{#if card.description}
				<div class="card-desc">{card.description}</div>
			{/if}
		</div>
		<div class="card-actions">
			<button
				class="action-btn edit"
				aria-label="Edit"
				onclick={() => (editing = true)}
				title="Edit"
			>
				<svg width="11" height="11" viewBox="0 0 11 11" fill="none" aria-hidden="true">
					<path
						d="M7.5 1L10 3.5L3.5 10L1 10.5L1.5 8L7.5 1Z"
						stroke="currentColor"
						stroke-width="1.2"
						stroke-linejoin="round"
					/>
				</svg>
			</button>
			<button
				class="action-btn delete"
				aria-label="Delete"
				onclick={() => ondelete(card.id)}
				title="Delete"
			>
				<svg width="11" height="11" viewBox="0 0 11 11" fill="none" aria-hidden="true">
					<path
						d="M1.5 3h8M4.5 3V2h2v1M3.5 3l.5 6.5h3.5l.5-6.5"
						stroke="currentColor"
						stroke-width="1.2"
						stroke-linecap="round"
						stroke-linejoin="round"
					/>
				</svg>
			</button>
		</div>
	{/if}
</div>

<style>
	.card {
		background: white;
		border-radius: 8px;
		padding: 0.6rem 0.7rem;
		box-shadow:
			0 1px 3px rgba(0, 0, 0, 0.1),
			0 1px 1px rgba(0, 0, 0, 0.06);
		border: 1px solid rgba(0, 0, 0, 0.06);
		cursor: grab;
		transition:
			box-shadow 0.15s,
			transform 0.12s;
		position: relative;
	}
	.card:hover {
		box-shadow:
			0 4px 14px rgba(0, 0, 0, 0.12),
			0 1px 3px rgba(0, 0, 0, 0.08);
		transform: translateY(-1px);
	}
	.card:active {
		cursor: grabbing;
	}
	.card :global(.dragging) {
		opacity: 0.35;
		transform: scale(0.96);
	}
	.card-body {
		padding-right: 3rem;
	}
	.card-title {
		font-size: 0.86rem;
		font-weight: 500;
		color: #172b4d;
		line-height: 1.45;
	}
	.card-desc {
		font-size: 0.76rem;
		color: #6b778c;
		margin-top: 0.3rem;
		line-height: 1.4;
		display: -webkit-box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}
	.card-actions {
		position: absolute;
		top: 0.4rem;
		right: 0.4rem;
		display: flex;
		gap: 0.15rem;
		opacity: 0;
		transition: opacity 0.15s;
	}
	.card:hover .card-actions {
		opacity: 1;
	}
	.action-btn {
		width: 1.5rem;
		height: 1.5rem;
		display: flex;
		align-items: center;
		justify-content: center;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		color: #6b778c;
		background: rgba(240, 244, 248, 0.95);
		transition:
			background 0.12s,
			color 0.12s;
	}
	.action-btn.edit:hover {
		background: #e8f0fe;
		color: #0052cc;
	}
	.action-btn.delete:hover {
		background: #fde8e8;
		color: #d32f2f;
	}
	.edit-input,
	.edit-textarea {
		width: 100%;
		font-size: 0.86rem;
		padding: 0.4rem 0.5rem;
		border: 1.5px solid #c4d0dc;
		border-radius: 6px;
		box-sizing: border-box;
		color: #172b4d;
		font-family: inherit;
		outline: none;
		transition:
			border-color 0.12s,
			box-shadow 0.12s;
	}
	.edit-input:focus,
	.edit-textarea:focus {
		border-color: #0052cc;
		box-shadow: 0 0 0 2px rgba(0, 82, 204, 0.14);
	}
	.edit-textarea {
		margin-top: 0.4rem;
		resize: vertical;
	}
	.edit-actions {
		display: flex;
		align-items: center;
		gap: 0.4rem;
		margin-top: 0.4rem;
	}
	.save-btn {
		background: #0052cc;
		color: white;
		border: none;
		padding: 0.35rem 0.75rem;
		border-radius: 5px;
		cursor: pointer;
		font-size: 0.8rem;
		font-weight: 700;
		transition: background 0.12s;
	}
	.save-btn:hover {
		background: #0041a3;
	}
	.cancel-link {
		background: none;
		border: none;
		cursor: pointer;
		font-size: 0.8rem;
		color: #6b778c;
		padding: 0.35rem 0.4rem;
		border-radius: 5px;
		transition:
			background 0.12s,
			color 0.12s;
	}
	.cancel-link:hover {
		background: #f0f4f8;
		color: #172b4d;
	}
</style>
