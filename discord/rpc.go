package discord

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
)

const (
    OpHandshake = 0
    OpFrame     = 1
    OpClose     = 2
    OpPing      = 3
    OpPong      = 4
)

type Client struct {
	sock net.Conn
	appid string
}

type activityPayload struct {
	Cmd   string         `json:"cmd"`
	Nonce string         `json:"nonce"`
	Args  activityArgs   `json:"args"`
}

type activityArgs struct {
	PID      int      `json:"pid"`
	Activity *Activity `json:"activity,omitempty"`
}

type Activity struct {
	Details string `json:"details,omitempty"`
	State   string `json:"state,omitempty"`

	Assets *Assets `json:"assets,omitempty"`

	Party      *Party      `json:"party,omitempty"`
	Timestamps *Timestamps `json:"timestamps,omitempty"`
	Secrets    *Secrets    `json:"secrets,omitempty"`

	Buttons []*Button `json:"buttons,omitempty"`
}

type Assets struct {
	LargeImage string `json:"large_image,omitempty"`
	LargeText  string `json:"large_text,omitempty"`

	SmallImage string `json:"small_image,omitempty"`
	SmallText  string `json:"small_text,omitempty"`
}

type Timestamps struct {
	Start int64 `json:"start,omitempty"`
	End   int64 `json:"end,omitempty"`
}

type Button struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

type Secrets struct {
	Match    string `json:"match,omitempty"`
	Join     string `json:"join,omitempty"`
	Spectate string `json:"spectate,omitempty"`
}

type Party struct {
	ID         string `json:"id,omitempty"`
	Players    int    `json:"size,omitempty"`
	MaxPlayers int    `json:"max,omitempty"`
}

func (client *Client) Connect(appid string) error {
	client.appid = appid
	connected := false

	runtime := os.Getenv("XDG_RUNTIME_DIR")
	if runtime == "" { runtime = "/tmp" }

	for i := range 10 {
        path := fmt.Sprintf("%s/discord-ipc-%d", runtime, i)

        conn, err := net.Dial("unix", path)

        if err == nil {
			client.sock = conn
			connected = true
			break
        }
    }

	if !connected { return errors.New("fuck") }

	return client.writePacket(OpHandshake, map[string]any{
        "v": 1,
        "client_id": client.appid,
    })
}

func (c *Client) SetActivity(a *Activity) error {
    return c.writePacket(OpFrame, activityPayload {
		Cmd: "SET_ACTIVITY",
		Nonce: "hentai",
		Args: activityArgs{
			PID: os.Getpid(),
			Activity: a,
		},
	})
}

func (c *Client) Close() error {
	c.SetActivity(nil)
	if c.sock == nil {
		return nil
	}

	err := c.sock.Close()
	c.sock = nil

	return err
}

func (c *Client) writePacket(op int32, v any) error {
    payload, err := json.Marshal(v)
    if err != nil { return err }

    header := make([]byte, 8)

    binary.LittleEndian.PutUint32(header[0:], uint32(op))
    binary.LittleEndian.PutUint32(header[4:], uint32(len(payload)))

    if _, err := c.sock.Write(header); err != nil {
        return err
    }

    _, err = c.sock.Write(payload)
    return err
}
