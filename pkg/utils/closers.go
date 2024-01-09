package utils

import (
	"errors"
	"io"
)

type MultiCloser struct {
	cls []func() error
}

func (m *MultiCloser) AddCloser(c io.Closer) {
	m.cls = append(m.cls, c.Close)
}

func (m *MultiCloser) AddFunc(fn func()) {
	m.cls = append(m.cls, func() error {
		fn()
		return nil
	})
}

func (m *MultiCloser) Close() error {
	var errs []error

	for i := range m.cls {
		if err := m.cls[len(m.cls)-1-i](); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
