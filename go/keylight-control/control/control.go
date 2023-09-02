package control

import (
	"errors"

	"github.com/rs/zerolog/log"
)

type KeylightCommand struct {
	Id      int
	Command LightCommand
}

type KeylightControl struct {
	finder    KeylightFinder
	adapter   KeylightAdapter
	store     KeylightStore
	keylights []Keylight
}

func New(finder KeylightFinder, adapter KeylightAdapter, store KeylightStore) KeylightControl {
	return KeylightControl{finder, adapter, store, []Keylight{}}
}

func (control *KeylightControl) LoadOrDiscoverKeylights() []Keylight {
	control.loadKeylights()
	if len(control.keylights) == 0 {
		control.DiscoverKeylights()
	}
	return control.keylights
}

func (control *KeylightControl) DiscoverKeylights() ([]Keylight, error) {
	keylights := control.finder.Discover(control.adapter, control.store)
	control.keylights = keylights
	isSuccess := control.discoverKeylights()
	if isSuccess {
		log.Debug().Msgf("Discovered %d keylights", len((keylights)))
		control.saveKeylights()
		return keylights, nil
	} else {
		log.Debug().Msg("Failed to discover keylights")
		return keylights, errors.New("Failed to load some lights")
	}
}

func (control *KeylightControl) discoverKeylights() bool {
	isSuccess := true
	for i := range control.keylights {
		keylight := &control.keylights[i]
		err := keylight.loadLights()
		if err != nil {
			isSuccess = false
		}
	}
	return isSuccess

}

func (control *KeylightControl) saveKeylights() error {
	err := control.store.Save(control.keylights)
	if err != nil {
		log.Debug().Msg("Failed to save keylights")
	} else {
		log.Debug().Msg("Saved keylights")
	}
	return err
}

func (control *KeylightControl) loadKeylights() bool {
	isSuccess := true
	keylights, err := control.store.Load(control.adapter)
	if err != nil {
		return false
	}
	for i := range keylights {
		keylight := &keylights[i]
		err = keylight.loadLights()
		if err != nil {
			isSuccess = false
		}
	}
	control.keylights = keylights
	log.Debug().Msgf("Loaded %d keylights: %+v", len(control.keylights), control.keylights)
	return isSuccess
}

func (control *KeylightControl) KeylightWithId(id int) *Keylight {
	return control.findKeylight(id)
}

func (control *KeylightControl) SendKeylightCommand(command KeylightCommand) error {
	log.Debug().Msgf("Send command: %+v", command)
	keylight := control.findKeylight(command.Id)
	if keylight == nil {
		return errors.New("Keylight not found")
	}
	keylight.setLight(command.Command)
	log.Debug().Msg("Send command success")
	return nil
}

func (control *KeylightControl) findKeylight(id int) *Keylight {
	var selectedKeylight *Keylight
	for _, light := range control.keylights {
		if light.Id == id {
			selectedKeylight = &light
			return selectedKeylight
		}
	}
	return nil

}
