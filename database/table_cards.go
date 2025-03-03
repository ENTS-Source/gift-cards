package database

import (
	"database/sql"
	"errors"
	"time"
)

type Card struct {
	Id        string
	AMemberId *int
	Created   time.Time
}

var cardSelect *sql.Stmt
var cardInsert *sql.Stmt

func (d *Database) prepareCards() error {
	var err error
	if cardSelect, err = d.db.Prepare("SELECT id, amember_id, created FROM cards WHERE id = $1;"); err != nil {
		return err
	}
	if cardInsert, err = d.db.Prepare("INSERT INTO cards (id, amember_id, created) VALUES ($1, NULL, DATETIME()) RETURNING id, amember_id, created;"); err != nil {
		return err
	}
	return nil
}

func (d *Database) GetOrCreateCard(id string) (*Card, error) {
	card := &Card{}
	row := cardSelect.QueryRow(id)
	if err := row.Scan(&card.Id, &card.AMemberId, &card.Created); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			row = cardInsert.QueryRow(id)
			if err = row.Scan(&card.Id, &card.AMemberId, &card.Created); err != nil {
				return nil, err
			}
			return card, nil
		}
		return nil, err
	}
	return card, nil
}
