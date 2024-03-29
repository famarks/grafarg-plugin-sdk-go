package live

import (
	"fmt"
	"time"

	"github.com/centrifugal/centrifuge-go"
	"github.com/famarks/grafarg-plugin-sdk-go/backend"
	"github.com/famarks/grafarg-plugin-sdk-go/backend/log"
)

// GrafargLiveClient connects to the GrafargLive server
type GrafargLiveClient struct {
	connected bool
	client    *centrifuge.Client
	lastWarn  time.Time
	channels  map[string]*GrafargLiveChannel
	Log       log.Logger
}

// GrafargLiveChannel allows access to a channel within the server
type GrafargLiveChannel struct {
	id         string
	subscribed bool
	client     *GrafargLiveClient
	sub        *centrifuge.Subscription
}

// Subscribe will subscribe to a channel within the server
func (c *GrafargLiveClient) Subscribe(addr ChannelAddress) (*GrafargLiveChannel, error) {
	id := addr.ToChannelID()
	if !addr.IsValid() {
		return nil, fmt.Errorf("invalid channel: %s", id)
	}

	sub, err := c.client.NewSubscription(id)
	if err != nil {
		return nil, err
	}

	ch := &GrafargLiveChannel{
		id:     id,
		client: c,
		sub:    sub,
	}

	sub.OnSubscribeSuccess(ch)
	sub.OnSubscribeError(ch)
	sub.OnUnsubscribe(ch)

	err = sub.Subscribe()
	if err != nil {
		return nil, err
	}
	return ch, nil
}

// Publish sends the data to the channel
func (c *GrafargLiveChannel) Publish(data []byte) {
	if !c.client.connected {
		if time.Since(c.client.lastWarn) > time.Second*5 {
			c.client.lastWarn = time.Now()
			c.client.Log.Warn("Grafarg live channel not connected", "id", c.id)
		}
		return
	}
	_, err := c.sub.Publish(data)
	if err != nil {
		c.client.Log.Info("error publishing", "error", err)
	}
}

//--------------------------------------------------------------------------------------
// CLIENT
//--------------------------------------------------------------------------------------

type liveClientHandler struct {
	client *GrafargLiveClient
}

func (h *liveClientHandler) OnConnect(c *centrifuge.Client, e centrifuge.ConnectEvent) {
	h.client.Log.Info("Connected to Grafarg live", "clientId", e.ClientID)
	h.client.connected = true
}

func (h *liveClientHandler) OnError(c *centrifuge.Client, e centrifuge.ErrorEvent) {
	h.client.Log.Warn("Grafarg live error", "error", e.Message)
}

func (h *liveClientHandler) OnDisconnect(c *centrifuge.Client, e centrifuge.DisconnectEvent) {
	h.client.Log.Info("Disconnected from Grafarg live", "reason", e.Reason)
	h.client.connected = false
}

//--------------------------------------------------------------------------------------
// Channel
//--------------------------------------------------------------------------------------

// OnSubscribeSuccess is called when the channel is subscribed
func (c *GrafargLiveChannel) OnSubscribeSuccess(sub *centrifuge.Subscription, e centrifuge.SubscribeSuccessEvent) {
	c.subscribed = true
	c.client.Log.Info("Subscribed", "channel", sub.Channel())
}

// OnSubscribeError is called when the channel has an error
func (c *GrafargLiveChannel) OnSubscribeError(sub *centrifuge.Subscription, e centrifuge.SubscribeErrorEvent) {
	c.subscribed = false
	c.client.Log.Warn("Subscription failed", "channel", sub.Channel(), "error", e.Error)
}

// OnUnsubscribe is called when the channel is unsubscribed
func (c *GrafargLiveChannel) OnUnsubscribe(sub *centrifuge.Subscription, e centrifuge.UnsubscribeEvent) {
	c.subscribed = false
	c.client.Log.Info("Unsubscribed", "channel", sub.Channel())
}

// InitGrafargLiveClient starts a chat server
func InitGrafargLiveClient(conn ConnectionInfo) (*GrafargLiveClient, error) {
	url, err := conn.ToWebSocketURL()
	if err != nil {
		return nil, err
	}
	c := centrifuge.New(url, centrifuge.DefaultConfig())

	glc := &GrafargLiveClient{
		client:   c,
		channels: make(map[string]*GrafargLiveChannel),
		Log:      backend.Logger,
	}
	handler := &liveClientHandler{
		client: glc,
	}
	c.OnConnect(handler)
	c.OnError(handler)
	c.OnDisconnect(handler)

	err = c.Connect()
	if err != nil {
		return nil, err
	}

	return glc, nil
}
