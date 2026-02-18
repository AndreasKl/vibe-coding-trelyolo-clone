import type { User, Board, Card, Column } from './types';

class ApiError extends Error {
	constructor(public status: number, message: string) {
		super(message);
	}
}

async function request<T>(method: string, path: string, body?: unknown): Promise<T> {
	const opts: RequestInit = {
		method,
		headers: { 'Content-Type': 'application/json' },
		credentials: 'include'
	};
	if (body) opts.body = JSON.stringify(body);

	const res = await fetch(path, opts);
	if (!res.ok) {
		const data = await res.json().catch(() => ({ error: res.statusText }));
		throw new ApiError(res.status, data.error || res.statusText);
	}
	if (res.status === 204) return undefined as T;
	return res.json();
}

// Auth
export const signup = (email: string, password: string, name: string) =>
	request<User>('POST', '/api/auth/signup', { email, password, name });

export const login = (email: string, password: string) =>
	request<User>('POST', '/api/auth/login', { email, password });

export const logout = () => request<void>('POST', '/api/auth/logout');

export const me = () => request<User>('GET', '/api/auth/me');

// Boards
export const listBoards = () => request<Board[]>('GET', '/api/boards');

export const createBoard = (name: string) =>
	request<Board>('POST', '/api/boards', { name });

export const getBoard = (id: string) => request<Board>('GET', `/api/boards/${id}`);

export const deleteBoard = (id: string) => request<void>('DELETE', `/api/boards/${id}`);

// Columns
export const createColumn = (boardId: string, name: string) =>
	request<Column>('POST', `/api/boards/${boardId}/columns`, { name });

export const updateColumn = (id: string, data: { name?: string; position?: number }) =>
	request<Column>('PATCH', `/api/columns/${id}`, data);

export const deleteColumn = (id: string) => request<void>('DELETE', `/api/columns/${id}`);

// Cards
export const createCard = (columnId: string, title: string, description = '') =>
	request<Card>('POST', `/api/columns/${columnId}/cards`, { title, description });

export const updateCard = (id: string, data: { title?: string; description?: string }) =>
	request<Card>('PATCH', `/api/cards/${id}`, data);

export const deleteCard = (id: string) => request<void>('DELETE', `/api/cards/${id}`);

export const moveCard = (id: string, columnId: string, position: number) =>
	request<Card>('POST', `/api/cards/${id}/move`, { column_id: columnId, position });

export const moveColumn = (id: string, position: number) =>
	request<Column>('POST', `/api/columns/${id}/move`, { position });
