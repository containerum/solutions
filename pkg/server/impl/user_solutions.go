package impl

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"strings"

	cherry "git.containerum.net/ch/kube-client/pkg/cherry/solutions"

	"net/url"

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

func (s *serverImpl) DownloadSolutionConfig(ctx context.Context, solutionReq stypes.UserSolution) (solutionFile []byte, solutionName *string, err error) {
	s.log.Infoln("Downloading solution config")
	solutionAvailable, err := s.svc.DB.GetAvailableSolution(ctx, solutionReq.Template)
	if err = s.handleDBError(err); err != nil {
		return nil, nil, err
	}
	if solutionAvailable == nil {
		return nil, nil, cherry.ErrSolutionNotExist()
	}

	solutionURL, err := url.Parse(solutionAvailable.URL)
	if err != nil {
		return nil, nil, err
	}

	sName := strings.TrimSpace(solutionURL.Path[1:])

	solutionF, err := s.svc.DownloadClient.DownloadFile(ctx, fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/.containerum.json", sName, solutionReq.Branch))
	if err != nil {
		return nil, nil, err
	}
	s.log.Infoln("Solution config downloaded")
	return solutionF, &sName, nil
}

func (s *serverImpl) ParseSolutionConfig(ctx context.Context, solutionBody []byte, solutionReq stypes.UserSolution) (solutionConfig *server.Solution, solutionUUID *string, err error) {
	s.log.Infoln("Parsing solution config")
	solutionTmpl, err := template.New("solution").Funcs(template.FuncMap{
		"rand_string": utils.RandString,
	}).Parse(string(solutionBody))
	if err != nil {
		return nil, nil, err
	}

	var solutionBuf bytes.Buffer
	err = solutionTmpl.Execute(&solutionBuf, nil)
	if err != nil {
		return nil, nil, err
	}

	err = jsoniter.Unmarshal(solutionBuf.Bytes(), &solutionConfig)
	if err != nil {
		return nil, nil, err
	}

	if len(solutionConfig.Env) == 0 {
		solutionReq.Env = make(map[string]string, 0)
	}

	s.log.Debugln("Setting env")
	solutionConfig.Env[NamespaceKey] = solutionReq.Namespace

	for k, v := range solutionReq.Env {
		solutionConfig.Env[k] = v
	}

	if _, set := solutionConfig.Env[VolumeKey]; !set { // use default volume name format if volume name not specified
		solutionConfig.Env[VolumeKey] = fmt.Sprintf("%s-volume", solutionReq.Namespace)
	}

	solutionConfig.Env[OwnerKey] = server.MustGetUserID(ctx)

	sUUID := uuid.New().String()
	environments, err := jsoniter.Marshal(solutionConfig.Env)
	if err != nil {
		return nil, nil, err
	}

	err = s.svc.DB.Transactional(ctx, func(ctx context.Context, tx models.DB) error {
		err = s.svc.DB.AddSolution(ctx, solutionReq, server.MustGetUserID(ctx), sUUID, string(environments))
		return err
	})
	if err = s.handleDBError(err); err != nil {
		return nil, nil, err
	}
	s.log.Infoln("Solution config parsed")
	return solutionConfig, &sUUID, nil
}

func (s *serverImpl) CreateSolutionResources(ctx context.Context, solutionConfig server.Solution, solutionReq stypes.UserSolution, solutionName string, solutionUUID string) error {
	s.log.Infoln("Creating solution resources")
	for _, f := range solutionConfig.Run {
		s.log.Debugf("Creating %s %s", f.Type, f.Name)

		resF, err := s.svc.DownloadClient.DownloadFile(ctx, fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s", solutionName, solutionReq.Branch, f.Name))
		if err != nil {
			return err
		}

		resTmpl, err := template.New("res").Parse(string(resF))
		if err != nil {
			return err
		}

		var resParsed bytes.Buffer
		err = resTmpl.Execute(&resParsed, solutionConfig.Env)
		if err != nil {
			return err
		}

		var resMetaJSON server.ResName
		err = jsoniter.Unmarshal(resParsed.Bytes(), &resMetaJSON)
		if err != nil {
			return err
		}

		switch f.Type {
		case "deployment":
			convertedDeploy, err := s.svc.ConverterClient.ConvertDeployment(ctx, resParsed.String())
			if err != nil {
				return err
			}

			err = s.svc.ResourceClient.CreateDeployment(ctx, solutionReq.Namespace, *convertedDeploy)
			if err != nil {
				return err
			}

			err = s.svc.DB.Transactional(ctx, func(ctx context.Context, tx models.DB) error {
				err = s.svc.DB.AddDeployment(ctx, resMetaJSON.Metadata.Name, solutionUUID)
				return err
			})
			if err := s.handleDBError(err); err != nil {
				return err
			}
		case "service":
			convertedService, err := s.svc.ConverterClient.ConvertService(ctx, resParsed.String())
			if err != nil {
				return err
			}
			err = s.svc.ResourceClient.CreateService(ctx, solutionReq.Namespace, *convertedService)
			if err != nil {
				return err
			}
			err = s.svc.DB.Transactional(ctx, func(ctx context.Context, tx models.DB) error {
				err = s.svc.DB.AddService(ctx, resMetaJSON.Metadata.Name, solutionUUID)
				return err
			})
			if err := s.handleDBError(err); err != nil {
				return err
			}
		default:
			s.log.Debugln("Unknown resource type: ", f.Type)
		}
		if err != nil {
			return err
		}
	}
	s.log.Infoln("All solution resources has been created")
	return nil
}

func (s *serverImpl) DeleteSolution(ctx context.Context, solution string) error {
	depl := make([]string, 0)
	svc := make([]string, 0)
	var ns *string

	err := s.svc.DB.Transactional(ctx, func(ctx context.Context, tx models.DB) error {
		var err error
		depl, ns, err = s.svc.DB.GetUserSolutionsDeployments(ctx, solution)
		return err
	})
	if err := s.handleDBError(err); err != nil {
		return err
	}

	err = s.svc.DB.Transactional(ctx, func(ctx context.Context, tx models.DB) error {
		var err error
		svc, _, err = s.svc.DB.GetUserSolutionsServices(ctx, solution)
		return err
	})
	if err := s.handleDBError(err); err != nil {
		return err
	}

	for _, r := range depl {
		err = s.svc.ResourceClient.DeleteDeployment(ctx, *ns, r)
		if err != nil {
			s.log.Infoln(err)
		}
	}

	for _, r := range svc {
		err = s.svc.ResourceClient.DeleteService(ctx, *ns, r)
		if err != nil {
			s.log.Infoln(err)
		}
	}

	err = s.svc.DB.Transactional(ctx, func(ctx context.Context, tx models.DB) error {
		var err error
		err = s.svc.DB.DeleteSolution(ctx, solution)
		return err
	})
	if err := s.handleDBError(err); err != nil {
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
	if err := s.handleDBError(err); err != nil {
		return nil, err
	}

	if ns == nil || len(depl) == 0 {
		return &stypes.DeploymentsList{make([]*interface{}, 0)}, nil
	}

	userdepl, err := s.svc.KubeAPIClient.GetUserDeployments(ctx, *ns, depl)
	if err != nil {
		return nil, err
	}

	return userdepl, nil
}

func (s *serverImpl) GetUserSolutionServices(ctx context.Context, solutionName string) (*stypes.ServicesList, error) {
	svc, ns, err := s.svc.DB.GetUserSolutionsServices(ctx, solutionName)
	if err := s.handleDBError(err); err != nil {
		return nil, err
	}

	if ns == nil || len(svc) == 0 {
		return &stypes.ServicesList{Services: make([]*interface{}, 0)}, nil
	}

	usersvc, err := s.svc.KubeAPIClient.GetUserServices(ctx, *ns, svc)
	if err != nil {
		return nil, err
	}

	return usersvc, nil
}
