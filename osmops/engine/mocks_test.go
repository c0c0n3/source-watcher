package engine

import (
	"context"
	"errors"
	"path/filepath"
	"sort"

	"github.com/fluxcd/source-watcher/osmops/cfg"
	"github.com/fluxcd/source-watcher/osmops/nbic"
	"github.com/go-logr/logr"
)

// logr.Logger implementation

type logCollector struct {
	entries []logEntry
}

type logEntry struct {
	msg    string
	params map[interface{}]interface{}
	err    error
}

func (c *logCollector) append(err error, msg string, kvs ...interface{}) {
	e := logEntry{msg: msg, params: map[interface{}]interface{}{}, err: err}
	for i := 0; i < len(kvs)-1; i += 2 {
		k, v := kvs[i], kvs[i+1]
		e.params[k] = v
	}

	c.entries = append(c.entries, e)
}

func (c *logCollector) Enabled() bool {
	return true
}

func (c *logCollector) Info(msg string, keysAndValues ...interface{}) {
	c.append(nil, msg, keysAndValues...)
}

func (c *logCollector) Error(err error, msg string, keysAndValues ...interface{}) {
	c.append(err, msg, keysAndValues...)
}

func (c *logCollector) V(level int) logr.Logger {
	return c
}

func (c *logCollector) WithValues(keysAndValues ...interface{}) logr.Logger {
	return c
}

func (c *logCollector) WithName(name string) logr.Logger {
	return c
}

// logr.Logger & context factory functions

func newCtx(logger *logCollector) context.Context {
	return logr.NewContext(context.TODO(), logger)
}

func newLogCollector() *logCollector {
	return &logCollector{entries: []logEntry{}}
}

// logCollector utils

func (c *logCollector) countEntries() int {
	return len(c.entries)
}

func (c *logCollector) msgAt(ix int) string {
	return c.entries[ix].msg
}

func (c *logCollector) errAt(ix int) error {
	return c.entries[ix].err
}

func (c *logCollector) sortProcessedFileNames() []string {
	names := []string{}
	for _, e := range c.entries {
		if e.msg == processingMsg {
			if path, ok := e.params[fileLogKey]; ok {
				if p, ok := path.(string); ok {
					name := filepath.Base(p)
					names = append(names, name)
				}
			}
		}
	}
	sort.Strings(names)
	return names
}

func (c *logCollector) sortErrorFileNames() []string {
	names := []string{}
	for _, e := range c.entries {
		if e.msg == processingErrMsg {
			if err, ok := e.err.(*cfg.VisitError); ok {
				name := filepath.Base(err.AbsPath)
				names = append(names, name)
			}
		}
	}
	sort.Strings(names)
	return names
}

// nbic.Workflow implementation

type mockCreateOrUpdate struct {
	dataMap map[string]*nbic.NsInstanceContent
}

func newMockNbicWorkflow() *mockCreateOrUpdate {
	return &mockCreateOrUpdate{
		dataMap: map[string]*nbic.NsInstanceContent{},
	}
}

func (m *mockCreateOrUpdate) CreateOrUpdateNsInstance(data *nbic.NsInstanceContent) error {
	m.dataMap[data.KduName] = data
	if data.KduName == "k2" {
		return errors.New("k2")
	}
	return nil
}

// mockCreateOrUpdate utils

func (m *mockCreateOrUpdate) hasProcessedKdu(name string) bool {
	_, ok := m.dataMap[name]
	return ok
}

func (m *mockCreateOrUpdate) dataFor(kduName string) *nbic.NsInstanceContent {
	if data, ok := m.dataMap[kduName]; ok {
		return data
	}
	return nil
}

func (m *mockCreateOrUpdate) lookupParam(kduName string, paramName string) interface{} {
	if data := m.dataFor(kduName); data != nil {
		if ps, ok := data.KduParams.(map[interface{}]interface{}); ok {
			if v, ok := ps[paramName]; ok {
				return v
			}
		}
	}
	return nil
}
