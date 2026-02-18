import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ fetch, url }) => {
	try {
		const res = await fetch('/api/auth/me');
		if (res.ok) {
			const user = await res.json();
			return { user };
		}
	} catch {
		// not authenticated
	}

	const publicPaths = ['/login', '/signup'];
	if (!publicPaths.includes(url.pathname)) {
		return { user: null, needsAuth: true };
	}

	return { user: null };
};
