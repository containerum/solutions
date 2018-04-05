package impl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/url"
	"strings"

	stypes "git.containerum.net/ch/json-types/solutions"
	"git.containerum.net/ch/solutions/pkg/models"
	"git.containerum.net/ch/solutions/pkg/server"
	"git.containerum.net/ch/solutions/pkg/utils"
	"github.com/google/uuid"
	"github.com/json-iterator/go"
)

const (
	NamespaceKey = "NS"
	VolumeKey    = "VOLUME"
	OwnerKey     = "OWNER"
)

func (s *serverImpl) RunSolution(ctx context.Context, solutionReq stypes.UserSolution) error {
	solutionAvailable, err := s.svc.DB.GetAvailableSolution(ctx, solutionReq.Template)
	if err != nil {
		return err
	}

	url, err := url.Parse(solutionAvailable.URL)
	if err != nil {
		return err
	}

	if solutionReq.Branch != "" {
		solutionReq.Branch = strings.TrimSpace(solutionReq.Branch)
	} else {
		solutionReq.Branch = "master"
	}
	sName := strings.TrimSpace(url.Path[1:])

	//TODO
	sName = "Iliad/redmine-postgresql-solution"

	solutionF, err := s.svc.DownloadClient.DownloadFile(ctx, fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/.containerum.json", sName, solutionReq.Branch))
	if err != nil {
		return err
	}

	solutionTmpl, err := template.New("solution").Funcs(template.FuncMap{
		"rand_string": utils.RandString,
	}).Parse(string(solutionF))
	if err != nil {
		return err
	}

	var solutionBuf bytes.Buffer
	err = solutionTmpl.Execute(&solutionBuf, nil)
	if err != nil {
		return err
	}

	var solutionJSON server.Solution
	err = json.Unmarshal(solutionBuf.Bytes(), &solutionJSON)

	if len(solutionJSON.Env) == 0 {
		solutionReq.Env = make(map[string]string, 0)
	}

	s.log.Debugln("Setting env")
	solutionJSON.Env[NamespaceKey] = solutionReq.Namespace

	for k, v := range solutionReq.Env {
		solutionJSON.Env[k] = v
	}

	if _, set := solutionJSON.Env[VolumeKey]; !set { // use default volume name format if volume name not specified
		solutionJSON.Env[VolumeKey] = fmt.Sprintf("%s-volume", solutionReq.Namespace)
	}

	solutionJSON.Env[OwnerKey] = server.MustGetUserID(ctx)

	solutionUUID := uuid.New().String()
	enviroments, err := jsoniter.Marshal(solutionJSON.Env)

	err = s.svc.DB.Transactional(ctx, func(ctx context.Context, tx models.DB) error {
		err := s.svc.DB.AddSolution(ctx, solutionReq, server.MustGetUserID(ctx), solutionUUID, string(enviroments))
		return err
	})
	if err != nil {
		return err
	}

	//Creating all resources from solution
	for _, f := range solutionJSON.Run {
		s.log.Debugf("Creating %s %s", f.Type, f.Name)
		resF, err := s.svc.DownloadClient.DownloadFile(ctx, fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s", sName, solutionReq.Branch, f.Name))
		if err != nil {
			return err
		}

		resTmpl, err := template.New("res").Parse(string(resF))
		if err != nil {
			return err
		}

		var resParsed bytes.Buffer
		err = resTmpl.Execute(&resParsed, solutionJSON.Env)
		if err != nil {
			return err
		}

		fmt.Println()

		var resMetaJSON server.ResName
		err = json.Unmarshal(resParsed.Bytes(), &resMetaJSON)

		switch f.Type {
		case "deployment":
			s.log.Debugln("Deployment sent to kube-api")
			err := s.svc.KubeAPI.CreateDeployment(ctx, solutionReq.Namespace, resParsed.String())
			if err != nil {
				return err
			}

			fmt.Println("TEST", resParsed.String())

			err = s.svc.DB.Transactional(ctx, func(ctx context.Context, tx models.DB) error {
				err := s.svc.DB.AddDeployment(ctx, resMetaJSON.Metadata.Name, solutionUUID)
				return err
			})
			if err != nil {
				return err
			}
		case "service":
			s.log.Debugln("Service sent to kube-api")
			err := s.svc.KubeAPI.CreateService(ctx, solutionReq.Namespace, resParsed.String())
			if err != nil {
				return err
			}
			err = s.svc.DB.Transactional(ctx, func(ctx context.Context, tx models.DB) error {
				err := s.svc.DB.AddService(ctx, resMetaJSON.Metadata.Name, solutionUUID)
				return err
			})
			if err != nil {
				return err
			}
		default:
			s.log.Debugln("Unknown resource type: ", f.Type)
		}
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func (s *serverImpl) DeleteSolution(ctx context.Context, solution string) error {
	depl := make([]string, 0)
	svc := make([]string, 0)
	var ns *string

	err := s.svc.DB.Transactional(ctx, func(ctx context.Context, tx models.DB) error {
		var err error
		depl, ns, err = s.svc.DB.GetUserSolutionsServices(ctx, solution)
		return err
	})
	if err != nil {
		return err
	}

	err = s.svc.DB.Transactional(ctx, func(ctx context.Context, tx models.DB) error {
		var err error
		svc, _, err = s.svc.DB.GetUserSolutionsServices(ctx, solution)
		return err
	})
	if err != nil {
		return err
	}

	for _, r := range depl {
		err := s.svc.KubeAPI.DeleteDeployment(ctx, *ns, r)
		if err != nil {
			s.log.Infoln(err)
		}
	}

	for _, r := range svc {
		err := s.svc.KubeAPI.DeleteService(ctx, *ns, r)
		if err != nil {
			s.log.Infoln(err)
		}
	}

	err = s.svc.DB.Transactional(ctx, func(ctx context.Context, tx models.DB) error {
		var err error
		err = s.svc.DB.DeleteSolution(ctx, solution)
		return err
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *serverImpl) GetUserSolutionsList(ctx context.Context) (*stypes.UserSolutionsList, error) {
	resp, err := s.svc.DB.GetUserSolutionsList(ctx, server.MustGetUserID(ctx))
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *serverImpl) GetUserSolutionDeployments(ctx context.Context, solutionName string) (*stypes.DeploymentsList, error) {
	depl, ns, err := s.svc.DB.GetUserSolutionsDeployments(ctx, solutionName)
	if err != nil {
		return nil, err
	}

	userdepl, err := s.svc.KubeAPI.GetUserDeployments(ctx, *ns, depl)
	if err != nil {
		return nil, err
	}

	return userdepl, nil
}

func (s *serverImpl) GetUserSolutionServices(ctx context.Context, solutionName string) (*stypes.ServicesList, error) {
	depl, ns, err := s.svc.DB.GetUserSolutionsServices(ctx, solutionName)
	if err != nil {
		return nil, err
	}

	usersvc, err := s.svc.KubeAPI.GetUserServices(ctx, *ns, depl)
	if err != nil {
		return nil, err
	}

	return usersvc, nil
}