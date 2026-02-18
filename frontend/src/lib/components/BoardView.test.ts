import { render, screen, fireEvent, waitFor, within } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import BoardView from './BoardView.svelte';
import type { Board } from '$lib/types';

const { mockGetBoard, mockCreateColumn, mockDeleteColumn, mockCreateCard, mockDeleteCard, mockUpdateCard, mockMoveCard, mockMoveColumn } =
	vi.hoisted(() => ({
		mockGetBoard: vi.fn(),
		mockCreateColumn: vi.fn(),
		mockDeleteColumn: vi.fn(),
		mockCreateCard: vi.fn(),
		mockDeleteCard: vi.fn(),
		mockUpdateCard: vi.fn(),
		mockMoveCard: vi.fn(),
		mockMoveColumn: vi.fn()
	}));

vi.mock('$lib/api', () => ({
	getBoard: mockGetBoard,
	createColumn: mockCreateColumn,
	deleteColumn: mockDeleteColumn,
	createCard: mockCreateCard,
	deleteCard: mockDeleteCard,
	updateCard: mockUpdateCard,
	moveCard: mockMoveCard,
	moveColumn: mockMoveColumn
}));

const board: Board = {
	id: 'b1',
	user_id: '1',
	name: 'My Board',
	created_at: '2024-01-01',
	columns: [
		{ id: 'col1', board_id: 'b1', name: 'Backlog', position: 0, created_at: '', cards: [] },
		{ id: 'col2', board_id: 'b1', name: 'In Progress', position: 1, created_at: '', cards: [] }
	]
};

describe('BoardView', () => {
	beforeEach(() => {
		mockGetBoard.mockReset().mockResolvedValue(board);
		mockCreateColumn.mockReset().mockResolvedValue({ id: 'col-new', board_id: 'b1', name: 'Done', position: 2, created_at: '', cards: [] });
		mockDeleteColumn.mockReset().mockResolvedValue(undefined);
		mockCreateCard.mockReset().mockResolvedValue({});
		mockDeleteCard.mockReset().mockResolvedValue(undefined);
		mockUpdateCard.mockReset().mockResolvedValue({});
		mockMoveCard.mockReset().mockResolvedValue({});
		mockMoveColumn.mockReset().mockResolvedValue({});
	});

	it('renders the board name', () => {
		render(BoardView, { board });
		expect(screen.getByText('My Board')).toBeInTheDocument();
	});

	it('renders all existing columns', () => {
		render(BoardView, { board });
		expect(screen.getByText('Backlog')).toBeInTheDocument();
		expect(screen.getByText('In Progress')).toBeInTheDocument();
	});

	it('shows add column input when "+ Add column" is clicked', async () => {
		render(BoardView, { board });
		await fireEvent.click(screen.getByText('+ Add column'));
		expect(screen.getByPlaceholderText('Column name...')).toBeInTheDocument();
	});

	it('creates a column and refreshes when Add column is clicked', async () => {
		render(BoardView, { board });
		await fireEvent.click(screen.getByText('+ Add column'));
		const input = screen.getByPlaceholderText('Column name...');
		input.value = 'Done';
		await fireEvent.input(input);
		await fireEvent.click(screen.getByText('Add column'));
		await waitFor(() => {
			expect(mockCreateColumn).toHaveBeenCalledWith('b1', 'Done');
			expect(mockGetBoard).toHaveBeenCalledWith('b1');
		});
	});

	it('hides add column form when cancel (×) is clicked', async () => {
		render(BoardView, { board });
		await fireEvent.click(screen.getByText('+ Add column'));
		const addColumnDiv = screen.getByPlaceholderText('Column name...').closest('div') as HTMLElement;
		await fireEvent.click(within(addColumnDiv).getByText('×'));
		expect(screen.queryByPlaceholderText('Column name...')).not.toBeInTheDocument();
	});

	it('does not call createColumn when input is empty', async () => {
		render(BoardView, { board });
		await fireEvent.click(screen.getByText('+ Add column'));
		await fireEvent.click(screen.getByText('Add column'));
		expect(mockCreateColumn).not.toHaveBeenCalled();
	});
});
