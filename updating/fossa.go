package updating

import (
	. "self/prelude"
	"time"
)

type fossaChannelResp struct {
	Channel struct {
		ID              string    `json:"id"`
		Login           string    `json:"login"`
		DisplayName     string    `json:"display_name"`
		Avatar          string    `json:"avatar"`
		Slug            string    `json:"slug"`
		BroadcasterType string    `json:"broadcaster_type"`
		Provider        string    `json:"provider"`
		ProviderID      string    `json:"provider_id"`
		CreatedAt       time.Time `json:"createdAt"`
	} `json:"channel"`
	Parent struct {
		ID       string `json:"id"`
		Type     string `json:"type"`
		ParentID string `json:"parent_id"`
	} `json:"parent"`
}

type fossaCommandsResp struct {
	Roles []struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Default bool   `json:"default"`
	} `json:"roles"`
	Commands []struct {
		EnabledOffline bool          `json:"enabled_offline"`
		EnabledOnline  bool          `json:"enabled_online"`
		ID             string        `json:"id"`
		Name           string        `json:"name"`
		Response       string        `json:"response"`
		Type           string        `json:"type"`
		Aliases        []interface{} `json:"aliases"`
		RoleIDs        []interface{} `json:"role_ids"`
	} `json:"commands"`
}

// FetchFossa loads the public commands from the fossabot api by user login.
func FetchFossa(uid UserId, login UserLogin) ([]Command, error) {
	channel, err := fetchFossaChannel(login)
	if err != nil {
		return nil, err
	}

	channelCmds, err := fetchFossaCommands(channel.Channel.ID)
	if err != nil {
		return nil, err
	}

	cmds := make([]Command, 0, len(channelCmds.Commands))

	for _, x := range channelCmds.Commands {
		cmds = append(cmds, Command{
			Prefix:      "!" + x.Name,
			Description: x.Response,
			Source:      "Fossabot",
		})
	}

	return cmds, nil
}

func fetchFossaChannel(login UserLogin) (*fossaChannelResp, error) {
	msg := &fossaChannelResp{}
	url := `https://fossabot.com/api/v2/cached/channels/by-slug/` + login

	if err := FetchJson(url, msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func fetchFossaCommands(fossaChannelId UserId) (*fossaCommandsResp, error) {
	msg := &fossaCommandsResp{}
	url := `https://fossabot.com/api/v2/cached/channels/` + fossaChannelId + `/commands`

	if err := FetchJson(url, msg); err != nil {
		return nil, err
	}

	return msg, nil
}
