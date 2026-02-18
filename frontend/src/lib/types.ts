export interface User {
	id: string;
	email: string;
	name: string;
	created_at: string;
}

export interface Board {
	id: string;
	user_id: string;
	name: string;
	created_at: string;
	columns?: Column[];
}

export interface Column {
	id: string;
	board_id: string;
	name: string;
	position: number;
	created_at: string;
	cards: Card[];
}

export interface Card {
	id: string;
	column_id: string;
	title: string;
	description: string;
	position: number;
	created_at: string;
}
