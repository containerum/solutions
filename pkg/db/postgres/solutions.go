package postgres

import (
	"context"

	"time"

	"git.containerum.net/ch/solutions/pkg/sErrors"
	stypes "github.com/containerum/kube-client/pkg/model"
	"github.com/json-iterator/go"
)

func (pgdb *pgDB) AddSolution(ctx context.Context, solution stypes.UserSolution, userID, templateID, uuid, env string) error {
	pgdb.log.Infoln("Saving solution")

	_, err := pgdb.qLog.QueryxContext(ctx, "INSERT INTO solutions (id, template_id, name, namespace, user_id) "+
		"VALUES ($1, $2, $3, $4, $5)", uuid, templateID, solution.Name, solution.Namespace, userID)
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

func (pgdb *pgDB) GetSolutionsList(ctx context.Context, userID string) (*stypes.UserSolutionsList, error) {
	pgdb.log.Infoln("Get solutions list")
	var ret stypes.UserSolutionsList

	ret.Solutions = make([]stypes.UserSolution, 0)

	rows, err := pgdb.qLog.QueryxContext(ctx, "SELECT templates.Name, solutions.name, solutions.namespace, parameters.env, parameters.branch "+
		"FROM solutions JOIN parameters ON solutions.id = parameters.solution_id JOIN templates ON solutions.template_id = templates.ID WHERE solutions.user_id=$1 AND solutions.is_deleted !='true'", userID)
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

func (pgdb *pgDB) DeleteSolution(ctx context.Context, name string, userID string) error {
	pgdb.log.Infoln("Deleting solution")

	res, err := pgdb.eLog.ExecContext(ctx, `UPDATE solutions SET is_deleted = 'true', deleted_at=$1 WHERE name=$2 AND user_id=$3 AND is_deleted != 'true'`, time.Now(), name, userID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if rows == 0 {
		return sErrors.ErrSolutionNotExist()
	}
	return nil
}

func (pgdb *pgDB) GetSolutionsDeployments(ctx context.Context, solutionName string, userID string) (deployments []string, ns *string, err error) {
	pgdb.log.Infoln("Get solution deployments")

	rows, err := pgdb.qLog.QueryxContext(ctx, "SELECT solutions.namespace, deployments.deploy_name "+
		"FROM solutions JOIN deployments ON solutions.id = deployments.solution_id WHERE solutions.name=$1 AND solutions.user_id=$2 AND solutions.is_deleted !='true'", solutionName, userID)
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

func (pgdb *pgDB) GetSolutionsServices(ctx context.Context, solutionName string, userID string) (services []string, ns *string, err error) {
	pgdb.log.Infoln("Get solution services")

	rows, err := pgdb.qLog.QueryxContext(ctx, "SELECT solutions.namespace, services.service_name "+
		"FROM solutions JOIN services ON solutions.id = services.solution_id WHERE solutions.name=$1 AND solutions.user_id=$2 AND solutions.is_deleted !='true'", solutionName, userID)
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
