import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import * as api from './api';

describe('api', () => {
	const mockFetch = vi.fn();

	beforeEach(() => {
		vi.stubGlobal('fetch', mockFetch);
	});

	afterEach(() => {
		vi.unstubAllGlobals();
		vi.clearAllMocks();
	});

	function ok(data: unknown, status = 200) {
		return { ok: true, status, json: async () => data } as Response;
	}

	function noContent() {
		return { ok: true, status: 204, json: async () => undefined } as unknown as Response;
	}

	const user = { id: '1', email: 'test@example.com', name: 'Test User', created_at: '2024-01-01' };
	const board = { id: 'b1', user_id: '1', name: 'My Board', created_at: '2024-01-01', columns: [] };
	const column = { id: 'col1', board_id: 'b1', name: 'To Do', position: 0, created_at: '2024-01-01', cards: [] };
	const card = { id: 'card1', column_id: 'col1', title: 'My Card', description: '', position: 0, created_at: '2024-01-01' };

	describe('auth', () => {
		it('signup sends credentials and returns user', async () => {
			mockFetch.mockResolvedValueOnce(ok(user));
			const result = await api.signup('test@example.com', 'password123', 'Test User');
			expect(mockFetch).toHaveBeenCalledWith(
				'/api/auth/signup',
				expect.objectContaining({
					method: 'POST',
					body: JSON.stringify({ email: 'test@example.com', password: 'password123', name: 'Test User' })
				})
			);
			expect(result).toEqual(user);
		});

		it('login sends credentials and returns user', async () => {
			mockFetch.mockResolvedValueOnce(ok(user));
			const result = await api.login('test@example.com', 'password123');
			expect(mockFetch).toHaveBeenCalledWith(
				'/api/auth/login',
				expect.objectContaining({
					method: 'POST',
					body: JSON.stringify({ email: 'test@example.com', password: 'password123' })
				})
			);
			expect(result).toEqual(user);
		});

		it('logout sends POST and returns undefined on 204', async () => {
			mockFetch.mockResolvedValueOnce(noContent());
			const result = await api.logout();
			expect(mockFetch).toHaveBeenCalledWith(
				'/api/auth/logout',
				expect.objectContaining({ method: 'POST' })
			);
			expect(result).toBeUndefined();
		});

		it('me returns the current user', async () => {
			mockFetch.mockResolvedValueOnce(ok(user));
			const result = await api.me();
			expect(mockFetch).toHaveBeenCalledWith(
				'/api/auth/me',
				expect.objectContaining({ method: 'GET' })
			);
			expect(result).toEqual(user);
		});
	});

	describe('boards', () => {
		it('listBoards returns array of boards', async () => {
			mockFetch.mockResolvedValueOnce(ok([board]));
			const result = await api.listBoards();
			expect(mockFetch).toHaveBeenCalledWith('/api/boards', expect.objectContaining({ method: 'GET' }));
			expect(result).toEqual([board]);
		});

		it('createBoard sends name and returns board', async () => {
			mockFetch.mockResolvedValueOnce(ok(board));
			const result = await api.createBoard('My Board');
			expect(mockFetch).toHaveBeenCalledWith(
				'/api/boards',
				expect.objectContaining({
					method: 'POST',
					body: JSON.stringify({ name: 'My Board' })
				})
			);
			expect(result).toEqual(board);
		});

		it('getBoard fetches board by id', async () => {
			mockFetch.mockResolvedValueOnce(ok(board));
			const result = await api.getBoard('b1');
			expect(mockFetch).toHaveBeenCalledWith(
				'/api/boards/b1',
				expect.objectContaining({ method: 'GET' })
			);
			expect(result).toEqual(board);
		});

		it('deleteBoard sends DELETE request', async () => {
			mockFetch.mockResolvedValueOnce(noContent());
			await api.deleteBoard('b1');
			expect(mockFetch).toHaveBeenCalledWith(
				'/api/boards/b1',
				expect.objectContaining({ method: 'DELETE' })
			);
		});
	});

	describe('columns', () => {
		it('createColumn sends name and returns column', async () => {
			mockFetch.mockResolvedValueOnce(ok(column));
			const result = await api.createColumn('b1', 'To Do');
			expect(mockFetch).toHaveBeenCalledWith(
				'/api/boards/b1/columns',
				expect.objectContaining({
					method: 'POST',
					body: JSON.stringify({ name: 'To Do' })
				})
			);
			expect(result).toEqual(column);
		});

		it('updateColumn sends patch data', async () => {
			mockFetch.mockResolvedValueOnce(ok({ ...column, name: 'In Progress' }));
			const result = await api.updateColumn('col1', { name: 'In Progress' });
			expect(mockFetch).toHaveBeenCalledWith(
				'/api/columns/col1',
				expect.objectContaining({
					method: 'PATCH',
					body: JSON.stringify({ name: 'In Progress' })
				})
			);
			expect(result.name).toBe('In Progress');
		});

		it('deleteColumn sends DELETE request', async () => {
			mockFetch.mockResolvedValueOnce(noContent());
			await api.deleteColumn('col1');
			expect(mockFetch).toHaveBeenCalledWith(
				'/api/columns/col1',
				expect.objectContaining({ method: 'DELETE' })
			);
		});
	});

	describe('cards', () => {
		it('createCard sends title and empty description by default', async () => {
			mockFetch.mockResolvedValueOnce(ok(card));
			const result = await api.createCard('col1', 'My Card');
			expect(mockFetch).toHaveBeenCalledWith(
				'/api/columns/col1/cards',
				expect.objectContaining({
					method: 'POST',
					body: JSON.stringify({ title: 'My Card', description: '' })
				})
			);
			expect(result).toEqual(card);
		});

		it('createCard sends custom description when provided', async () => {
			mockFetch.mockResolvedValueOnce(ok({ ...card, description: 'Details' }));
			await api.createCard('col1', 'My Card', 'Details');
			expect(mockFetch).toHaveBeenCalledWith(
				'/api/columns/col1/cards',
				expect.objectContaining({
					body: JSON.stringify({ title: 'My Card', description: 'Details' })
				})
			);
		});

		it('updateCard sends patch fields', async () => {
			mockFetch.mockResolvedValueOnce(ok({ ...card, title: 'Updated' }));
			const result = await api.updateCard('card1', { title: 'Updated' });
			expect(mockFetch).toHaveBeenCalledWith(
				'/api/cards/card1',
				expect.objectContaining({
					method: 'PATCH',
					body: JSON.stringify({ title: 'Updated' })
				})
			);
			expect(result.title).toBe('Updated');
		});

		it('deleteCard sends DELETE request', async () => {
			mockFetch.mockResolvedValueOnce(noContent());
			await api.deleteCard('card1');
			expect(mockFetch).toHaveBeenCalledWith(
				'/api/cards/card1',
				expect.objectContaining({ method: 'DELETE' })
			);
		});

		it('moveCard sends column_id and position', async () => {
			mockFetch.mockResolvedValueOnce(ok(card));
			await api.moveCard('card1', 'col2', 3);
			expect(mockFetch).toHaveBeenCalledWith(
				'/api/cards/card1/move',
				expect.objectContaining({
					method: 'POST',
					body: JSON.stringify({ column_id: 'col2', position: 3 })
				})
			);
		});
	});

	describe('error handling', () => {
		it('throws ApiError with status and message on non-ok response', async () => {
			mockFetch.mockResolvedValueOnce({
				ok: false,
				status: 401,
				statusText: 'Unauthorized',
				json: async () => ({ error: 'not authenticated' })
			} as Response);
			await expect(api.me()).rejects.toThrow('not authenticated');
		});
	});
});
