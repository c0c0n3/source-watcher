package cmd

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	jsoniter "github.com/json-iterator/go"

	"github.com/fluxcd/source-watcher/osmops/cfg"
	u "github.com/fluxcd/source-watcher/osmops/util"
)

func buildScanner(repoRootDir u.AbsPath) (*cfg.KduNsActionRepoScanner, error) {
	if store, err := cfg.NewStore(repoRootDir); err != nil {
		return nil, err
	} else {
		return cfg.NewKduNsActionRepoScanner(store), nil
	}
}

type processor struct {
	ctx context.Context
}

func (p *processor) log() logr.Logger {
	return logr.FromContext(p.ctx)
}

func kduModelToJson(file *cfg.KduNsActionFile) ([]byte, error) {
	if file.Content.Kdu.Model == nil {
		return nil, nil
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary // (*)
	return json.Marshal(file.Content.Kdu.Model)

	// (*) can't use Go's built-in json lib since it will blow up w/
	//    json: unsupported type: map[interface {}]interface{}
	// In fact, the YAML lib deserialises the Model block into a
	//    map[interface {}]interface{}
	// which the built-in json doesn't know how to handle.
	// See:
	// - https://stackoverflow.com/questions/35377477
}

func (p *processor) Process(file *cfg.KduNsActionFile) error {
	kduModel, err := kduModelToJson(file)
	if err != nil {
		return err
	}

	cmd := ""
	if kduModel != nil {
		cmdFmt := "osm ns-action %s --vnf_name %s --kdu_name %s " +
			"--action_name %s --params '%s'"
		cmd = fmt.Sprintf(cmdFmt, file.Content.Name, file.Content.VnfName,
			file.Content.Kdu.Name, file.Content.Action, kduModel)
	} else {
		cmdFmt := "osm ns-action %s --vnf_name %s --kdu_name %s " +
			"--action_name %s"
		cmd = fmt.Sprintf(cmdFmt, file.Content.Name, file.Content.VnfName,
			file.Content.Kdu.Name, file.Content.Action)
	}

	p.log().Info("Processed", "file", file.FilePath.Value(), "command", cmd)

	return nil
}

// TODO dig deep into OSM client code
//   git clone https://osm.etsi.org/gerrit/osm/osmclient
// it looks like we're better off writing the YAML to a file and then using
// '--params_file' instead of '--params'.
// TODO should we use the REST API directly? More work, but it could simplify
// configuration in the end since potentially we could just work w/ plain
// OSM model files instead of having to roll out our own, e.g. KduNsActionFile.

func OsmReconcile(ctx context.Context, repoRootDir string) {
	visitor := &processor{ctx: ctx}

	rootDir, err := u.ParseAbsPath(repoRootDir)
	if err != nil {
		visitor.log().Error(err, "can't convert to abs path")
		return
	}

	scanner, err := buildScanner(rootDir)
	if err != nil {
		visitor.log().Error(err, "can't build scanner")
		return
	}

	errors := scanner.Visit(visitor)
	if len(errors) > 0 {
		for k, e := range errors {
			visitor.log().Error(e, "processing errors", "error", k)
		}
	}
}
