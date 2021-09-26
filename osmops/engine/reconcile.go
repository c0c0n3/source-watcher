package engine

import (
	"context"

	"github.com/go-logr/logr"
	jsoniter "github.com/json-iterator/go"

	"github.com/fluxcd/source-watcher/osmops/cfg"
	"github.com/fluxcd/source-watcher/osmops/nbic"
	u "github.com/fluxcd/source-watcher/osmops/util"
)

type Engine struct {
	ctx       context.Context
	opsConfig *cfg.Store
	nbic      nbic.Workflow
}

func newNbic(opsConfig *cfg.Store) (nbic.Workflow, error) {
	hp, err := u.ParseHostAndPort(opsConfig.OsmConnection().Hostname)
	if err != nil {
		return nil, err
	}

	conn := nbic.Connection{
		Address: *hp,
		Secure:  false,
	}
	usrCreds := nbic.UserCredentials{
		Username: opsConfig.OsmConnection().User,
		Password: opsConfig.OsmConnection().Password,
		Project:  opsConfig.OsmConnection().Project,
	}

	return nbic.New(conn, usrCreds)
}

func newProcessor(ctx context.Context, repoRootDir string) (*Engine, error) {
	rootDir, err := u.ParseAbsPath(repoRootDir)
	if err != nil {
		return nil, err
	}

	store, err := cfg.NewStore(rootDir)
	if err != nil {
		return nil, err
	}

	client, err := newNbic(store)
	if err != nil {
		return nil, err
	}

	return &Engine{
		ctx:       ctx,
		opsConfig: store,
		nbic:      client,
	}, nil
}

func log(ctx context.Context) logr.Logger {
	return logr.FromContext(ctx)
}

func (p *Engine) log() logr.Logger {
	return log(p.ctx)
}

func (p *Engine) repoScanner() *cfg.KduNsActionRepoScanner {
	return cfg.NewKduNsActionRepoScanner(p.opsConfig)
}

func kduParamsToJson(file *cfg.KduNsActionFile) ([]byte, error) {
	if file.Content.Kdu.Params == nil {
		return nil, nil
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary // (*)
	return json.Marshal(file.Content.Kdu.Params)

	// (*) can't use Go's built-in json lib since it will blow up w/
	//    json: unsupported type: map[interface {}]interface{}
	// In fact, the YAML lib deserialises the Params block into a
	//    map[interface {}]interface{}
	// which the built-in json doesn't know how to handle.
	// See:
	// - https://stackoverflow.com/questions/35377477
}

const (
	processingMsg    = "processing"
	fileLogKey       = "file"
	engineInitErrMsg = "can't initialize reconcile engine"
	processingErrMsg = "processing errors"
	errorLogKey      = "error"
)

func (p *Engine) Process(file *cfg.KduNsActionFile) error {
	p.log().Info(processingMsg, fileLogKey, file.FilePath.Value())

	kduParams, err := kduParamsToJson(file) // TODO bug! zap this!!
	if err != nil {
		return err
	}

	data := nbic.NsInstanceContent{
		Name:           file.Content.Name,
		Description:    file.Content.Description,
		NsdName:        file.Content.NsdName,
		VnfName:        file.Content.VnfName,
		VimAccountName: file.Content.VimAccountName,
		KduName:        file.Content.Kdu.Name,
		KduParams:      kduParams,
	}
	return p.nbic.CreateOrUpdateNsInstance(&data)
}

// New instantiates an Engine to reconcile the state of the OSM deployment
// with that declared in the OSM GitOps files found in the specified repo.
func New(ctx context.Context, repoRootDir string) (*Engine, error) {
	engine, err := newProcessor(ctx, repoRootDir)
	if err != nil {
		log(ctx).Error(err, engineInitErrMsg)
		return nil, err
	}
	return engine, nil
}

// Reconcile looks for OSM GitOps files in the repo and, for each file
// found, it calls OSM NBI to reach the deployment state declared in the
// file.
func (p *Engine) Reconcile() {
	errors := p.repoScanner().Visit(p)
	if len(errors) > 0 {
		for k, e := range errors {
			p.log().Error(e, processingErrMsg, errorLogKey, k)
		}
	}
}
