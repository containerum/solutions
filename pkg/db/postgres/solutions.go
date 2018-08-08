package postgres

import (
	"context"

	"time"

	"git.containerum.net/ch/solutions/pkg/sErrors"
	kube_types "github.com/containerum/kube-client/pkg/model"
	"github.com/json-iterator/go"
)

func (pgdb *pgDB) AddSolution(ctx context.Context, solution kube_types.Solution, userID, templateID, uuid, env string) error {
	pgdb.log.Infoln("Saving solution")

	if _, err := pgdb.eLog.ExecContext(ctx, "INSERT INTO solutions (id, template_id, name, namespace, user_id) "+
		"VALUES ($1, $2, $3, $4, $5)", uuid, templateID, solution.Name, solution.Namespace, userID); err != nil {
		return err
	}

	if _, err := pgdb.eLog.ExecContext(ctx, "INSERT INTO parameters (solution_id, branch, env) "+
		"VALUES ($1, $2, $3)", uuid, solution.Branch, env); err != nil {
		return err
	}
	return nil
}

func (pgdb *pgDB) GetSolutionsList(ctx context.Context, userID string) (*kube_types.SolutionsList, error) {
	pgdb.log.Infoln("Get solutions list")
	var ret kube_types.SolutionsList

	ret.Solutions = make([]kube_types.Solution, 0)

	rows, err := pgdb.qLog.QueryxContext(ctx, "SELECT templates.name, templates.url, solutions.id, solutions.name, solutions.namespace, parameters.env, parameters.branch "+
		"FROM solutions JOIN parameters ON solutions.id = parameters.solution_id JOIN templates ON solutions.template_id = templates.ID WHERE solutions.user_id=$1 AND solutions.is_deleted !='true'", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		solution := kube_types.Solution{}
		var env string
		err := rows.Scan(&solution.Template, &solution.URL, &solution.ID, &solution.Name, &solution.Namespace, &env, &solution.Branch)
		if err != nil {
			return nil, err
		}
		if err := jsoniter.UnmarshalFromString(env, &solution.Env); err != nil {
			return nil, err
		}

		solution.URL = solution.URL + "/tree/" + solution.Branch
		ret.Solutions = append(ret.Solutions, solution)
	}

	return &ret, rows.Err()
}

func (pgdb *pgDB) GetNamespaceSolutionsList(ctx context.Context, namespace string) (*kube_types.SolutionsList, error) {
	pgdb.log.Infoln("Get solutions list")
	var ret kube_types.SolutionsList

	ret.Solutions = make([]kube_types.Solution, 0)

	rows, err := pgdb.qLog.QueryxContext(ctx, "SELECT templates.name, templates.url, solutions.id, solutions.name, solutions.namespace, parameters.env, parameters.branch "+
		"FROM solutions JOIN parameters ON solutions.id = parameters.solution_id JOIN templates ON solutions.template_id = templates.ID WHERE solutions.namespace=$1 AND solutions.is_deleted !='true'", namespace)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		solution := kube_types.Solution{}
		var env string
		err := rows.Scan(&solution.Template, &solution.URL, &solution.ID, &solution.Name, &solution.Namespace, &env, &solution.Branch)
		if err != nil {
			return nil, err
		}
		if err := jsoniter.UnmarshalFromString(env, &solution.Env); err != nil {
			return nil, err
		}

		solution.URL = solution.URL + "/tree/" + solution.Branch
		ret.Solutions = append(ret.Solutions, solution)
	}

	return &ret, rows.Err()
}

func (pgdb *pgDB) GetSolution(ctx context.Context, namespace, solutionName string) (*kube_types.Solution, error) {
	pgdb.log.Infoln("Get solution")

	var solution kube_types.Solution

	rows, err := pgdb.qLog.QueryxContext(ctx, "SELECT templates.name, templates.url, solutions.id, solutions.name, solutions.namespace, parameters.env, parameters.branch "+
		"FROM solutions JOIN parameters ON solutions.id = parameters.solution_id JOIN templates ON solutions.template_id = templates.ID WHERE solutions.name=$1 AND solutions.namespace=$2 AND solutions.is_deleted !='true'", solutionName, namespace)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	if !rows.Next() {
		if rows.Err() != nil {
			return nil, rows.Err()
		}
		return nil, sErrors.ErrSolutionNotExist()
	}

	var env string
	err = rows.Scan(&solution.Template, &solution.URL, &solution.ID, &solution.Name, &solution.Namespace, &env, &solution.Branch)
	if err != nil {
		return nil, err
	}
	if err := jsoniter.UnmarshalFromString(env, &solution.Env); err != nil {
		return nil, err
	}

	solution.URL = solution.URL + "/tree/" + solution.Branch

	return &solution, rows.Err()
}

func (pgdb *pgDB) DeleteSolution(ctx context.Context, namespace, solutionName string) error {
	pgdb.log.Infoln("Deleting solution")

	res, err := pgdb.eLog.ExecContext(ctx, `UPDATE solutions SET is_deleted = 'true', deleted_at=$1 WHERE name=$2 AND namespace=$3 AND is_deleted != 'true'`, time.Now(), solutionName, namespace)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if rows == 0 {
		return sErrors.ErrSolutionNotExist()
	}
	return err
}

func (pgdb *pgDB) CompletelyDeleteSolution(ctx context.Context, namespace, solutionName string) error {
	pgdb.log.Infoln("Deleting solution")

	res, err := pgdb.eLog.ExecContext(ctx, "DELETE FROM solutions WHERE name=$1 AND namespace=$2", solutionName, namespace)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if rows == 0 {
		return sErrors.ErrSolutionNotExist()
	}
	return err
}

func (pgdb *pgDB) CompletelyDeleteSolutions(ctx context.Context, userID string) error {
	pgdb.log.Infoln("Deleting user solutions")

	if _, err := pgdb.eLog.ExecContext(ctx, "DELETE FROM solutions WHERE user_id=$1", userID); err != nil {
		return err
	}
	return nil
}

func (pgdb *pgDB) CompletelyDeleteNamespaceSolutions(ctx context.Context, namespace string) error {
	pgdb.log.Infoln("Deleting namespace solution")

	if _, err := pgdb.eLog.ExecContext(ctx, "DELETE FROM solutions WHERE namespace=$1", namespace); err != nil {
		return err
	}
	return nil
}
