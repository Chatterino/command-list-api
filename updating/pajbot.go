package updating

import (
	"encoding/json"
	"errors"
	"io/ioutil"
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

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var msg struct {
		Commands []struct {
			ID                    int         `json:"id"`
			Level                 int         `json:"level"`
			MainAlias             string      `json:"main_alias"`
			Aliases               []string    `json:"aliases"`
			Description           interface{} `json:"description"`
			LongDescription       string      `json:"long_description"`
			CdAll                 int         `json:"cd_all"`
			CdUser                int         `json:"cd_user"`
			Enabled               bool        `json:"enabled"`
			Cost                  int         `json:"cost"`
			TokensCost            int         `json:"tokens_cost"`
			CanExecuteWithWhisper bool        `json:"can_execute_with_whisper"`
			SubOnly               bool        `json:"sub_only"`
			ModOnly               bool        `json:"mod_only"`
			ResolveString         string      `json:"resolve_string"`
			Examples              []struct {
				ID          interface{} `json:"id"`
				CommandID   int         `json:"command_id"`
				Title       string      `json:"title"`
				Description string      `json:"description"`
				Messages    []struct {
					Source struct {
						Type string      `json:"type"`
						From string      `json:"from"`
						To   interface{} `json:"to"`
					} `json:"source"`
					Message string `json:"message"`
				} `json:"messages"`
			} `json:"examples"`
			Data struct {
				NumUses      int       `json:"num_uses"`
				AddedBy      string    `json:"added_by"`
				EditedBy     string    `json:"edited_by"`
				LastDateUsed time.Time `json:"last_date_used"`
			} `json:"data"`
		} `json:"commands"`
	}

	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}

	cmds := make([]Command, 0, len(msg.Commands))

	for _, x := range msg.Commands {

		cmds = append(cmds, Command{
			Prefix:      x.MainAlias,
			Description: "",
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

	data, err := ioutil.ReadAll(resp.Body)
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
