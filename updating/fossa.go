package updating

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	. "self/prelude"
)

type fossaResp struct {
	Channelname string `json:"channelName"`
	Channel     struct {
		Avatar      string `json:"avatar"`
		Provider    string `json:"provider"`
		Providerid  string `json:"providerId"`
		Displayname string `json:"displayName"`
		Login       string `json:"login"`
	} `json:"channel"`
	Commands []struct {
		Minuserlevel int    `json:"minUserlevel"`
		ID           string `json:"_id"`
		Name         string `json:"name"`
		Response     string `json:"response"`
	} `json:"commands"`
}

// FetchFossa loads the public commands from the fossabot api by user login.
func FetchFossa(uid UserId, login UserLogin) ([]Command, error) {
	resp, err := http.Get(
		`https://api-v1.fossabot.com/api/v1/` + login + `/public-commands`)

	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var msg fossaResp
	if err = json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}

	cmds := make([]Command, 0, len(msg.Commands))

	for _, x := range msg.Commands {
		cmds = append(cmds, Command{
			Prefix:      "!" + x.Name,
			Description: x.Response,
			Source:      "Fossabot",
		})
	}

	return cmds, nil
}
