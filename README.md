### Twitch CLI Testing

#### Prerequisite(s)
Requires to have TwitchCLI installed. You can download the Twitch CLI [here](https://dev.twitch.tv/docs/cli/).

#### Live Testing
Start a websocket server with `twitch event websocket start-server`. See [documentation](https://dev.twitch.tv/docs/cli/websocket-event-command/) for use of the CLI.

Start `govern` with `go run govern.go -addr=localhost:8080`.

Events:
    Online Event:
    `twitch event trigger stream.online --transport=websocket`
    Force Reconnect:
    `twitch event websocket reconnect`