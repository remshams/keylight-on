package utils

import (
	"github.com/remshams/device-control/keylight-control/control"
	"strconv"
)

func FindKeylightWithId(keylights []control.Keylight, keylightId string) *control.Keylight {
	id, err := strconv.Atoi(keylightId)
	if err != nil {
		return nil
	}
	return control.FindKeylightWithId(keylights, id)
}
