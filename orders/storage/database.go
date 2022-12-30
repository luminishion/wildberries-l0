package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

type Database struct {
	db                  *sqlx.DB
	stmtOrderSelect     *sqlx.Stmt
	stmtItemSelect      *sqlx.Stmt
	stmtOrderInsert     *sqlx.NamedStmt
	stmtItemInsert      *sqlx.NamedStmt
	stmtOrderListSelect *sqlx.Stmt
}

func NewDatabase(databaseURL string) (*Database, error) {
	db, err := sqlx.Connect("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("database new connect: %w", err)
	}

	if _, err := db.Exec(sqlDbSchema); err != nil {
		return nil, fmt.Errorf("database new schema: %w", err)
	}

	stmtOrderSelect, err := db.Preparex(sqlOrderSelect)
	if err != nil {
		return nil, fmt.Errorf("database new prepare 0: %w", err)
	}

	stmtItemSelect, err := db.Preparex(sqlItemSelect)
	if err != nil {
		return nil, fmt.Errorf("database new prepare 1: %w", err)
	}

	stmtOrderInsert, err := db.PrepareNamed(sqlOrderInsert)
	if err != nil {
		return nil, fmt.Errorf("database new prepare 2: %w", err)
	}

	stmtItemInsert, err := db.PrepareNamed(sqlItemInsert)
	if err != nil {
		return nil, fmt.Errorf("database new prepare 3: %w", err)
	}

	stmtOrderListSelect, err := db.Preparex(sqlOrderListSelect)
	if err != nil {
		return nil, fmt.Errorf("database new prepare 4: %w", err)
	}

	d := &Database{
		db:                  db,
		stmtOrderSelect:     stmtOrderSelect,
		stmtItemSelect:      stmtItemSelect,
		stmtOrderInsert:     stmtOrderInsert,
		stmtItemInsert:      stmtItemInsert,
		stmtOrderListSelect: stmtOrderListSelect,
	}

	return d, nil
}

func (d *Database) LoadList() ([]string, error) {
	rows, err := d.stmtOrderListSelect.Queryx()
	if err != nil {
		return nil, fmt.Errorf("database list select: %w", err)
	}
	defer rows.Close()

	var orders []string

	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("database list row: %w", err)
		}

		orders = append(orders, id)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("database list end: %w", err)
	}

	return orders, nil
}

func (d *Database) Store(order Order) error {
	tx, err := d.db.Beginx()
	if err != nil {
		return fmt.Errorf("database store begin: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.NamedStmt(d.stmtOrderInsert).Exec(&order); err != nil {
		return fmt.Errorf("database store order: %w", err)
	}

	stmtIns := tx.NamedStmt(d.stmtItemInsert)
	orderId := order.Id

	for _, item := range order.Items {
		item.OrderId = orderId

		if _, err := stmtIns.Exec(&item); err != nil {
			return fmt.Errorf("database store item: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("database store commit: %w", err)
	}

	return nil
}

func (d *Database) Load(id string) (*Order, error) {
	row := d.stmtOrderSelect.QueryRowx(id)

	var order Order
	err := row.StructScan(&order)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("database load order: %w", err)
	}

	rows, err := d.stmtItemSelect.Queryx(id)
	if err != nil {
		return nil, fmt.Errorf("database load items: %w", err)
	}
	defer rows.Close()

	var items []Item

	for rows.Next() {
		var item Item
		err := rows.StructScan(&item)
		if err != nil {
			return nil, fmt.Errorf("database load row: %w", err)
		}

		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("database load end: %w", err)
	}

	order.Items = items

	return &order, nil
}

func (d *Database) Close() {
	d.db.Close()
}
