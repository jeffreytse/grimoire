package cmd

import (
	"testing"

	"github.com/jeffreytse/grimoire/internal/config"
)

func TestApplyKey_CoreTrace(t *testing.T) {
	var fs config.FileConfig
	for _, v := range []string{"on", "off", ""} {
		if err := applyKey(&fs, "core.trace", v); err != nil || fs.Core.Trace != v {
			t.Errorf("set %q: err=%v trace=%q", v, err, fs.Core.Trace)
		}
	}
	if err := applyKey(&fs, "core.trace", "loud"); err == nil {
		t.Error("core.trace = loud: want error")
	}

	r := config.Merge([]config.FileConfig{{Core: config.CoreSection{Trace: "off"}}}, []string{"/g.toml"})
	if v, src, err := getKeyResolved(&r, "core.trace"); err != nil || v != "off" || src != "/g.toml" {
		t.Errorf("get: v=%q src=%q err=%v", v, src, err)
	}
}

func TestApplyKey_CoreDisplayStyle(t *testing.T) {
	var fs config.FileConfig
	for _, v := range []string{"rich", "plain", ""} {
		if err := applyKey(&fs, "core.display-style", v); err != nil || fs.Core.DisplayStyle != v {
			t.Errorf("set %q: err=%v style=%q", v, err, fs.Core.DisplayStyle)
		}
	}
	if err := applyKey(&fs, "core.display-style", "fancy"); err == nil {
		t.Error("core.display-style = fancy: want error")
	}
}
