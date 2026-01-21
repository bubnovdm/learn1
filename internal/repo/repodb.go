package repo

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

/*
GetItems(userID int) []*repo.Item
AddItem(userID int, items []*repo.Item)
RemoveItem(userID int, skuID int)
ClearCart(userID int)
*/

type DBrepo struct {
	db *sqlx.DB
}

func NewPostgresRepo(db *sqlx.DB) *DBrepo {
	return &DBrepo{db: db}
}

func (d *DBrepo) GetItems(userID int) ([]*Item, error) {
	query := sq.Select("sku_id", "count").From("cart_table").Where(sq.Eq{"user_id": userID}).PlaceholderFormat(sq.Dollar)
	items, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	var result []*Item
	if err := d.db.Select(&result, items, args...); err != nil {
		return nil, err
	}
	return result, nil
}

// TODO:impliment
func (d *DBrepo) AddItem(userID int, items []*Item) error {
	return nil
}

func (d *DBrepo) RemoveItem(userID int, skuID int) error {
	return nil
}
func (d *DBrepo) ClearCart(userID int) error {
	return nil
}
