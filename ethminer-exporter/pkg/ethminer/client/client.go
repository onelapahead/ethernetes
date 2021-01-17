package client

type Client struct {
	Host string
	Port int
}

type jsonApiObject struct {
	ID      int    `json:"id"`
	JsonRPC string `json:"jsonrpc"`
}

type jsonRequest struct {
	jsonApiObject
	Method string `json:"method"`
}

type jsonResponse struct {
	jsonApiObject
	Result *jsonResult `json:"result"`
}

type jsonResult map[string]interface{}

func (c *Client) Init() {

}

func (c *Client) Ping() *jsonResponse {

	return nil
}

func (c *Client) GetDetailedStats() *jsonResponse {

	return nil
}
