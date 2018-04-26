package postgres

import (
	"context"

	stypes "git.containerum.net/ch/solutions/pkg/models"
	"github.com/json-iterator/go"
)

func (pgdb *pgDB) AddSolution(ctx context.Context, solution stypes.UserSolution, userID string, uuid string, env string) error {
	pgdb.log.Infoln("Saving solutions list")

	_, err := pgdb.qLog.QueryxContext(ctx, "INSERT INTO solutions (id, template, name, namespace, user_id) "+
		"VALUES ($1, $2, $3, $4, $5)", uuid, solution.Template, solution.Name, solution.Namespace, userID)
	if err != nil {
		return err
	}

	_, err = pgdb.qLog.QueryxContext(ctx, "INSERT INTO parameters (solution_id, branch, env) "+
		"VALUES ($1, $2, $3)", uuid, solution.Branch, env)
	if err != nil {
		return err
	}
	return err
}

func (pgdb *pgDB) AddDeployment(ctx context.Context, name string, solutionID string) error {
	pgdb.log.Infoln("Adding deployment")

	_, err := pgdb.qLog.QueryxContext(ctx, "INSERT INTO deployments (deploy_name, solution_id) "+
		"VALUES ($1, $2)", name, solutionID)
	if err != nil {
		return err
	}
	return err
}

func (pgdb *pgDB) AddService(ctx context.Context, name string, solutionID string) error {
	pgdb.log.Infoln("Adding service")

	_, err := pgdb.qLog.QueryxContext(ctx, "INSERT INTO services (service_name, solution_id) "+
		"VALUES ($1, $2)", name, solutionID)
	if err != nil {
		return err
	}
	return err
}

func (pgdb *pgDB) GetUserSolutionsList(ctx context.Context, userID string) (*stypes.UserSolutionsList, error) {
	pgdb.log.Infoln("Get solutions list")
	var ret stypes.UserSolutionsList

	rows, err := pgdb.qLog.QueryxContext(ctx, "SELECT solutions.template, solutions.name, solutions.namespace, parameters.env, parameters.branch "+
		"FROM solutions JOIN parameters ON solutions.id = parameters.solution_id WHERE solutions.user_id=$1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		solution := stypes.UserSolution{}
		var env string
		err := rows.Scan(&solution.Template, &solution.Name, &solution.Namespace, &env, &solution.Branch)
		if err != nil {
			return nil, err
		}
		if err := jsoniter.UnmarshalFromString(env, &solution.Env); err != nil {
			return nil, err
		}

		ret.Solutions = append(ret.Solutions, solution)
	}

	return &ret, rows.Err()
}

func (pgdb *pgDB) GetUserSolution(ctx context.Context, solutionName string) (*stypes.UserSolution, error) {
	pgdb.log.Infoln("Get solutions list")

	rows, err := pgdb.qLog.QueryxContext(ctx, "SELECT solutions.template, solutions.name, solutions.namespace, parameters.env, parameters.branch "+
		"FROM solutions JOIN parameters ON solutions.id = parameters.solution_id WHERE solution.name=$1", solutionName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, rows.Err()
	}
	solution := stypes.UserSolution{}
	var env string
	err = rows.Scan(&solution.Template, &solution.Name, &solution.Namespace, &env, &solution.Branch)
	if err != nil {
		return nil, err
	}
	if err := jsoniter.UnmarshalFromString(env, &solution.Env); err != nil {
		return nil, err
	}
	return &solution, rows.Err()
}

func (pgdb *pgDB) DeleteSolution(ctx context.Context, name string) error {
	pgdb.log.Infoln("Deleting solution")

	res, err := pgdb.eLog.ExecContext(ctx, "DELETE FROM solutions WHERE name=$1", name)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if rows == 0 {
		return solutionsErrorsErrSolutionNotExist()
	}
	return nil
}

func (pgdb *pgDB) GetUserSolutionsDeployments(ctx context.Context, solutionName string) (deployments []string, ns *string, err error) {
	pgdb.log.Infoln("Get solution deployments")

	rows, err := pgdb.qLog.QueryxContext(ctx, "SELECT solutions.namespace, deployments.deploy_name "+
		"FROM solutions JOIN deployments ON solutions.id = deployments.solution_id WHERE solutions.name=$1", solutionName)
	if err != nil {
		return nil, nil, err
	}
	deployments = make([]string, 0)

	defer rows.Close()
	for rows.Next() {
		var deploy string
		err := rows.Scan(&ns, &deploy)
		if err != nil {
			return nil, nil, err
		}

		deployments = append(deployments, deploy)
	}
	return deployments, ns, rows.Err()
}

func (pgdb *pgDB) GetUserSolutionsServices(ctx context.Context, solutionName string) (services []string, ns *string, err error) {
	pgdb.log.Infoln("Get solution services")

	rows, err := pgdb.qLog.QueryxContext(ctx, "SELECT solutions.namespace, services.service_name "+
		"FROM solutions JOIN services ON solutions.id = services.solution_id WHERE solutions.name=$1", solutionName)
	if err != nil {
		return nil, nil, err
	}
	services = make([]string, 0)

	defer rows.Close()
	for rows.Next() {
		var deploy string
		err := rows.Scan(&ns, &deploy)
		if err != nil {
			return nil, nil, err
		}

		services = append(services, deploy)
	}
	return services, ns, rows.Err()
}
