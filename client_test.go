package bingx

func newTestClient(do doFunc) *Client {
	client := NewClient("DummyKey", "DummySecret")
	client.Debug = true
	client.do = do
	return client
}
