package impl

import (
	"io"
	"reflect"

	"errors"

	"git.containerum.net/ch/solutions/pkg/server"

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

func (u *serverImpl) Close() error {
	var errs []error
	s := reflect.ValueOf(u.svc)
	closer := reflect.TypeOf((*io.Closer)(nil)).Elem()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
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