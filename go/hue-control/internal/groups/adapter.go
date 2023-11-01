package groups

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hue-control/internal/bridges"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

const path = "http://%s/api/%s/groups"

type GroupDtoById = map[string]GroupDto

type GroupDto struct {
	Name   string
	Lights []string
}

func (groupDto GroupDto) toGroup(id string) Group {
	return Group{
		id:     id,
		name:   groupDto.Name,
		lights: groupDto.Lights,
	}
}

type GroupHttpAdapter struct {
	bridge bridges.Bridge
}

func InitGroupHttpAdapter(bridge bridges.Bridge) GroupHttpAdapter {
	return GroupHttpAdapter{bridge}
}

func (adapter GroupHttpAdapter) All() ([]Group, error) {
	req, client, cancel, err := adapter.requestWithTimeout(
		http.MethodGet,
		fmt.Sprintf(path, adapter.bridge.GetIp(), adapter.bridge.GetApiKey()),
		nil,
		nil,
	)
	defer cancel()
	var response *http.Response
	if err == nil {
		response, err = client.Do(req)
	}
	if err != nil || response.StatusCode >= 300 {
		log.Error().Msg("Could not load groups")
		return nil, errors.New("Could not load groups")
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error().Msg("Could not load lights")
		return nil, errors.New("Could not load groups")
	}
	defer response.Body.Close()
	var groupResponseDto GroupDtoById
	err = json.Unmarshal(body, &groupResponseDto)
	if err != nil {
		log.Error().Msg("Could not parse groups")
		return nil, errors.New("Could not parse groups")
	}
	groups := []Group{}
	if len(groupResponseDto) > 0 {
		for id, groupDto := range groupResponseDto {
			groups = append(groups, groupDto.toGroup(id))
		}
	}
	return groups, nil
}

func (adapter GroupHttpAdapter) requestWithTimeout(method string, url string, body io.Reader, timeout *time.Duration) (*http.Request, *http.Client, context.CancelFunc, error) {
	defaultTimeout := 2 * time.Second
	requestTimeout := timeout
	if requestTimeout == nil {
		requestTimeout = &defaultTimeout
	}
	ctx, cancel := context.WithTimeout(context.Background(), *requestTimeout)
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	return req, client, cancel, err
}