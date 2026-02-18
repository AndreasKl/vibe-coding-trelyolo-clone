import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import BoardColumn from './BoardColumn.svelte';
import type { Column } from '$lib/types';

const { mockCreateCard, mockDeleteCard, mockUpdateCard, mockDeleteColumn, mockMoveCard } =
	vi.hoisted(() => ({
		mockCreateCard: vi.fn(),
		mockDeleteCard: vi.fn(),
		mockUpdateCard: vi.fn(),
		mockDeleteColumn: vi.fn(),
		mockMoveCard: vi.fn()
	}));

vi.mock('$lib/api', () => ({
	createCard: mockCreateCard,
	deleteCard: mockDeleteCard,
	updateCard: mockUpdateCard,
	deleteColumn: mockDeleteColumn,
	moveCard: mockMoveCard
}));

const column: Column = {
	id: 'col1',
	board_id: 'b1',
	name: 'To Do',
	position: 0,
	created_at: '2024-01-01',
	cards: [
		{ id: 'card1', column_id: 'col1', title: 'First Card', description: '', position: 0, created_at: '' },
		{ id: 'card2', column_id: 'col1', title: 'Second Card', description: 'details', position: 1, created_at: '' }
	]
};

describe('BoardColumn', () => {
	beforeEach(() => {
		mockCreateCard.mockReset().mockResolvedValue({ id: 'new', title: 'New', description: '', column_id: 'col1', position: 2, created_at: '' });
		mockDeleteCard.mockReset().mockResolvedValue(undefined);
		mockDeleteColumn.mockReset().mockResolvedValue(undefined);
		mockUpdateCard.mockReset().mockResolvedValue({});
		mockMoveCard.mockReset().mockResolvedValue({});
	});

	it('renders the column name', () => {
		render(BoardColumn, { column, onrefresh: vi.fn() });
		expect(screen.getByText('To Do')).toBeInTheDocument();
	});

	it('renders all cards in the column', () => {
		render(BoardColumn, { column, onrefresh: vi.fn() });
		expect(screen.getByText('First Card')).toBeInTheDocument();
		expect(screen.getByText('Second Card')).toBeInTheDocument();
	});

	it('shows add card input when "+ Add a card" is clicked', async () => {
		render(BoardColumn, { column, onrefresh: vi.fn() });
		await fireEvent.click(screen.getByText('+ Add a card'));
		expect(screen.getByPlaceholderText('Card title...')).toBeInTheDocument();
	});

	it('creates a card and calls onrefresh', async () => {
		const onrefresh = vi.fn();
		render(BoardColumn, { column, onrefresh });
		await fireEvent.click(screen.getByText('+ Add a card'));
		const input = screen.getByPlaceholderText('Card title...');
		input.value = 'New Card';
		await fireEvent.input(input);
		await fireEvent.click(screen.getByText('Add'));
		await waitFor(() => {
			expect(mockCreateCard).toHaveBeenCalledWith('col1', 'New Card');
			expect(onrefresh).toHaveBeenCalled();
		});
	});

	it('hides add form when cancel button is clicked', async () => {
		render(BoardColumn, { column, onrefresh: vi.fn() });
		await fireEvent.click(screen.getByText('+ Add a card'));
		const cancelBtns = screen.getAllByText('×');
		await fireEvent.click(cancelBtns[cancelBtns.length - 1]);
		expect(screen.queryByPlaceholderText('Card title...')).not.toBeInTheDocument();
	});

	it('deletes the column and calls onrefresh', async () => {
		const onrefresh = vi.fn();
		render(BoardColumn, { column, onrefresh });
		await fireEvent.click(screen.getByTitle('Delete column'));
		await waitFor(() => {
			expect(mockDeleteColumn).toHaveBeenCalledWith('col1');
			expect(onrefresh).toHaveBeenCalled();
		});
	});

	it('deletes a card and calls onrefresh', async () => {
		const onrefresh = vi.fn();
		render(BoardColumn, { column, onrefresh });
		const deleteButtons = screen.getAllByText('Delete');
		await fireEvent.click(deleteButtons[0]);
		await waitFor(() => {
			expect(mockDeleteCard).toHaveBeenCalledWith('card1');
			expect(onrefresh).toHaveBeenCalled();
		});
	});
});
