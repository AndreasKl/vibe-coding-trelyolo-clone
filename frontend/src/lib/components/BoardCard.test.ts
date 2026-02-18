import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import BoardCard from './BoardCard.svelte';
import type { Card } from '$lib/types';

const { mockUpdateCard } = vi.hoisted(() => ({
	mockUpdateCard: vi.fn()
}));

vi.mock('$lib/api', () => ({ updateCard: mockUpdateCard }));

const card: Card = {
	id: 'card1',
	column_id: 'col1',
	title: 'Test Card',
	description: 'A description',
	position: 0,
	created_at: '2024-01-01'
};

describe('BoardCard', () => {
	beforeEach(() => {
		mockUpdateCard.mockReset().mockResolvedValue({ ...card });
	});

	it('renders the card title', () => {
		render(BoardCard, { card, onupdate: vi.fn(), ondelete: vi.fn() });
		expect(screen.getByText('Test Card')).toBeInTheDocument();
	});

	it('renders the description when present', () => {
		render(BoardCard, { card, onupdate: vi.fn(), ondelete: vi.fn() });
		expect(screen.getByText('A description')).toBeInTheDocument();
	});

	it('does not render description when empty', () => {
		render(BoardCard, { card: { ...card, description: '' }, onupdate: vi.fn(), ondelete: vi.fn() });
		expect(screen.queryByText('A description')).not.toBeInTheDocument();
	});

	it('shows edit form with current values when Edit is clicked', async () => {
		render(BoardCard, { card, onupdate: vi.fn(), ondelete: vi.fn() });
		await fireEvent.click(screen.getByRole('button', { name: 'Edit' }));
		expect(screen.getByDisplayValue('Test Card')).toBeInTheDocument();
		expect(screen.getByDisplayValue('A description')).toBeInTheDocument();
	});

	it('calls ondelete with the card id when Delete is clicked', async () => {
		const ondelete = vi.fn();
		render(BoardCard, { card, onupdate: vi.fn(), ondelete });
		await fireEvent.click(screen.getByRole('button', { name: 'Delete' }));
		expect(ondelete).toHaveBeenCalledWith('card1');
	});

	it('calls updateCard and onupdate when Save is clicked', async () => {
		const onupdate = vi.fn();
		render(BoardCard, { card, onupdate, ondelete: vi.fn() });
		await fireEvent.click(screen.getByRole('button', { name: 'Edit' }));
		await fireEvent.click(screen.getByText('Save'));
		await waitFor(() => {
			expect(mockUpdateCard).toHaveBeenCalledWith('card1', {
				title: 'Test Card',
				description: 'A description'
			});
			expect(onupdate).toHaveBeenCalled();
		});
	});

	it('returns to view mode when Cancel is clicked', async () => {
		render(BoardCard, { card, onupdate: vi.fn(), ondelete: vi.fn() });
		await fireEvent.click(screen.getByRole('button', { name: 'Edit' }));
		await fireEvent.click(screen.getByText('Cancel'));
		expect(screen.getByText('Test Card')).toBeInTheDocument();
		expect(screen.queryByDisplayValue('Test Card')).not.toBeInTheDocument();
	});
});
