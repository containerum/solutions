package postgres

import (
	"context"

	"git.containerum.net/ch/solutions/pkg/solerrors"
	kube_types "github.com/containerum/kube-client/pkg/model"
	"github.com/json-iterator/go"
)

func (pgdb *pgDB) CreateTemplate(ctx context.Context, solution kube_types.SolutionTemplate) error {
	pgdb.log.Infoln("Saving solution template")

	images, _ := jsoniter.Marshal(solution.Images)

	if _, err := pgdb.eLog.ExecContext(ctx,
		`INSERT INTO templates (name, cpu, ram, images, url, active) VALUES ($1, $2, $3, $4, $5, $6);`, solution.Name, solution.Limits.CPU, solution.Limits.RAM, string(images), solution.URL, "true"); err != nil {
		return err
	}

	return nil
}

func (pgdb *pgDB) UpdateTemplate(ctx context.Context, solution kube_types.SolutionTemplate) error {
	pgdb.log.Infoln("Updating solution template")

	images, _ := jsoniter.Marshal(solution.Images)

	res, err := pgdb.eLog.ExecContext(ctx,
		`UPDATE templates SET (cpu, ram, images, url) = ($2, $3, $4, $5) 
				WHERE name = $1`, solution.Name, solution.Limits.CPU, solution.Limits.RAM, string(images), solution.URL)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if rows == 0 {
		return solerrors.ErrTemplateNotExist()
	}
	return err
}

func (pgdb *pgDB) ActivateTemplate(ctx context.Context, solution string) error {
	pgdb.log.Infoln("Activating solution template")

	res, err := pgdb.eLog.ExecContext(ctx,
		`UPDATE templates SET active = 'true' 
				WHERE name = $1 AND active = 'false'`, solution)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if rows == 0 {
		return solerrors.ErrTemplateNotExist()
	}
	return err
}

func (pgdb *pgDB) DeactivateTemplate(ctx context.Context, solution string) error {
	pgdb.log.Infoln("Deactivating solution template")

	res, err := pgdb.eLog.ExecContext(ctx,
		`UPDATE templates SET active = 'false' 
				WHERE name = $1 AND active = 'true'`, solution)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if rows == 0 {
		return solerrors.ErrSolutionNotExist()
	}
	return err
}

func (pgdb *pgDB) DeleteTemplate(ctx context.Context, solution string) error {
	pgdb.log.Infoln("deleting solution template")

	res, err := pgdb.eLog.ExecContext(ctx,
		`DELETE FROM templates WHERE name = $1`, solution)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if rows == 0 {
		return solerrors.ErrSolutionNotExist()
	}
	return err
}

func (pgdb *pgDB) GetTemplatesList(ctx context.Context, isAdmin bool) (*kube_types.SolutionsTemplatesList, error) {
	pgdb.log.Infoln("Get solutions templates list")
	var ret kube_types.SolutionsTemplatesList

	query := "SELECT name, id, cpu, ram, images, url, active FROM templates"

	if !isAdmin {
		query = query + " WHERE active = 'true'"
	}

	rows, err := pgdb.qLog.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		solution := kube_types.SolutionTemplate{Limits: &kube_types.SolutionLimits{}}
		var images string
		err := rows.Scan(&solution.Name, &solution.ID, &solution.Limits.CPU, &solution.Limits.RAM, &images, &solution.URL, &solution.Active)
		if err != nil {
			return nil, err
		}
		if err := jsoniter.UnmarshalFromString(images, &solution.Images); err != nil {
			return nil, err
		}

		ret.Solutions = append(ret.Solutions, solution)
	}

	return &ret, rows.Err()
}

func (pgdb *pgDB) GetTemplate(ctx context.Context, name string) (*kube_types.SolutionTemplate, error) {
	pgdb.log.Infoln("Get solution template ", name)
	rows, err := pgdb.qLog.QueryxContext(ctx, "SELECT id, name, cpu, ram, images, url FROM templates WHERE name = $1 AND active = 'true'", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, solerrors.ErrTemplateNotExist()
	}

	solution := kube_types.SolutionTemplate{Limits: &kube_types.SolutionLimits{}}
	var images string
	err = rows.Scan(&solution.ID, &solution.Name, &solution.Limits.CPU, &solution.Limits.RAM, &images, &solution.URL)
	if err != nil {
		return nil, err
	}
	if err = jsoniter.UnmarshalFromString(images, &solution.Images); err != nil {
		return nil, err
	}

	return &solution, err
}
