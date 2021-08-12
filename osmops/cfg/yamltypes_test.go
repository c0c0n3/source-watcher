package cfg

import (
	"testing"
)

var opsConfigValidationFailFixtures = []OpsConfig{
	{TargetDir: "", ConnectionFile: ""},
	{TargetDir: " ", ConnectionFile: "\n"},
	{TargetDir: "valid", ConnectionFile: "\n"},
	{TargetDir: "\t", ConnectionFile: "./val/id"},
}

func TestOpsConfigValidationFail(t *testing.T) {
	for k, d := range opsConfigValidationFailFixtures {
		if got := d.Validate(); got == nil {
			t.Errorf("[%d] want: error; got: valid", k)
		}
	}
}

var opsConfigValidationOkFixtures = []OpsConfig{
	{ConnectionFile: "./"},
	{TargetDir: "", ConnectionFile: "./"},
	{TargetDir: ".", ConnectionFile: "./"},
	{TargetDir: " /a/", ConnectionFile: "/a/b "},
	{TargetDir: "valid", ConnectionFile: "./val/id"},
	{TargetDir: "\tval/id\n", ConnectionFile: "\n/val/id/\t"},
}

func TestOpsConfigValidationOk(t *testing.T) {
	for k, d := range opsConfigValidationOkFixtures {
		if got := d.Validate(); got != nil {
			t.Errorf("[%d] want: valid; got: %s", k, got)
		}
	}
}

var osmConnectionValidationFailFixtures = []OsmConnection{
	{Hostname: "", User: "u", Password: "p"},
	{}, {Hostname: "h", Password: "p"},
}

func TestOsmConnectionValidationFail(t *testing.T) {
	for k, d := range osmConnectionValidationFailFixtures {
		if got := d.Validate(); got == nil {
			t.Errorf("[%d] want: error; got: valid", k)
		}
	}
}

var osmConnectionValidationOkFixtures = []OsmConnection{
	{Hostname: "h:0", User: "u", Password: "p"}, {Hostname: "h:1", User: "u"},
}

func TestOsmConnectionValidationOk(t *testing.T) {
	for k, d := range osmConnectionValidationOkFixtures {
		if got := d.Validate(); got != nil {
			t.Errorf("[%d] want: valid; got: %s", k, got)
		}
	}
}

var kduNsActionValidationFailFixtures = []KduNsAction{
	{},
	{
		Kind:    KduNsActionKind,
		Name:    "x",
		Action:  NsAction.Create,
		VnfName: "x",
	},
	{
		Kind:    KduNsActionKind,
		Name:    "x",
		Action:  NsAction.Create,
		VnfName: "x",
		Kdu: Kdu{
			Model: "x",
		},
	},
	{
		Name:    "x",
		Action:  NsAction.Create,
		VnfName: "x",
		Kdu: Kdu{
			Name:  "x",
			Model: "x",
		},
	},
	{
		Kind:    "ain't right",
		Name:    "x",
		Action:  NsAction.Create,
		VnfName: "x",
		Kdu: Kdu{
			Name:  "x",
			Model: "x",
		},
	},
	{
		Kind:    KduNsActionKind,
		Action:  NsAction.Create,
		VnfName: "x",
		Kdu: Kdu{
			Name:  "x",
			Model: "x",
		},
	},
	{
		Kind:    KduNsActionKind,
		Name:    "x",
		VnfName: "x",
		Kdu: Kdu{
			Name:  "x",
			Model: "x",
		},
	},
	{
		Kind:    KduNsActionKind,
		Name:    "x",
		Action:  "ain't right",
		VnfName: "x",
		Kdu: Kdu{
			Name:  "x",
			Model: "x",
		},
	},
	{
		Kind:   KduNsActionKind,
		Action: NsAction.Create,
		Name:   "x",
		Kdu: Kdu{
			Name:  "x",
			Model: "x",
		},
	},
}

func TestKduNsActionValidationFail(t *testing.T) {
	for k, d := range kduNsActionValidationFailFixtures {
		if got := d.Validate(); got == nil {
			t.Errorf("[%d] want: error; got: valid", k)
		}
	}
}

var kduNsActionValidationOkFixtures = []KduNsAction{
	{
		Kind:    KduNsActionKind,
		Name:    "x",
		Action:  NsAction.Create,
		VnfName: "x",
		Kdu: Kdu{
			Name:  "x",
			Model: "x",
		},
	},
	{
		Kind:    KduNsActionKind,
		Name:    "x",
		Action:  NsAction.Upgrade,
		VnfName: "x",
		Kdu: Kdu{
			Name: "x",
		},
	},
	{
		Kind:    KduNsActionKind,
		Name:    "x",
		Action:  NsAction.Delete,
		VnfName: "x",
		Kdu: Kdu{
			Name: "x",
		},
	},
}

func TestKduNsActionValidationOk(t *testing.T) {
	for k, d := range kduNsActionValidationOkFixtures {
		if got := d.Validate(); got != nil {
			t.Errorf("[%d] want: valid; got: %s", k, got)
		}
	}
}
