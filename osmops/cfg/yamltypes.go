package cfg

import (
	"fmt"
	"strings"

	v "github.com/go-ozzo/ozzo-validation"

	u "github.com/fluxcd/source-watcher/osmops/util"
)

// OpsConfig holds the configuration data needed to scan a repo to find
// supported OSM Ops deployment files and run OSM commands according to
// the data found in those files.
type OpsConfig struct {
	// TargetDir is a path, relative to the repo root, pointing to the
	// directory containing OSM Ops YAML files. Defaults to the repo root
	// if omitted.
	TargetDir string `yaml:"targetDir, omitempty"`

	// ConnectionFile is a path to the file containing OSM connection data.
	// (See `OsmConnection` structure.) This is an absolute path to a separate
	// YAML config file mounted on the pod running OSM Ops, typically through
	// a K8s secret.
	ConnectionFile string `yaml:"connectionFile"`
}

// Validate OpsConfig data read from a YAML file.
// An instance is valid if:
// * TargetDir is not present or if present isn't empty and is a valid path.
// * ConnectionFile isn't empty and is a valid path.
func (d OpsConfig) Validate() error {
	validTargetDir := func(value interface{}) error { // (*)
		s, _ := value.(string)
		if len(s) == 0 {
			return nil
		}
		return u.IsStringPath(value)
	}
	return v.ValidateStruct(&d,
		v.Field(&d.TargetDir, v.By(validTargetDir)),
		v.Field(&d.ConnectionFile, v.By(u.IsStringPath)),
	)

	// (*) the latest ozzo-validation (GH/master) comes w/ conditional
	// validation rules, so when they release it, we could replace our
	// custom validTargetDir w/ e.g.
	//     v.When(d.TargetDir != "", v.By(u.IsStringPath)).Else(v.Nil)
}

// OsmConnection holds the data the OSM client needs to connect to the OSM
// north-bound interface.
type OsmConnection struct {
	Hostname string `yaml:"hostname"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// Validate OsmConnection data read from a YAML file.
// The hostname field must be in the form h:p where h is a DNS name or IP
// address and p is a valid port number---i.e. between 0 and 65535. IP6
// addresses are accepted too but have to be enclosed in square brackets---e.g.
// "[::1]:80", "[::1%lo0]:80".
// The user field must not be empty whereas there are no restrictions for
// the password field.
func (d OsmConnection) Validate() error {
	return v.ValidateStruct(&d,
		v.Field(&d.Hostname, v.By(u.IsHostAndPort)),
		v.Field(&d.User, v.Required),
	)
}

const KduNsActionKind = "KduNsAction"

var NsAction = struct {
	Create  string
	Upgrade string
	Delete  string
}{
	Create:  "create",
	Upgrade: "upgrade",
	Delete:  "delete",
}

func NsActions() []string {
	return []string{NsAction.Create, NsAction.Upgrade, NsAction.Delete}
}

type Kdu struct {
	Name  string      `yaml:"name"`
	Model interface{} `yaml:"model"`
}

func (d Kdu) Validate() error {
	return v.ValidateStruct(&d, v.Field(&d.Name, v.Required))
}

// KduNsAction holds the data in a YAML file that instructs OSM Ops to run
// an NS action on a KDU.
type KduNsAction struct {
	Kind    string `yaml:"kind"`
	Name    string `yaml:"name"`
	Action  string `yaml:"action"`
	VnfName string `yaml:"vnfName"`
	Kdu     Kdu    `yaml:"kdu"`
}

// Validate KduNsAction data read from a YAML file.
// An instance is valid if:
// * Kind has a value of KduNsActionKind.
// * Name, VnfName and Kdu.Name are not empty.
// * Action is one of the NsAction constants---case doesn't matter.
func (d KduNsAction) Validate() error {
	validKind := func(value interface{}) error {
		s, _ := value.(string)
		if s != KduNsActionKind {
			return fmt.Errorf("not %s", KduNsActionKind)
		}
		return nil
	}
	validAction := func(value interface{}) error {
		s, _ := value.(string)
		s = strings.ToLower(s)
		for _, v := range NsActions() { // (*)
			if v == s {
				return nil
			}
		}
		return fmt.Errorf("not an action: %s", s)

		// (*) ideally it should rather be
		//     return v.In(NsActions()...).Validate(s)
		// but I couldn't get it right. When the action is missing,
		// validation passes!
	}
	return v.ValidateStruct(&d,
		v.Field(&d.Kind, v.By(validKind)), // (*)
		v.Field(&d.Name, v.Required),
		v.Field(&d.Action, v.By(validAction)),
		v.Field(&d.VnfName, v.Required),
		v.Field(&d.Kdu),
	)

	// (*) ideally it'd be the In rule, but I couldn't get it right, if
	// there's no Kind, validation passes!
}
