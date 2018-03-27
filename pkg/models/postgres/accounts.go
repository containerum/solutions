package postgres

import (
	"context"
)

func (db *pgDB) DoSMTH(ctx context.Context) error {
	db.log.Infoln("DO SMTH")

	return nil
}
