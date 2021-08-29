package updating

import (
	"log"
	. "self/prelude"

	"github.com/go-redis/redis/v8"
)

var (
	allProviders = []Provider{
		{id: "fossa", fetch: FetchFossa},
		{id: "streamelements", fetch: FetchStreamElements},
	}
)

type Provider struct {
	id    string
	fetch func(uid UserId, login UserLogin) ([]Command, error)
}

// GetId returns the unique id of this provider.
func (p Provider) GetId() string {
	return p.id
}

// Fetch fetches the commands of this provider/user.
func (p Provider) Fetch(uid UserId, login UserLogin) ([]Command, error) {
	return p.fetch(uid, login)
}

// GetRedisKey gets the redis key under which the commands of this provider/user are stored.
func (p Provider) GetRedisKey(uid UserId) string {
	return GetRedisKey4("twitch", uid, "commands", p.id)
}

func updateCommands(rdb *redis.Client, uid UserId, login UserLogin) {
	for _, prv := range allProviders {
		cmds, err := prv.Fetch(uid, login)
		if err != nil {
			log.Printf("Error loading %s commands for %s %s", prv.GetId(), uid, login)
		}

		persistCommands(rdb, prv.GetRedisKey(uid), cmds)
	}

	// merge all command json blobs into one
	keys := []string{}
	for _, prv := range allProviders {
		keys = append(keys, prv.GetRedisKey(uid))
	}

	mergeCommandLists(rdb, keys, GetRedisKey4("twitch", uid, "commands", "all"))
}
