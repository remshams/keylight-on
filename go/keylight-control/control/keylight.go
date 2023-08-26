package control

import (
	"net"
)

type LightCommand struct {
	On          *bool
	Brightness  *int
	Temperature *int
}

type Light struct {
	On          bool
	Brightness  int
	Temperature int
}

type Keylight struct {
	Id      int
	Name    string
	Ip      []net.IP
	Port    int
	light   *Light
	adapter KeylightAdapter
}

func (keylight *Keylight) loadLights() error {
	lights, err := keylight.adapter.Load(keylight.Ip, keylight.Port)
	if err != nil {
		return err
	}
	if len(lights) > 0 {
		keylight.light = &lights[0]
	}
	return nil
}

func (keylight *Keylight) setLight(lightCommand LightCommand) error {
	on := lightCommand.On
	if on == nil {
		on = &keylight.light.On
	}
	brightness := lightCommand.Brightness
	if brightness == nil {
		brightness = &keylight.light.Brightness
	}
	temperature := lightCommand.Temperature
	if temperature == nil {
		temperature = &keylight.light.Temperature
	}
	light := Light{
		On:          *on,
		Temperature: *temperature,
		Brightness:  *brightness,
	}
	err := keylight.adapter.Set(keylight.Ip, keylight.Port, []Light{light})
	if err == nil {
		keylight.light = &light
	}
	return err
}
