package iotaconn

const (
	AlphanetEndpointURL = "https://api.iota-rebased-alphanet.iota.cafe"
	LocalnetEndpointURL = "http://localhost:9000"

	AlphanetWebsocketEndpointURL = "wss://api.iota-rebased-alphanet.iota.cafe"
	LocalnetWebsocketEndpointURL = "ws://localhost:9000"

	AlphanetFaucetURL = "https://faucet.iota-rebased-alphanet.iota.cafe/gas"
	LocalnetFaucetURL = "http://localhost:9123/gas"
)

const (
	ChainIdentifierAlphanet = "ed959793"
	// localnet doesn't have a fixed ChainIdentifier
)

func FaucetURL(apiURL string) string {
	switch apiURL {
	case AlphanetEndpointURL:
		return AlphanetFaucetURL
	case LocalnetEndpointURL:
		return LocalnetFaucetURL
	default:
		panic("unspecified FaucetURL")
	}
}
