package bingx

const (
	baseWsUrl        = "wss://open-api-swap.bingx.com/swap-market"
	baseAccountWsUrl = "wss://open-api-swap.bingx.com/swap-market?listenKey="
)

var (
	wsClients map[string]*WsClient
)

func getWsEndpoint() string {
	return baseWsUrl
}

func getAccountWsEndpoint(listenKey string) string {
	return baseAccountWsUrl + listenKey
}

func getWsClient(endpoint string) *WsClient {
	if c := wsClients[endpoint]; c != nil {
		return c
	}

	return &WsClient{
		config: newWsConfig(endpoint),
	}
}

type DepthEvent struct {
	DataType string      `json:"dataType"`
	Data     interface{} `json:"data"`
	Asks     interface{} `json:"asks"`
	Bids     interface{} `json:"bids"`
	Price    float64     `json:"p"`
	Volume   float64     `json:"v"`
}

type WsDepthHandler func(*DepthEvent)

func NewMarketDethWs(symbol string, level, interval int, handler WsDepthHandler) (*WsClient, error) {
	client := getWsClient(getWsEndpoint())
	if client.conn == nil {
		err := client.serve()
		if err != nil {
			return nil, err
		}
	}

	
}
