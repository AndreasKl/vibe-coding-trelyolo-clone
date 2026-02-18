import { render, screen, fireEvent } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import Navbar from './Navbar.svelte';
import type { User } from '$lib/types';

const { mockGoto, mockLogout } = vi.hoisted(() => ({
	mockGoto: vi.fn(),
	mockLogout: vi.fn()
}));

vi.mock('$app/navigation', () => ({ goto: mockGoto }));
vi.mock('$app/paths', () => ({ resolve: (path: string) => path }));
vi.mock('$lib/api', () => ({ logout: mockLogout }));

const user: User = { id: '1', email: 'alice@example.com', name: 'Alice', created_at: '2024-01-01' };

describe('Navbar', () => {
	beforeEach(() => {
		mockGoto.mockReset();
		mockLogout.mockReset().mockResolvedValue(undefined);
	});

	it('renders the logo link', () => {
		render(Navbar, { user: null });
		expect(screen.getByText('Trello Clone')).toBeInTheDocument();
	});

	it('does not show user section when logged out', () => {
		render(Navbar, { user: null });
		expect(screen.queryByText('Logout')).not.toBeInTheDocument();
	});

	it('shows user name and logout button when logged in', () => {
		render(Navbar, { user });
		expect(screen.getByText('Alice')).toBeInTheDocument();
		expect(screen.getByText('Logout')).toBeInTheDocument();
	});

	it('falls back to email when user has no name', () => {
		render(Navbar, { user: { ...user, name: '' } });
		expect(screen.getByText('alice@example.com')).toBeInTheDocument();
	});

	it('calls logout API on logout button click', async () => {
		render(Navbar, { user });
		await fireEvent.click(screen.getByText('Logout'));
		expect(mockLogout).toHaveBeenCalled();
	});
});
