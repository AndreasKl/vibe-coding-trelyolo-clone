package board

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	DB *sql.DB
}

// Boards

func (s *Store) CreateBoard(ctx context.Context, userID, name string) (*Board, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	b := &Board{}
	err = tx.QueryRowContext(ctx,
		`INSERT INTO boards (user_id, name) VALUES ($1, $2) RETURNING id, user_id, name, created_at`,
		userID, name,
	).Scan(&b.ID, &b.UserID, &b.Name, &b.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Auto-create default columns
	for i, colName := range []string{"Todo", "Doing", "Done"} {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO board_columns (board_id, name, position) VALUES ($1, $2, $3)`,
			b.ID, colName, i,
		)
		if err != nil {
			return nil, err
		}
	}

	return b, tx.Commit()
}

func (s *Store) ListBoards(ctx context.Context, userID string) ([]Board, error) {
	rows, err := s.DB.QueryContext(ctx,
		`SELECT id, user_id, name, created_at FROM boards WHERE user_id=$1 ORDER BY created_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var boards []Board
	for rows.Next() {
		var b Board
		if err := rows.Scan(&b.ID, &b.UserID, &b.Name, &b.CreatedAt); err != nil {
			return nil, err
		}
		boards = append(boards, b)
	}
	return boards, rows.Err()
}

func (s *Store) GetBoard(ctx context.Context, id, userID string) (*Board, error) {
	b := &Board{}
	err := s.DB.QueryRowContext(ctx,
		`SELECT id, user_id, name, created_at FROM boards WHERE id=$1 AND user_id=$2`, id, userID,
	).Scan(&b.ID, &b.UserID, &b.Name, &b.CreatedAt)
	if err != nil {
		return nil, err
	}

	cols, err := s.listColumns(ctx, b.ID)
	if err != nil {
		return nil, err
	}
	b.Columns = cols

	for i := range b.Columns {
		cards, err := s.listCards(ctx, b.Columns[i].ID)
		if err != nil {
			return nil, err
		}
		b.Columns[i].Cards = cards
	}

	return b, nil
}

func (s *Store) DeleteBoard(ctx context.Context, id, userID string) error {
	res, err := s.DB.ExecContext(ctx, `DELETE FROM boards WHERE id=$1 AND user_id=$2`, id, userID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// Columns

func (s *Store) listColumns(ctx context.Context, boardID string) ([]Column, error) {
	rows, err := s.DB.QueryContext(ctx,
		`SELECT id, board_id, name, position, created_at FROM board_columns WHERE board_id=$1 ORDER BY position`, boardID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cols []Column
	for rows.Next() {
		var c Column
		if err := rows.Scan(&c.ID, &c.BoardID, &c.Name, &c.Position, &c.CreatedAt); err != nil {
			return nil, err
		}
		c.Cards = []Card{}
		cols = append(cols, c)
	}
	return cols, rows.Err()
}

func (s *Store) CreateColumn(ctx context.Context, boardID, name string) (*Column, error) {
	c := &Column{}
	err := s.DB.QueryRowContext(ctx,
		`INSERT INTO board_columns (board_id, name, position)
		 VALUES ($1, $2, COALESCE((SELECT MAX(position)+1 FROM board_columns WHERE board_id=$1), 0))
		 RETURNING id, board_id, name, position, created_at`,
		boardID, name,
	).Scan(&c.ID, &c.BoardID, &c.Name, &c.Position, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	c.Cards = []Card{}
	return c, nil
}

func (s *Store) UpdateColumn(ctx context.Context, id string, name *string, position *int) (*Column, error) {
	c := &Column{}
	err := s.DB.QueryRowContext(ctx,
		`UPDATE board_columns SET
			name = COALESCE($2, name),
			position = COALESCE($3, position)
		 WHERE id=$1
		 RETURNING id, board_id, name, position, created_at`,
		id, name, position,
	).Scan(&c.ID, &c.BoardID, &c.Name, &c.Position, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Store) DeleteColumn(ctx context.Context, id string) error {
	_, err := s.DB.ExecContext(ctx, `DELETE FROM board_columns WHERE id=$1`, id)
	return err
}

// ColumnBoardOwner checks column belongs to user's board
func (s *Store) ColumnBoardOwner(ctx context.Context, columnID, userID string) (string, error) {
	var boardID string
	err := s.DB.QueryRowContext(ctx,
		`SELECT bc.board_id FROM board_columns bc
		 JOIN boards b ON b.id = bc.board_id
		 WHERE bc.id=$1 AND b.user_id=$2`,
		columnID, userID,
	).Scan(&boardID)
	return boardID, err
}

// Cards

func (s *Store) listCards(ctx context.Context, columnID string) ([]Card, error) {
	rows, err := s.DB.QueryContext(ctx,
		`SELECT id, column_id, title, description, position, created_at FROM cards WHERE column_id=$1 ORDER BY position`, columnID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cards := []Card{}
	for rows.Next() {
		var c Card
		if err := rows.Scan(&c.ID, &c.ColumnID, &c.Title, &c.Description, &c.Position, &c.CreatedAt); err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}
	return cards, rows.Err()
}

func (s *Store) CreateCard(ctx context.Context, columnID, title, description string) (*Card, error) {
	c := &Card{}
	err := s.DB.QueryRowContext(ctx,
		`INSERT INTO cards (column_id, title, description, position)
		 VALUES ($1, $2, $3, COALESCE((SELECT MAX(position)+1 FROM cards WHERE column_id=$1), 0))
		 RETURNING id, column_id, title, description, position, created_at`,
		columnID, title, description,
	).Scan(&c.ID, &c.ColumnID, &c.Title, &c.Description, &c.Position, &c.CreatedAt)
	return c, err
}

func (s *Store) UpdateCard(ctx context.Context, id string, title, description *string) (*Card, error) {
	c := &Card{}
	err := s.DB.QueryRowContext(ctx,
		`UPDATE cards SET
			title = COALESCE($2, title),
			description = COALESCE($3, description)
		 WHERE id=$1
		 RETURNING id, column_id, title, description, position, created_at`,
		id, title, description,
	).Scan(&c.ID, &c.ColumnID, &c.Title, &c.Description, &c.Position, &c.CreatedAt)
	return c, err
}

func (s *Store) DeleteCard(ctx context.Context, id string) error {
	_, err := s.DB.ExecContext(ctx, `DELETE FROM cards WHERE id=$1`, id)
	return err
}

// CardOwner checks card belongs to user's board
func (s *Store) CardOwner(ctx context.Context, cardID, userID string) error {
	var exists bool
	err := s.DB.QueryRowContext(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM cards c
			JOIN board_columns bc ON bc.id = c.column_id
			JOIN boards b ON b.id = bc.board_id
			WHERE c.id=$1 AND b.user_id=$2
		)`,
		cardID, userID,
	).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("not found")
	}
	return nil
}

func (s *Store) MoveCard(ctx context.Context, cardID, targetColumnID string, targetPosition int) (*Card, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Get current position
	var srcColumnID string
	var srcPosition int
	err = tx.QueryRowContext(ctx,
		`SELECT column_id, position FROM cards WHERE id=$1`, cardID,
	).Scan(&srcColumnID, &srcPosition)
	if err != nil {
		return nil, err
	}

	// Close gap in source column
	_, err = tx.ExecContext(ctx,
		`UPDATE cards SET position = position - 1 WHERE column_id=$1 AND position > $2`,
		srcColumnID, srcPosition,
	)
	if err != nil {
		return nil, err
	}

	// Open gap in target column
	_, err = tx.ExecContext(ctx,
		`UPDATE cards SET position = position + 1 WHERE column_id=$1 AND position >= $2`,
		targetColumnID, targetPosition,
	)
	if err != nil {
		return nil, err
	}

	// Move card
	c := &Card{}
	err = tx.QueryRowContext(ctx,
		`UPDATE cards SET column_id=$2, position=$3 WHERE id=$1
		 RETURNING id, column_id, title, description, position, created_at`,
		cardID, targetColumnID, targetPosition,
	).Scan(&c.ID, &c.ColumnID, &c.Title, &c.Description, &c.Position, &c.CreatedAt)
	if err != nil {
		return nil, err
	}

	return c, tx.Commit()
}
