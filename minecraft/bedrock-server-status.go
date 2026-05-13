package minecraft

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"time"
)

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

func (b *Bedrock) ServerStatus(host string, port int) (*BedrockServerStatusStruct, error) {
	addr := fmt.Sprintf("%s:%d", host, port)

	conn, err := net.DialTimeout("udp", addr, b.Timeout)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(b.Timeout))

	var packet bytes.Buffer
	packet.WriteByte(0x01)

	timestamp := time.Now().UnixMilli()
	_ = binary.Write(&packet, binary.BigEndian, timestamp)

	packet.Write([]byte{
		0x00, 0xff, 0xff, 0x00,
		0xfe, 0xfe, 0xfe, 0xfe,
		0xfd, 0xfd, 0xfd, 0xfd,
		0x12, 0x34, 0x56, 0x78,
	})

	_ = binary.Write(&packet, binary.BigEndian, int64(1))

	start := time.Now()

	if _, err := conn.Write(packet.Bytes()); err != nil {
		return nil, err
	}

	buf := make([]byte, 2048)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	ping := time.Since(start)

	return parseBedrockStatusResponse(buf[:n], ping)
}

func parseBedrockStatusResponse(packet []byte, ping time.Duration) (*BedrockServerStatusStruct, error) {
	if len(packet) < 35 || packet[0] != 0x1c {
		return nil, fmt.Errorf("invalid response")
	}

	offset := 1 + 8 + 8 + 16
	if offset+2 > len(packet) {
		return nil, fmt.Errorf("invalid packet")
	}

	strLen := int(binary.BigEndian.Uint16(packet[offset : offset+2]))
	offset += 2
	if offset+strLen > len(packet) {
		return nil, fmt.Errorf("invalid string length")
	}

	parts := bytes.Split(packet[offset:offset+strLen], []byte(";"))
	if len(parts) < 12 {
		return nil, fmt.Errorf("invalid bedrock response")
	}
	if string(parts[0]) != "MCPE" {
		return nil, fmt.Errorf("invalid edition: %q", string(parts[0]))
	}

	parseInt := func(field string, raw []byte) (int, error) {
		value, err := strconv.Atoi(string(raw))
		if err != nil {
			return 0, fmt.Errorf("invalid %s: %q", field, string(raw))
		}
		return value, nil
	}

	protocol, err := parseInt("protocol", parts[2])
	if err != nil {
		return nil, err
	}
	playersOnline, err := parseInt("players online", parts[4])
	if err != nil {
		return nil, err
	}
	playersMax, err := parseInt("players max", parts[5])
	if err != nil {
		return nil, err
	}
	gameModeNumeric, err := parseInt("game mode numeric", parts[9])
	if err != nil {
		return nil, err
	}
	ipv4Port, err := parseInt("IPv4 port", parts[10])
	if err != nil {
		return nil, err
	}
	ipv6Port, err := parseInt("IPv6 port", parts[11])
	if err != nil {
		return nil, err
	}

	return &BedrockServerStatusStruct{
		Ping:            ping,
		MOTD:            string(parts[1]),
		Protocol:        protocol,
		Version:         string(parts[3]),
		PlayersOnline:   playersOnline,
		PlayersMax:      playersMax,
		ServerID:        string(parts[6]),
		LevelName:       string(parts[7]),
		GameMode:        string(parts[8]),
		GameModeNumeric: gameModeNumeric,
		IPv4Port:        ipv4Port,
		IPv6Port:        ipv6Port,
	}, nil
}
