package minecraft

import (
	"testing"
	"time"
)

func TestBedrockServerStatus(t *testing.T) {
	b := &Bedrock{
		Timeout: 5 * time.Second,
	}

	status, err := b.ServerStatus("play.nethergames.org", 19132)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if status == nil {
		t.Fatal("status is nil")
	}

	if status.Ping <= 0 {
		t.Error("invalid ping")
	}

	if status.MOTD == "" {
		t.Error("empty MOTD")
	}

	if status.Protocol <= 0 {
		t.Error("invalid protocol")
	}

	if status.Version == "" {
		t.Error("empty version")
	}

	if status.PlayersMax < 0 {
		t.Error("invalid players max")
	}

	if status.PlayersOnline < 0 {
		t.Error("invalid players online")
	}

	if status.ServerID == "" {
		t.Error("empty server id")
	}

	if status.LevelName == "" {
		t.Error("empty level name")
	}

	if status.GameMode == "" {
		t.Error("empty game mode")
	}

	if status.IPv4Port <= 0 {
		t.Error("invalid IPv4 port")
	}

	if status.IPv6Port < 0 {
		t.Error("invalid IPv6 port")
	}
}
