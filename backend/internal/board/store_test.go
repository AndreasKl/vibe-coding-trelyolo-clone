package board

import (
	"context"
	"database/sql"
	"testing"

	"trello-clone/internal/auth"
	"trello-clone/internal/testutil"
)

func createUser(t *testing.T, db *sql.DB, email string) *auth.User {
	t.Helper()
	s := &auth.Store{DB: db}
	hash, _ := auth.HashPassword("password123")
	u, err := s.CreateUser(context.Background(), email, hash, "Test User")
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	return u
}

func TestCreateBoard(t *testing.T) {
	db := testutil.SetupDB(t)
	s := &Store{DB: db}
	u := createUser(t, db, "board@example.com")

	b, err := s.CreateBoard(context.Background(), u.ID, "My Board")
	if err != nil {
		t.Fatalf("create board: %v", err)
	}
	if b.ID == "" {
		t.Fatal("expected non-empty board ID")
	}
	if b.Name != "My Board" {
		t.Fatalf("name = %q, want My Board", b.Name)
	}

	// Verify 3 default columns were created
	full, err := s.GetBoard(context.Background(), b.ID, u.ID)
	if err != nil {
		t.Fatalf("get board: %v", err)
	}
	if len(full.Columns) != 3 {
		t.Fatalf("columns = %d, want 3", len(full.Columns))
	}
	names := []string{full.Columns[0].Name, full.Columns[1].Name, full.Columns[2].Name}
	if names[0] != "Todo" || names[1] != "Doing" || names[2] != "Done" {
		t.Fatalf("column names = %v, want [Todo Doing Done]", names)
	}
}

func TestListBoards(t *testing.T) {
	db := testutil.SetupDB(t)
	s := &Store{DB: db}
	u := createUser(t, db, "list@example.com")
	ctx := context.Background()

	s.CreateBoard(ctx, u.ID, "Board 1")
	s.CreateBoard(ctx, u.ID, "Board 2")

	boards, err := s.ListBoards(ctx, u.ID)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(boards) != 2 {
		t.Fatalf("count = %d, want 2", len(boards))
	}
}

func TestListBoardsEmpty(t *testing.T) {
	db := testutil.SetupDB(t)
	s := &Store{DB: db}
	u := createUser(t, db, "empty@example.com")

	boards, err := s.ListBoards(context.Background(), u.ID)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if boards != nil && len(boards) != 0 {
		t.Fatalf("expected nil or empty, got %d boards", len(boards))
	}
}

func TestGetBoard(t *testing.T) {
	db := testutil.SetupDB(t)
	s := &Store{DB: db}
	u := createUser(t, db, "get@example.com")
	ctx := context.Background()

	b, _ := s.CreateBoard(ctx, u.ID, "Get Board")
	full, _ := s.GetBoard(ctx, b.ID, u.ID)

	// Add a card to the first column
	col := full.Columns[0]
	card, err := s.CreateCard(ctx, col.ID, "Card 1", "Description")
	if err != nil {
		t.Fatalf("create card: %v", err)
	}

	// Get board again with nested data
	full2, err := s.GetBoard(ctx, b.ID, u.ID)
	if err != nil {
		t.Fatalf("get board: %v", err)
	}
	if len(full2.Columns) != 3 {
		t.Fatalf("columns = %d, want 3", len(full2.Columns))
	}
	if len(full2.Columns[0].Cards) != 1 {
		t.Fatalf("cards in first column = %d, want 1", len(full2.Columns[0].Cards))
	}
	if full2.Columns[0].Cards[0].ID != card.ID {
		t.Fatal("card ID mismatch")
	}
}

func TestGetBoardWrongUser(t *testing.T) {
	db := testutil.SetupDB(t)
	s := &Store{DB: db}
	u1 := createUser(t, db, "owner@example.com")
	u2 := createUser(t, db, "other@example.com")
	ctx := context.Background()

	b, _ := s.CreateBoard(ctx, u1.ID, "Private Board")
	_, err := s.GetBoard(ctx, b.ID, u2.ID)
	if err == nil {
		t.Fatal("expected error for wrong user")
	}
}

