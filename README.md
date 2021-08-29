# Command List API

Accumulates a list of all chat commands available in twitch channels.
Updates commands by using chat bot APIs.

## Architecture

- `updating` decides when to call the chat bot apis. Then it fetches the data and saves it to redis.

- `serve` reads the values out of redis and serves it to the user.

- `prelude` is helper functions.

## Running

### Docker

- Create an `.env` file with your [twitch app](https://dev.twitch.tv/docs/authentication#registration) tokens:
  ```
  TWITCH_CLIENT_ID=abcdefg123456
  TWITCH_CLIENT_SECRET=hijklmnop7890
  ```
- Run `docker compose build` to build.
- Run `docker compose up` to run.

Tip: Use `docker compose run redis redis-cli -h redis monitor` to monitor all redis changes and `docker compose run redis redis-cli -h redis` to clear all redis data.

### Manual

- Run `redis-server` (see [redis.io](https://redis.io/)) in a terminal.
- Set (in a different terminal) the env variables `TWITCH_CLIENT_ID` and `TWITCH_CLIENT_SECRET` to the according values [see](https://dev.twitch.tv/docs/authentication#registration). You can use something like `dotenv` or a start script which sets the variables if you like.
- Run `go build && ./self` to build and run the application.

Tip: Use `redis-cli monitor` to monitor all redis changes and `redis-cli FLUSHALL` to clear all redis data.
