package updating

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	. "self/prelude"
	"sync"
	"time"
)

var (
	pajbotInitialized  = false
	pajbotChannelMutex = sync.Mutex{}
	// channel id -> domain
	pajbotChannels      = map[string]string{}
	errDoesntHavePajbot = errors.New("doesnt have pajbot")
	errPajbotNotReady   = errors.New("pajbot not ready")
)

func init() {
	go func() {
		for {
			if updatePajbotBotList() == nil {
				log.Print("updated pajbot bot list")

				pajbotInitialized = true
				time.Sleep(time.Hour)
			} else {
				time.Sleep(time.Minute * 10)
			}
		}
	}()
}

func FetchPajbot(uid UserId, login UserLogin) ([]Command, error) {
	if !pajbotInitialized {
		return nil, errPajbotNotReady
	}

	domain, ok := getPajbotChannelDomain(uid)
	if !ok {
		return nil, errDoesntHavePajbot
	}

	return fetchPajbotCommands(domain)
}

func fetchPajbotCommands(domain string) ([]Command, error) {
	resp, err := http.Get("https://" + domain + "/api/v1/commands")
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var msg struct {
		Commands []struct {
			Aliases               []string `json:"aliases"`
			CanExecuteWithWhisper bool     `json:"can_execute_with_whisper"`
			CdUser                int      `json:"cd_user"`
			CdAll                 int      `json:"cd_all"`
			Cost                  int      `json:"cost"`
			Data                  struct {
				AddedBy      string    `json:"added_by"`
				EditedBy     string    `json:"edited_by"`
				LastDateUsed time.Time `json:"last_date_used"`
				NumUses      int       `json:"num_uses"`
			} `json:"data"`
			Description interface{} `json:"description"`
			Enabled     bool        `json:"enabled"`
			Examples    []struct {
				CommandID   int         `json:"command_id"`
				Description string      `json:"description"`
				ID          interface{} `json:"id"`
				Messages    []struct {
					Message string `json:"message"`
					Source  struct {
						From string      `json:"from"`
						To   interface{} `json:"to"`
						Type string      `json:"type"`
					} `json:"source"`
				} `json:"messages"`
				Title string `json:"title"`
			} `json:"examples"`
			ID                int         `json:"id"`
			JsonDescription   interface{} `json:"json_description"`
			Level             int         `json:"level"`
			LongDescription   string      `json:"long_description"`
			MainAlias         string      `json:"main_alias"`
			ModOnly           bool        `json:"mod_only"`
			ParsedDescription string      `json:"parsed_description"`
			ResolveString     string      `json:"resolve_string"`
			SubOnly           bool        `json:"sub_only"`
			TokensCost        int         `json:"tokens_cost"`
		} `json:"commands"`
	}

	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}

	cmds := make([]Command, 0, len(msg.Commands))

	for _, x := range msg.Commands {

		cmds = append(cmds, Command{
			Prefix:      x.MainAlias,
			Description: x.ParsedDescription,
			Source:      "Pajbot",
		})
	}

	return cmds, nil
}

func getPajbotChannelDomain(channelID string) (string, bool) {
	pajbotChannelMutex.Lock()
	defer pajbotChannelMutex.Unlock()

	domain, ok := pajbotChannels[channelID]
	return domain, ok
}

func updatePajbotBotList() error {
	resp, err := http.Get(`https://raw.githubusercontent.com/pajbot/pajbot.com/master/static/bots.json`)
	if err != nil {
		log.Println("Failed to fetch pajbot.com bot list: ", err)
		return err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read pajbot.com bot list: ", err)
		return err
	}

	var bots []struct {
		Hoster      string `json:"Hoster"`
		ChannelName string `json:"ChannelName"`
		ChannelID   string `json:"ChannelID"`
		Domain      string `json:"Domain"`
	}
	if err = json.Unmarshal(data, &bots); err != nil {
		log.Println("Failed to parse pajbot.com bot list: ", err)
		return err
	}

	pajbotChannelMutex.Lock()
	defer pajbotChannelMutex.Unlock()

	pajbotChannels = map[string]string{}
	for _, x := range bots {
		pajbotChannels[x.ChannelID] = x.Domain
	}

	return nil
}