func TestDeleteBoard(t *testing.T) {
	db := testutil.SetupDB(t)
	s := &Store{DB: db}
	u := createUser(t, db, "delete@example.com")
	ctx := context.Background()

	b, _ := s.CreateBoard(ctx, u.ID, "Delete Board")
	if err := s.DeleteBoard(ctx, b.ID, u.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	_, err := s.GetBoard(ctx, b.ID, u.ID)
	if err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestDeleteBoardWrongUser(t *testing.T) {
	db := testutil.SetupDB(t)
	s := &Store{DB: db}
	u1 := createUser(t, db, "delowner@example.com")
	u2 := createUser(t, db, "delother@example.com")
	ctx := context.Background()

	b, _ := s.CreateBoard(ctx, u1.ID, "Board")
	err := s.DeleteBoard(ctx, b.ID, u2.ID)
	if err != sql.ErrNoRows {
		t.Fatalf("err = %v, want sql.ErrNoRows", err)
	}
}

func TestCreateColumn(t *testing.T) {
	db := testutil.SetupDB(t)
	s := &Store{DB: db}
	u := createUser(t, db, "col@example.com")
	ctx := context.Background()

	b, _ := s.CreateBoard(ctx, u.ID, "Board")
	col, err := s.CreateColumn(ctx, b.ID, "Extra")
	if err != nil {
		t.Fatalf("create column: %v", err)
	}
	if col.Name != "Extra" {
		t.Fatalf("name = %q, want Extra", col.Name)
	}
	// Default board has 3 columns (pos 0,1,2), so new one should be pos 3
	if col.Position != 3 {
		t.Fatalf("position = %d, want 3", col.Position)
	}
}

func TestUpdateColumn(t *testing.T) {
	db := testutil.SetupDB(t)
	s := &Store{DB: db}
	u := createUser(t, db, "upd@example.com")
	ctx := context.Background()

	b, _ := s.CreateBoard(ctx, u.ID, "Board")
	full, _ := s.GetBoard(ctx, b.ID, u.ID)
	col := full.Columns[0]

	newName := "Renamed"
	updated, err := s.UpdateColumn(ctx, col.ID, &newName, nil)
	if err != nil {
		t.Fatalf("update column: %v", err)
	}
	if updated.Name != "Renamed" {
		t.Fatalf("name = %q, want Renamed", updated.Name)
	}
}

func TestDeleteColumn(t *testing.T) {
	db := testutil.SetupDB(t)
	s := &Store{DB: db}
	u := createUser(t, db, "delcol@example.com")
	ctx := context.Background()

	b, _ := s.CreateBoard(ctx, u.ID, "Board")
	full, _ := s.GetBoard(ctx, b.ID, u.ID)
	col := full.Columns[0]

	// Add a card so we can verify cascade delete
	s.CreateCard(ctx, col.ID, "Card in deleted column", "")

	if err := s.DeleteColumn(ctx, col.ID); err != nil {
		t.Fatalf("delete column: %v", err)
	}

	// Verify column and its cards are gone
	full2, _ := s.GetBoard(ctx, b.ID, u.ID)
	if len(full2.Columns) != 2 {
		t.Fatalf("columns = %d, want 2", len(full2.Columns))
	}
}

func TestColumnBoardOwner(t *testing.T) {
	db := testutil.SetupDB(t)
	s := &Store{DB: db}
	u := createUser(t, db, "colown@example.com")
	ctx := context.Background()

	b, _ := s.CreateBoard(ctx, u.ID, "Board")
	full, _ := s.GetBoard(ctx, b.ID, u.ID)
	col := full.Columns[0]

	boardID, err := s.ColumnBoardOwner(ctx, col.ID, u.ID)
	if err != nil {
		t.Fatalf("column board owner: %v", err)
	}
	if boardID != b.ID {
		t.Fatalf("boardID = %s, want %s", boardID, b.ID)
	}
}

func TestColumnBoardOwnerWrongUser(t *testing.T) {
	db := testutil.SetupDB(t)
	s := &Store{DB: db}
	u1 := createUser(t, db, "colown1@example.com")
	u2 := createUser(t, db, "colown2@example.com")
	ctx := context.Background()

	b, _ := s.CreateBoard(ctx, u1.ID, "Board")
	full, _ := s.GetBoard(ctx, b.ID, u1.ID)
	col := full.Columns[0]

	_, err := s.ColumnBoardOwner(ctx, col.ID, u2.ID)
	if err == nil {
		t.Fatal("expected error for wrong user")
	}
}

func TestCreateCard(t *testing.T) {
	db := testutil.SetupDB(t)
	s := &Store{DB: db}
	u := createUser(t, db, "card@example.com")
	ctx := context.Background()

	b, _ := s.CreateBoard(ctx, u.ID, "Board")
	full, _ := s.GetBoard(ctx, b.ID, u.ID)
	col := full.Columns[0]

	card, err := s.CreateCard(ctx, col.ID, "My Card", "Some description")
	if err != nil {
		t.Fatalf("create card: %v", err)
	}
	if card.Title != "My Card" {
		t.Fatalf("title = %q, want My Card", card.Title)
	}
	if card.Position != 0 {
		t.Fatalf("position = %d, want 0", card.Position)
	}

	// Second card should auto-position at 1
	card2, _ := s.CreateCard(ctx, col.ID, "Card 2", "")
	if card2.Position != 1 {
		t.Fatalf("position = %d, want 1", card2.Position)
	}
}

func TestUpdateCard(t *testing.T) {
	db := testutil.SetupDB(t)
	s := &Store{DB: db}
	u := createUser(t, db, "updcard@example.com")
	ctx := context.Background()

	b, _ := s.CreateBoard(ctx, u.ID, "Board")
	full, _ := s.GetBoard(ctx, b.ID, u.ID)
	col := full.Columns[0]

	card, _ := s.CreateCard(ctx, col.ID, "Original", "Orig desc")
	newTitle := "Updated"
	newDesc := "New desc"
	updated, err := s.UpdateCard(ctx, card.ID, &newTitle, &newDesc)
	if err != nil {
		t.Fatalf("update card: %v", err)
	}
	if updated.Title != "Updated" {
		t.Fatalf("title = %q, want Updated", updated.Title)
	}
	if updated.Description != "New desc" {
		t.Fatalf("description = %q, want New desc", updated.Description)
	}
}

func TestDeleteCard(t *testing.T) {
	db := testutil.SetupDB(t)
	s := &Store{DB: db}
	u := createUser(t, db, "delcard@example.com")
	ctx := context.Background()

	b, _ := s.CreateBoard(ctx, u.ID, "Board")
	full, _ := s.GetBoard(ctx, b.ID, u.ID)
	col := full.Columns[0]

	card, _ := s.CreateCard(ctx, col.ID, "To Delete", "")
	if err := s.DeleteCard(ctx, card.ID); err != nil {
		t.Fatalf("delete card: %v", err)
	}

	// Verify card is gone
	full2, _ := s.GetBoard(ctx, b.ID, u.ID)
	if len(full2.Columns[0].Cards) != 0 {
		t.Fatalf("cards = %d, want 0", len(full2.Columns[0].Cards))
	}
}

func TestCardOwner(t *testing.T) {
	db := testutil.SetupDB(t)
	s := &Store{DB: db}
	u := createUser(t, db, "cardown@example.com")
	ctx := context.Background()

	b, _ := s.CreateBoard(ctx, u.ID, "Board")
	full, _ := s.GetBoard(ctx, b.ID, u.ID)
	col := full.Columns[0]
	card, _ := s.CreateCard(ctx, col.ID, "Card", "")

	if err := s.CardOwner(ctx, card.ID, u.ID); err != nil {
		t.Fatalf("card owner: %v", err)
	}
}

func TestCardOwnerWrongUser(t *testing.T) {
	db := testutil.SetupDB(t)
	s := &Store{DB: db}
	u1 := createUser(t, db, "cardown1@example.com")
	u2 := createUser(t, db, "cardown2@example.com")
	ctx := context.Background()

	b, _ := s.CreateBoard(ctx, u1.ID, "Board")
	full, _ := s.GetBoard(ctx, b.ID, u1.ID)
	col := full.Columns[0]
	card, _ := s.CreateCard(ctx, col.ID, "Card", "")

	if err := s.CardOwner(ctx, card.ID, u2.ID); err == nil {
		t.Fatal("expected error for wrong user")
	}
}

func TestMoveCardSameColumn(t *testing.T) {
	db := testutil.SetupDB(t)
	s := &Store{DB: db}
	u := createUser(t, db, "movesame@example.com")
	ctx := context.Background()

	b, _ := s.CreateBoard(ctx, u.ID, "Board")
	full, _ := s.GetBoard(ctx, b.ID, u.ID)
	col := full.Columns[0]

	// Create 3 cards: pos 0, 1, 2
	c0, _ := s.CreateCard(ctx, col.ID, "Card 0", "")
	c1, _ := s.CreateCard(ctx, col.ID, "Card 1", "")
	c2, _ := s.CreateCard(ctx, col.ID, "Card 2", "")

	// Move card 0 to position 2
	moved, err := s.MoveCard(ctx, c0.ID, col.ID, 2)
	if err != nil {
		t.Fatalf("move card: %v", err)
	}
	if moved.Position != 2 {
		t.Fatalf("moved position = %d, want 2", moved.Position)
	}

	// Verify positions: c1 should be 0, c2 should be 2 (was 1 after gap close, then +1 from gap open), c0 should be 2
	full2, _ := s.GetBoard(ctx, b.ID, u.ID)
	cards := full2.Columns[0].Cards
	positions := make(map[string]int)
	for _, c := range cards {
		positions[c.ID] = c.Position
	}
	// After move: gap close shifts c1 to 0, c2 to 1. Then gap open at 2 shifts nothing. c0 placed at 2.
	if positions[c1.ID] != 0 {
		t.Fatalf("c1 position = %d, want 0", positions[c1.ID])
	}
	if positions[c2.ID] != 1 {
		t.Fatalf("c2 position = %d, want 1", positions[c2.ID])
	}
	if positions[c0.ID] != 2 {
		t.Fatalf("c0 position = %d, want 2", positions[c0.ID])
	}
}

func TestMoveCardAcrossColumns(t *testing.T) {
	db := testutil.SetupDB(t)
	s := &Store{DB: db}
	u := createUser(t, db, "moveacross@example.com")
	ctx := context.Background()

	b, _ := s.CreateBoard(ctx, u.ID, "Board")
	full, _ := s.GetBoard(ctx, b.ID, u.ID)
	srcCol := full.Columns[0]
	dstCol := full.Columns[1]

	// Create cards in source column
	c0, _ := s.CreateCard(ctx, srcCol.ID, "Src Card 0", "")
	c1, _ := s.CreateCard(ctx, srcCol.ID, "Src Card 1", "")

	// Create a card in destination column
	d0, _ := s.CreateCard(ctx, dstCol.ID, "Dst Card 0", "")

	// Move c0 from source to dest at position 0
	moved, err := s.MoveCard(ctx, c0.ID, dstCol.ID, 0)
	if err != nil {
		t.Fatalf("move card: %v", err)
	}
	if moved.ColumnID != dstCol.ID {
		t.Fatalf("column_id = %s, want %s", moved.ColumnID, dstCol.ID)
	}
	if moved.Position != 0 {
		t.Fatalf("position = %d, want 0", moved.Position)
	}

	// Source column: c1 should be at position 0 (gap closed)
	full2, _ := s.GetBoard(ctx, b.ID, u.ID)
	srcCards := full2.Columns[0].Cards
	if len(srcCards) != 1 {
		t.Fatalf("src cards = %d, want 1", len(srcCards))
	}
	if srcCards[0].ID != c1.ID || srcCards[0].Position != 0 {
		t.Fatalf("src card: id=%s pos=%d, want id=%s pos=0", srcCards[0].ID, srcCards[0].Position, c1.ID)
	}

	// Dest column: c0 at 0, d0 at 1
	dstCards := full2.Columns[1].Cards
	if len(dstCards) != 2 {
		t.Fatalf("dst cards = %d, want 2", len(dstCards))
	}
	dstPositions := make(map[string]int)
	for _, c := range dstCards {
		dstPositions[c.ID] = c.Position
	}
	if dstPositions[c0.ID] != 0 {
		t.Fatalf("c0 position = %d, want 0", dstPositions[c0.ID])
	}
	if dstPositions[d0.ID] != 1 {
		t.Fatalf("d0 position = %d, want 1", dstPositions[d0.ID])
	}
}
