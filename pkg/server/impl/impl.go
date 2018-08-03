package impl

import (
	"io"
	"reflect"

	"errors"

	"git.containerum.net/ch/solutions/pkg/db"
	"git.containerum.net/ch/solutions/pkg/sErrors"
	"git.containerum.net/ch/solutions/pkg/server"

	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type serverImpl struct {
	svc server.Services
	log *logrus.Entry
}

// NewSolutionsImpl returns a main Solutions implementation
func NewSolutionsImpl(services server.Services) server.SolutionsService {
	return &serverImpl{
		svc: services,
		log: logrus.WithField("component", "solutions_impl"),
	}
}

func (s *serverImpl) Close() error {
	var errs []error
	sv := reflect.ValueOf(s.svc)
	closer := reflect.TypeOf((*io.Closer)(nil)).Elem()
	for i := 0; i < sv.NumField(); i++ {
		f := sv.Field(i)
		if f.Type().ConvertibleTo(closer) {
			errs = append(errs, f.Convert(closer).Interface().(io.Closer).Close())
		}
	}
	var strerr string
	for _, v := range errs {
		if v != nil {
			strerr += v.Error() + ";"
		}
	}
	return errors.New(strerr)
}

func (s *serverImpl) handleDBError(err error) error {
	switch err {
	case nil:
		return nil
	case db.ErrTransactionRollback, db.ErrTransactionCommit, db.ErrTransactionBegin:
		s.log.WithError(err).Error("db transaction error")
		return err
	default:
		if pqerr, ok := err.(*pq.Error); ok {
			switch pqerr.Code {
			case "23505": //unique_violation
				return sErrors.ErrResourceAlreadyExists()
			default:
				s.log.WithError(pqerr)
			}
		}
		s.log.WithError(err).Error("db error")
		return err
	}
}
