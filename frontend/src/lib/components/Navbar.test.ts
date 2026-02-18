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
		expect(screen.getByText('FlowBoard')).toBeInTheDocument();
	});

	it('does not show user section when logged out', () => {
		render(Navbar, { user: null });
		expect(screen.queryByText('Sign out')).not.toBeInTheDocument();
	});

	it('shows avatar and sign out button when logged in', () => {
		render(Navbar, { user });
		expect(screen.getByTitle('Alice')).toBeInTheDocument();
		expect(screen.getByText('Sign out')).toBeInTheDocument();
	});

	it('falls back to email when user has no name', () => {
		render(Navbar, { user: { ...user, name: '' } });
		expect(screen.getByTitle('alice@example.com')).toBeInTheDocument();
	});

	it('calls logout API on sign out button click', async () => {
		render(Navbar, { user });
		await fireEvent.click(screen.getByText('Sign out'));
		expect(mockLogout).toHaveBeenCalled();
	});
});
