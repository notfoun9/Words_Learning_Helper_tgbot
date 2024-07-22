package telegram

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"telegram-bot/lib/e"
)

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func (c *Client) Update(offset int, limit int) ([]Update, error) {
	query := url.Values{}
	query.Add("offset", strconv.Itoa(offset))
	query.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, query)
	if err != nil {
		return nil, err
	}

	var res UpdateResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *Client) SendMessage(chatID int, text string) error {
	query := url.Values{}
	query.Add("chat_id", strconv.Itoa(chatID))
	query.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, query)

	return e.Wrap("cant send message", err)
}

func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	defer func() { err = e.Wrap("can't do request", err) }()

	u2 := "https://" + path.Join(c.host, c.basePath, method)

	req, err := http.NewRequest(http.MethodGet, u2, nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	log.Println(req.URL)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func newBasePath(token string) string {
	return "bot" + token
}
