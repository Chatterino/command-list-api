package updating

import (
	. "self/prelude"
	"time"
)

type seChannelResp struct {
	Profile struct {
		Headerimage string `json:"headerImage"`
	} `json:"profile"`
	Provider        string `json:"provider"`
	Broadcastertype string `json:"broadcasterType"`
	Suspended       bool   `json:"suspended"`
	ID              string `json:"_id"`
	Ispartner       bool   `json:"isPartner"`
	Providerid      string `json:"providerId"`
	Displayname     string `json:"displayName"`
	Username        string `json:"username"`
	Avatar          string `json:"avatar"`
	Alias           string `json:"alias"`
	Inactive        bool   `json:"inactive"`
}

type seDefaultCommandsResp []struct {
	Commandid      string `json:"commandId"`
	Command        string `json:"command"`
	AccessLevel    int    `json:"accessLevel"`
	Enabled        bool   `json:"enabled"`
	EnabledOnline  bool   `json:"enabledOnline"`
	EnabledOffline bool   `json:"enabledOffline"`
	ModuleEnabled  bool   `json:"moduleEnabled,omitempty"`
	Moduleid       string `json:"moduleId"`
	Cost           int    `json:"cost"`
	Cooldown       struct {
		User   int `json:"user"`
		Global int `json:"global"`
	} `json:"cooldown"`
	Aliases     []interface{} `json:"aliases"`
	Regex       string        `json:"regex"`
	Description string        `json:"description"`
	Subcommands []string      `json:"subCommands,omitempty"`
}

type seChannelCommandsResp []struct {
	Cooldown struct {
		User   int `json:"user"`
		Global int `json:"global"`
	} `json:"cooldown"`
	Aliases        []interface{} `json:"aliases"`
	Keywords       []interface{} `json:"keywords"`
	Enabled        bool          `json:"enabled"`
	EnabledOnline  bool          `json:"enabledOnline"`
	EnabledOffline bool          `json:"enabledOffline"`
	Hidden         bool          `json:"hidden"`
	Cost           int           `json:"cost"`
	Type           string        `json:"type"`
	AccessLevel    int           `json:"accessLevel"`
	ID             string        `json:"_id"`
	UpdatedAt      time.Time     `json:"updatedAt"`
	CreatedAt      time.Time     `json:"createdAt"`
	Command        string        `json:"command"`
	Reply          string        `json:"reply"`
	Channel        string        `json:"channel"`
	Regex          string        `json:"regex,omitempty"`
}

// FetchStreamElements loads the public commands from the streamelements api by user id.
func FetchStreamElements(uid UserId, login UserLogin) ([]Command, error) {
	channel, err := fetchSeChannel(login)
	if err != nil {
		return nil, err
	}

	defaultCmds, err := fetchSeDefaultCommands(channel.ID)
	if err != nil {
		return nil, err
	}

	channelCmds, err := fetchSeChannelCommands(channel.ID)
	if err != nil {
		return nil, err
	}

	// FOLD
	cmds := make([]Command, 0, len(defaultCmds)+len(channelCmds))

	for _, cmd := range defaultCmds {
		if !cmd.Enabled {
			continue
		}

		cmds = append(cmds, Command{
			Prefix:      "!" + cmd.Command,
			Description: cmd.Description,
			Source:      "StreamElements Default Command",
		})
	}

	for _, cmd := range channelCmds {
		if !cmd.Enabled {
			continue
		}

		if cmd.Hidden {
			continue
		}

		cmds = append(cmds, Command{
			Prefix:      "!" + cmd.Command,
			Description: cmd.Reply,
			Source:      "StreamElements Channel Command",
		})
	}

	return cmds, nil
}

func fetchSeChannel(login UserLogin) (*seChannelResp, error) {
	msg := &seChannelResp{}
	url := `https://api.streamelements.com/kappa/v2/channels/` + login

	if err := FetchJson(url, msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func fetchSeDefaultCommands(seChannelId UserId) (seDefaultCommandsResp, error) {
	var msg seDefaultCommandsResp
	url := `https://api.streamelements.com/kappa/v2/bot/commands/` + seChannelId + `/default`

	if err := FetchJson(url, &msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func fetchSeChannelCommands(seChannelId UserId) (seChannelCommandsResp, error) {
	var msg seChannelCommandsResp
	url := `https://api.streamelements.com/kappa/v2/bot/commands/` + seChannelId

	if err := FetchJson(url, &msg); err != nil {
		return nil, err
	}

	return msg, nil
}
