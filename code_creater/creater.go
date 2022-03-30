// Author: mike.qiao
// File:creater
// Date:2022/3/21 17:33

package code_creater

import (
	"fmt"
	"golang.org/x/tools/go/packages"
	"strings"
)

func Do() {
	cfg := &packages.Config{
		Mode: packages.NeedSyntax,
		// TODO: Need to think about constants in test files. Maybe write type_string_test.go
		// in a separate pass? For later.
		Tests:      false,
		BuildFlags: []string{fmt.Sprintf("-tags=%s", strings.Join(tags, " "))},
	}
	pkgs, err := packages.Load(cfg, patterns...)
}
