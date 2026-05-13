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

	if n < 35 || buf[0] != 0x1c {
		return nil, fmt.Errorf("invalid response")
	}

	offset := 1 + 8 + 8 + 16

	if offset+2 > n {
		return nil, fmt.Errorf("invalid packet")
	}

	strLen := binary.BigEndian.Uint16(buf[offset : offset+2])
	offset += 2

	if int(offset)+int(strLen) > n {
		return nil, fmt.Errorf("invalid string length")
	}

	raw := string(buf[offset : offset+int(strLen)])
	parts := bytes.Split([]byte(raw), []byte(";"))

	if len(parts) < 12 {
		return nil, fmt.Errorf("invalid bedrock response")
	}

	parseInt := func(b []byte) int {
		v, _ := strconv.Atoi(string(b))
		return v
	}

	return &BedrockServerStatusStruct{
		Ping:            ping,
		MOTD:            string(parts[1]),
		Protocol:        parseInt(parts[2]),
		Version:         string(parts[3]),
		PlayersOnline:   parseInt(parts[4]),
		PlayersMax:      parseInt(parts[5]),
		ServerID:        string(parts[6]),
		LevelName:       string(parts[7]),
		GameMode:        string(parts[8]),
		GameModeNumeric: parseInt(parts[9]),
		IPv4Port:        parseInt(parts[10]),
		IPv6Port:        parseInt(parts[11]),
	}, nil
}
