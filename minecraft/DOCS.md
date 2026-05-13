# minecraft

`minecraft` currently provides a small client for querying Minecraft Bedrock server status over UDP.

## Quick Start

Use the package directly:

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/barluscuda/dextools/minecraft"
)

func main() {
	client := &minecraft.Bedrock{
		Timeout: 5 * time.Second,
	}

	status, err := client.ServerStatus("play.nethergames.org", 19132)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("MOTD:", status.MOTD)
	fmt.Println("Version:", status.Version)
	fmt.Printf("Players: %d/%d\n", status.PlayersOnline, status.PlayersMax)
	fmt.Println("Ping:", status.Ping)
}
```

## Behavior

- `Bedrock.ServerStatus` sends a UDP Bedrock status ping to the target host and port.
- `Bedrock.Timeout` is used for both dialing and the read deadline.
- The response parser expects a valid Bedrock unconnected pong payload with the `MCPE` edition marker.
- Invalid numeric fields in the server response return an error instead of silently becoming `0`.
- `Ping` is measured from request write to response read.

## Available Types

The package currently exposes:

- `Bedrock`
- `BedrockServerStatusStruct`

## Client

### `Bedrock`

`Bedrock` is the client used to query a Bedrock server.

```go
type Bedrock struct {
	Timeout time.Duration
}
```

## Status Response

### `BedrockServerStatusStruct`

`ServerStatus` returns a `*BedrockServerStatusStruct` with the parsed server response.

```go
type BedrockServerStatusStruct struct {
	Ping            time.Duration
	MOTD            string
	Protocol        int
	Version         string
	PlayersOnline   int
	PlayersMax      int
	ServerID        string
	LevelName       string
	GameMode        string
	GameModeNumeric int
	IPv4Port        int
	IPv6Port        int
}
```

Field meanings:

- `Ping`: round-trip duration for the status request.
- `MOTD`: server message of the day.
- `Protocol`: Bedrock protocol version number.
- `Version`: Bedrock version string reported by the server.
- `PlayersOnline`: currently connected player count.
- `PlayersMax`: maximum configured player count.
- `ServerID`: server identifier from the Bedrock response.
- `LevelName`: world or level name.
- `GameMode`: game mode name.
- `GameModeNumeric`: numeric game mode identifier.
- `IPv4Port`: IPv4 listener port reported by the server.
- `IPv6Port`: IPv6 listener port reported by the server.

## API

### `func (b *Bedrock) ServerStatus(host string, port int) (*BedrockServerStatusStruct, error)`

`ServerStatus` connects to `host:port`, sends a Bedrock status ping, parses the response, and returns the structured result.

Common error cases:

- UDP dial failures
- request timeout or read timeout
- malformed status packets
- unexpected response payloads

## Testing

Unit tests for the Bedrock status client live in [bedrock-server-status_test.go](./bedrock-server-status_test.go).
