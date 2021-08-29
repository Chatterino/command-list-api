package updating

import (
	"encoding/json"
	"log"
	. "self/prelude"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	updateReqs = make(chan UserId, 100)
	printOnce  = false
)

// HintUpdate considers updating the commands of the given user.
func HintUpdate(uid UserId) {
	select {
	case updateReqs <- uid:
	default:
		if !printOnce {
			printOnce = true
			log.Println("update hits overflowing")
		}
	}
}

// Run receives updates from the `updateReqs` channel and updates the commands accordingly.
func Run(rdb *redis.Client) {
	for {
		uid := <-updateReqs

		if RedisCooldown(rdb, GetRedisKey3("twitch", uid, "cmd-cd")) != nil {
			continue
		}

		login, err := getLogin(rdb, uid)
		if err != nil {
			log.Printf("error getting login for %s: %s\n", uid, err)
			continue
		}

		updateCommands(rdb, uid, login)

		log.Println("updated all commands for ", uid)
	}
}

// Takes the `cmds` and serializes them to redis as json.
func persistCommands(rdb *redis.Client, redisKey string, cmds []Command) {
	b, err := json.Marshal(cmds)

	if err != nil {
		log.Println("Error while marshalling commands to json", err)
		return
	}

	rdb.Set(RedisCtx, redisKey, b, time.Hour*100)
	log.Println("persisted", redisKey)
}

// Takes json blobs out of the `sourceKeys` and merges them into `targetKey`.
func mergeCommandLists(rdb *redis.Client, sourceKeys []string, targetKey string) {
	cmds := make([]Command, 0)

	for _, key := range sourceKeys {
		// get from redis
		jsonBlob, err := rdb.Get(RedisCtx, key).Result()

		if err != nil {
			continue
		}

		var tmpCmds []Command
		err = json.Unmarshal([]byte(jsonBlob), &tmpCmds)

		if err != nil {
			log.Println("error while marshalling commands to json", err)
			continue
		}

		cmds = append(cmds, tmpCmds...)
	}

	b, err := json.Marshal(cmds)
	if err != nil {
		log.Println("error while marshalling commands to json", err)

		return
	}

	rdb.Set(RedisCtx, targetKey, string(b), time.Hour*24)
}
