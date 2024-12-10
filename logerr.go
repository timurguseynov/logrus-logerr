package logerr

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

func Entry(err error, fields ...logrus.Fields) *logrus.Entry {
	return logrus.WithFields(getFields(WithFields(err, fields...)))
}

func WithFields(err error, fields ...logrus.Fields) error {
	if err == nil {
		return nil
	}

	return &logrusError{
		err:    fmt.Errorf("%w", err),
		fields: MergeFields(fields...),
	}
}

func MergeFields(fields ...logrus.Fields) logrus.Fields {
	merged := make(logrus.Fields)
	for i := range fields {
		for k, v := range fields[i] {
			merged[k] = v
		}
	}
	return merged
}

type logrusErrorer interface {
	error
	GetFields() logrus.Fields
}

type logrusError struct {
	err    error
	fields logrus.Fields
}

func (e *logrusError) GetFields() logrus.Fields {
	if e == nil {
		return nil
	}
	return e.fields
}

func (e *logrusError) Error() string {
	if e == nil {
		return ""
	}
	return e.err.Error()
}

func (e *logrusError) Unwrap() error {
	if e == nil || e.err == nil {
		return nil
	}
	return e.err
}

func getFields(err error) logrus.Fields {
	var e logrusErrorer
	if errors.As(err, &e) {
		return mergeFieldsWithFuncs(getFields(errors.Unwrap(e)), e.GetFields())
	}
	return nil
}

func mergeFieldsWithFuncs(fields ...logrus.Fields) logrus.Fields {
	merged := make(logrus.Fields)

	var funcs []string
	for i := range fields {
		for k, v := range fields[i] {
			if k != "func" {
				merged[k] = v
				continue
			}

			funk, ok := v.(string)
			if ok && funk != "" {
				funcs = append(funcs, funk)
			}
		}
	}

	merged["func"] = strings.Join(funcs, ",")

	return merged
}
