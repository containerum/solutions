package impl

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/url"

	"git.containerum.net/ch/solutions/pkg/db"
	"git.containerum.net/ch/solutions/pkg/sErrors"
	"git.containerum.net/ch/solutions/pkg/server"
	"git.containerum.net/ch/solutions/pkg/utils"
	kube_types "github.com/containerum/kube-client/pkg/model"
	"github.com/containerum/utils/httputil"
	"github.com/google/uuid"
	"github.com/json-iterator/go"
)

const (
	unableToCreate = "unable to create %s %s: %s"
)

func parseSolutionConfig(ctx context.Context, s *serverImpl, solutionPath string, solutionReq kube_types.UserSolution) (*server.Solution, error) {
	solutionConfigFile, err := s.svc.DownloadClient.DownloadFile(ctx, fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/.containerum.json", solutionPath, solutionReq.Branch))
	if err != nil {
		return nil, err
	}

	solutionRandEnvTmpl, err := template.New("solution_rand").Funcs(template.FuncMap{
		"rand_string":       utils.RandString,
		"rand_string_minus": utils.RandStringWithMinus,
	}).Parse(string(solutionConfigFile))
	if err != nil {
		return nil, err
	}

	var solutionBuf bytes.Buffer
	if err := solutionRandEnvTmpl.Execute(&solutionBuf, nil); err != nil {
		return nil, err
	}

	var solutionConfig *server.Solution
	if err = jsoniter.Unmarshal(solutionBuf.Bytes(), &solutionConfig); err != nil {
		return nil, err
	}

	if len(solutionConfig.Env) == 0 {
		solutionReq.Env = make(map[string]string, 0)
	}

	for k, v := range solutionReq.Env {
		solutionConfig.Env[k] = v
	}
	return solutionConfig, nil
}

func createSolution(ctx context.Context, s *serverImpl, solutionConfig *server.Solution, templateID, solutionUUID string, solutionReq kube_types.UserSolution) error {
	solutionEnvironments, err := jsoniter.MarshalToString(solutionConfig.Env)
	if err != nil {
		return err
	}

	s.log.Debugln("Creating solution")
	if err := s.svc.DB.Transactional(ctx, func(ctx context.Context, tx db.DB) error {
		return s.svc.DB.AddSolution(ctx, solutionReq, httputil.MustGetUserID(ctx), templateID, solutionUUID, solutionEnvironments)
	}); err != nil {
		return s.handleDBError(err)
	}
	return nil
}

func parseResource(ctx context.Context, s *serverImpl, resourceConfig *server.ConfigFile, solutionConfig *server.Solution, solutionPath string, solutionReq kube_types.UserSolution) (*bytes.Buffer, error) {
	s.log.Infof("Creating %s %s", resourceConfig.Type, resourceConfig.Name)
	s.log.Debugln("Downloading resource")
	resF, err := s.svc.DownloadClient.DownloadFile(ctx, fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s", solutionPath, solutionReq.Branch, resourceConfig.Name))
	if err != nil {
		s.log.Debugln(err)
		return nil, fmt.Errorf(unableToCreate, resourceConfig.Type, resourceConfig.Name, err)
	}

	s.log.Debugln("Setting envs to resource config")
	resTmpl, err := template.New("res").Parse(string(resF))
	if err != nil {
		s.log.Debugln(err)
		return nil, fmt.Errorf(unableToCreate, resourceConfig.Type, resourceConfig.Name, err)
	}

	var resParsed bytes.Buffer
	err = resTmpl.Execute(&resParsed, solutionConfig.Env)
	if err != nil {
		s.log.Debugln(err)
		return nil, fmt.Errorf(unableToCreate, resourceConfig.Type, resourceConfig.Name, err)
	}
	return &resParsed, nil
}

func createDeployment(ctx context.Context, s *serverImpl, resourceConfig *server.ConfigFile, solutionName, solutionNamespace string, parsedRes bytes.Buffer) error {
	var parsedDeploy kube_types.Deployment
	err := jsoniter.Unmarshal(parsedRes.Bytes(), &parsedDeploy)
	if err != nil {
		s.log.Debugln(err)
		return fmt.Errorf(unableToCreate, resourceConfig.Type, resourceConfig.Name, err)
	}
	parsedDeploy.SolutionID = solutionName
	if err = s.svc.ResourceClient.CreateDeployment(ctx, solutionNamespace, parsedDeploy); err != nil {
		s.log.Debugln(err)
		return fmt.Errorf(unableToCreate, resourceConfig.Type, resourceConfig.Name, err)
	}
	return err
}

func createService(ctx context.Context, s *serverImpl, resourceConfig *server.ConfigFile, solutionName, solutionNamespace string, parsedRes bytes.Buffer) error {
	var parsedService kube_types.Service
	err := jsoniter.Unmarshal(parsedRes.Bytes(), &parsedService)
	if err != nil {
		s.log.Debugln(err)
		return fmt.Errorf(unableToCreate, resourceConfig.Type, resourceConfig.Name, err)
	}
	parsedService.SolutionID = solutionName
	if err = s.svc.ResourceClient.CreateService(ctx, solutionNamespace, parsedService); err != nil {
		s.log.Debugln(err)
		return fmt.Errorf(unableToCreate, resourceConfig.Type, resourceConfig.Name, err)
	}
	return err
}

func rollbackSolution(ctx context.Context, s *serverImpl, solutionName, solutionNamespace string) {
	s.log.Infoln("No resources was created. Deleting solution...")
	if err := s.svc.DB.Transactional(ctx, func(ctx context.Context, tx db.DB) error {
		err := s.svc.DB.CompletelyDeleteSolution(ctx, solutionNamespace, solutionName)
		return err
	}); err != nil {
		s.log.Errorln(err)
	}
}

func (s *serverImpl) RunSolution(ctx context.Context, solutionReq kube_types.UserSolution) (*kube_types.RunSolutionResponse, error) {
	s.log.Infoln("Running solution ", solutionReq.Name)
	s.log.Debugln("Getting template info from DB")
	solutionTemplate, err := s.svc.DB.GetTemplate(ctx, solutionReq.Template)
	if err = s.handleDBError(err); err != nil {
		return nil, err
	}

	solutionURL, err := url.Parse(solutionTemplate.URL)
	if err != nil {
		return nil, err
	}

	solutionPath := solutionURL.Path[1:]

	s.log.Debugln("Parsing solution config")
	solutionConfig, err := parseSolutionConfig(ctx, s, solutionPath, solutionReq)
	if err != nil {
		return nil, err
	}

	solutionUUID := uuid.New().String()

	err = createSolution(ctx, s, solutionConfig, solutionTemplate.ID, solutionUUID, solutionReq)
	if err != nil {
		return nil, err
	}

	ret := kube_types.RunSolutionResponse{
		Errors:     []string{},
		Created:    0,
		NotCreated: 0,
	}

	s.log.Debugln("Creating solution resources")
	for _, f := range solutionConfig.Run {
		parsedRes, err := parseResource(ctx, s, &f, solutionConfig, solutionPath, solutionReq)
		if err != nil {
			ret.Errors = append(ret.Errors, fmt.Sprintf(unableToCreate, f.Type, f.Name, err))
			continue
		}
		switch f.Type {
		case "deployment":
			if err := createDeployment(ctx, s, &f, solutionReq.Name, solutionReq.Namespace, *parsedRes); err != nil {
				ret.Errors = append(ret.Errors, fmt.Sprintf(unableToCreate, f.Type, f.Name, err))
				continue
			}
		case "service":
			if err := createService(ctx, s, &f, solutionReq.Name, solutionReq.Namespace, *parsedRes); err != nil {
				ret.Errors = append(ret.Errors, fmt.Sprintf(unableToCreate, f.Type, f.Name, err))
				continue
			}
		default:
			ret.Errors = append(ret.Errors, fmt.Sprintf("Unknown resource type: %v. Skipping.", f.Type))
			continue
		}
		ret.Created++
	}

	if ret.Created == 0 {
		rollbackSolution(ctx, s, solutionReq.Name, solutionReq.Namespace)
		return nil, sErrors.ErrUnableCreateSolution().AddDetails(ret.Errors...)
	}

	ret.NotCreated = len(ret.Errors)

	s.log.Infoln("Solution has been created")
	return &ret, nil
}

func (s *serverImpl) DeleteSolution(ctx context.Context, namespace, solutionName string) error {
	s.log.Infoln("Deleting solution ", solutionName)
	solution, err := s.svc.DB.GetSolution(ctx, namespace, solutionName)
	if err := s.handleDBError(err); err != nil {
		return err
	}

	if err := s.svc.ResourceClient.DeleteDeployments(ctx, solution.Namespace, solution.Name); err != nil {
		return err
	}

	if err := s.svc.ResourceClient.DeleteServices(ctx, solution.Namespace, solution.Name); err != nil {
		return err
	}

	s.log.Debugln("Deleting solution")
	if err := s.svc.DB.Transactional(ctx, func(ctx context.Context, tx db.DB) error {
		return s.svc.DB.DeleteSolution(ctx, solution.Namespace, solution.Name)
	}); err != nil {
		return s.handleDBError(err)
	}

	s.log.Debugln("Solution deleted")
	return nil
}

func (s *serverImpl) GetSolutionsList(ctx context.Context, isAdmin bool) (*kube_types.UserSolutionsList, error) {
	resp, err := s.svc.DB.GetSolutionsList(ctx, httputil.MustGetUserID(ctx))
	if err != nil {
		return nil, err
	}

	if !isAdmin {
		for i := range resp.Solutions {
			resp.Solutions[i].ID = ""
		}
	}

	return resp, nil
}

func (s *serverImpl) GetNamespaceSolutionsList(ctx context.Context, namespace string, isAdmin bool) (*kube_types.UserSolutionsList, error) {
	resp, err := s.svc.DB.GetNamespaceSolutionsList(ctx, namespace)
	if err != nil {
		return nil, err
	}

	if !isAdmin {
		for i := range resp.Solutions {
			resp.Solutions[i].ID = ""
		}
	}

	return resp, nil
}

func (s *serverImpl) GetSolution(ctx context.Context, namespace, solutionName string, isAdmin bool) (*kube_types.UserSolution, error) {
	resp, err := s.svc.DB.GetSolution(ctx, namespace, solutionName)
	if err != nil {
		return nil, err
	}

	if !isAdmin {
		resp.ID = ""
	}

	return resp, nil
}

func (s *serverImpl) GetSolutionDeployments(ctx context.Context, namespace, solutionName string) (*kube_types.DeploymentsList, error) {
	solution, err := s.svc.DB.GetSolution(ctx, namespace, solutionName)
	if err := s.handleDBError(err); err != nil {
		return nil, err
	}

	userdepl, err := s.svc.KubeAPIClient.GetUserDeployments(ctx, solution.Namespace, solution.Name)
	if err != nil {
		return nil, err
	}

	return userdepl, nil
}

func (s *serverImpl) GetSolutionServices(ctx context.Context, namespace, solutionName string) (*kube_types.ServicesList, error) {
	solution, err := s.svc.DB.GetSolution(ctx, namespace, solutionName)
	if err := s.handleDBError(err); err != nil {
		return nil, err
	}

	usersvc, err := s.svc.KubeAPIClient.GetUserServices(ctx, solution.Namespace, solution.Name)
	if err != nil {
		return nil, err
	}

	return usersvc, nil
}

func (s *serverImpl) DeleteUserSolutions(ctx context.Context) error {
	if err := s.svc.DB.Transactional(ctx, func(ctx context.Context, tx db.DB) error {
		return s.svc.DB.CompletelyDeleteUserSolutions(ctx, httputil.MustGetUserID(ctx))
	}); err != nil {
		return s.handleDBError(err)
	}

	s.log.Debugln("Solutions deleted")
	return nil
}

func (s *serverImpl) DeleteNamespaceSolutions(ctx context.Context, namespace string) error {
	if err := s.svc.DB.Transactional(ctx, func(ctx context.Context, tx db.DB) error {
		return s.svc.DB.CompletelyDeleteNamespaceSolutions(ctx, namespace)
	}); err != nil {
		return s.handleDBError(err)
	}

	s.log.Debugln("Solutions deleted")
	return nil
}
