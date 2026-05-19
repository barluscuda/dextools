package minecraft

import (
	"bytes"
	"encoding/binary"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestBedrockServerStatus(t *testing.T) {
	host := os.Getenv("MINECRAFT_MCSERVER_DEMO_IP")
	if host == "" {
		t.Skip("skipping live Bedrock status test: MINECRAFT_MCSERVER_DEMO_IP is not set")
	}

	portValue := os.Getenv("MINECRAFT_MCSERVER_DEMO_PORT")
	if portValue == "" {
		t.Skip("skipping live Bedrock status test: MINECRAFT_MCSERVER_DEMO_PORT is not set")
	}

	port, err := strconv.Atoi(portValue)
	if err != nil {
		t.Fatalf("invalid MINECRAFT_MCSERVER_DEMO_PORT %q: %v", portValue, err)
	}

	b := &Bedrock{
		Timeout: 5 * time.Second,
	}

	status, err := b.ServerStatus(host, port)
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

func TestParseBedrockStatusResponse(t *testing.T) {
	packet := buildBedrockStatusPacket(t, "MCPE;NetherGames;594;1.20.81;12;500;123456789;Lobby;Survival;1;19132;19133;")

	status, err := parseBedrockStatusResponse(packet, 42*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if status.Ping != 42*time.Millisecond {
		t.Fatalf("unexpected ping: %v", status.Ping)
	}
	if status.MOTD != "NetherGames" {
		t.Fatalf("unexpected MOTD: %q", status.MOTD)
	}
	if status.Protocol != 594 {
		t.Fatalf("unexpected protocol: %d", status.Protocol)
	}
	if status.Version != "1.20.81" {
		t.Fatalf("unexpected version: %q", status.Version)
	}
	if status.PlayersOnline != 12 || status.PlayersMax != 500 {
		t.Fatalf("unexpected player counts: %d/%d", status.PlayersOnline, status.PlayersMax)
	}
	if status.ServerID != "123456789" {
		t.Fatalf("unexpected server id: %q", status.ServerID)
	}
	if status.LevelName != "Lobby" {
		t.Fatalf("unexpected level name: %q", status.LevelName)
	}
	if status.GameMode != "Survival" || status.GameModeNumeric != 1 {
		t.Fatalf("unexpected game mode: %q/%d", status.GameMode, status.GameModeNumeric)
	}
	if status.IPv4Port != 19132 || status.IPv6Port != 19133 {
		t.Fatalf("unexpected ports: %d/%d", status.IPv4Port, status.IPv6Port)
	}
}

func TestParseBedrockStatusResponseRejectsInvalidNumericFields(t *testing.T) {
	packet := buildBedrockStatusPacket(t, "MCPE;NetherGames;bad;1.20.81;12;500;123456789;Lobby;Survival;1;19132;19133;")

	_, err := parseBedrockStatusResponse(packet, 0)
	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "invalid protocol") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func buildBedrockStatusPacket(t *testing.T, payload string) []byte {
	t.Helper()

	var packet bytes.Buffer
	packet.WriteByte(0x1c)
	if err := binary.Write(&packet, binary.BigEndian, int64(123)); err != nil {
		t.Fatalf("write timestamp: %v", err)
	}
	if err := binary.Write(&packet, binary.BigEndian, int64(456)); err != nil {
		t.Fatalf("write server id: %v", err)
	}
	packet.Write([]byte{
		0x00, 0xff, 0xff, 0x00,
		0xfe, 0xfe, 0xfe, 0xfe,
		0xfd, 0xfd, 0xfd, 0xfd,
		0x12, 0x34, 0x56, 0x78,
	})
	if err := binary.Write(&packet, binary.BigEndian, uint16(len(payload))); err != nil {
		t.Fatalf("write string length: %v", err)
	}
	packet.WriteString(payload)

	return packet.Bytes()
}
