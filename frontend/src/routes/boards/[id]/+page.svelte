<script lang="ts">
	import { page } from '$app/state';
	import BoardView from '$lib/components/BoardView.svelte';

	const data = $derived(page.data);
</script>

{#if data.board}
	<div class="board-page">
		<BoardView board={data.board} />
	</div>
{/if}

<style>
	.board-page {
		flex: 1;
		display: flex;
		flex-direction: column;
		background: linear-gradient(135deg, #0079bf, #5067c5);
	}
</style>
