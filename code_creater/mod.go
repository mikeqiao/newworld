// Author: mike.qiao
// File:mod
// Date:2022/3/15 18:05

package code_creater

import (
	"errors"
	"github.com/mikeqiao/newworld/net"
)

type Mod_Name struct {
}

func Route_To_Mod(modName, funcName string, data *net.UserData, req []byte, processor net.Processor) (err error) {
	switch modName {
	case "ModTestOne":
		err = ModName_Handler(funcName, data, req, processor)
	default:
		err = errors.New("nil function")
	}
	return
}
