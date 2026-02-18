import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params, fetch }) => {
	const res = await fetch(`/api/boards/${params.id}`);
	if (!res.ok) {
		return { board: null };
	}
	const board = await res.json();
	return { board };
};
