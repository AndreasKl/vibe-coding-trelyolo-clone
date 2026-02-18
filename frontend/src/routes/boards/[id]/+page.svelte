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
		background: linear-gradient(160deg, #0f2d50 0%, #1a4a8a 45%, #1d3c7a 100%);
	}
</style>
