package gateio

import (
	"encoding/json"
	"errors"
	"fmt"
	"gocryptotrader/common"
	"gocryptotrader/common/convert"
	"gocryptotrader/currency"
	"gocryptotrader/exchange/asset"
	"gocryptotrader/exchange/order"
	"gocryptotrader/exchange/stream"
	"gocryptotrader/exchange/trade"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	gateioWebsocketEndpoint  = "wss://api.gateio.ws/ws/v4/"
	gateioWebsocketRateLimit = 120
)

// WsConnect initiates a websocket connection
func (g *Gateio) WsConnect() error {
	fmt.Println("gateio ws connect................................")
	if !g.Websocket.IsEnabled() || !g.IsEnabled() {
		return errors.New(stream.WebsocketNotEnabled)
	}
	var dialer websocket.Dialer
	err := g.Websocket.Conn.Dial(&dialer, http.Header{})
	if err != nil {
		return err
	}

	g.Websocket.Wg.Add(1)
	go g.wsReadData()

	// if g.GetAuthenticatedAPISupport(exchange.WebsocketAuthentication) {
	// 	err = g.wsServerSignIn()
	// 	if err != nil {
	// 		g.Websocket.DataHandler <- err
	// 		g.Websocket.SetCanUseAuthenticatedEndpoints(false)
	// 	} else {
	// 		var authsubs []stream.ChannelSubscription
	// 		authsubs, err = g.GenerateAuthenticatedSubscriptions()
	// 		if err != nil {
	// 			g.Websocket.DataHandler <- err
	// 			g.Websocket.SetCanUseAuthenticatedEndpoints(false)
	// 		} else {
	// 			err = g.Websocket.SubscribeToChannels(authsubs)
	// 			if err != nil {
	// 				g.Websocket.DataHandler <- err
	// 				g.Websocket.SetCanUseAuthenticatedEndpoints(false)
	// 			}
	// 		}
	// 	}
	// }

	return nil
}

// func (g *Gateio) wsServerSignIn() error {
// 	nonce := int(time.Now().Unix() * 1000)
// 	sigTemp, err := g.GenerateSignature(strconv.Itoa(nonce))
// 	if err != nil {
// 		return err
// 	}
// 	signature := crypto.Base64Encode(sigTemp)
// 	signinWsRequest := WebsocketRequest{
// 		ID:     g.Websocket.Conn.GenerateMessageID(false),
// 		Channel: "server.sign",
// 		Params: []interface{}{g.API.Credentials.Key, signature, nonce},
// 	}
// 	resp, err := g.Websocket.Conn.SendMessageReturnResponse(signinWsRequest.ID,
// 		signinWsRequest)
// 	if err != nil {
// 		g.Websocket.SetCanUseAuthenticatedEndpoints(false)
// 		return err
// 	}
// 	var response WebsocketAuthenticationResponse
// 	err = json.Unmarshal(resp, &response)
// 	if err != nil {
// 		g.Websocket.SetCanUseAuthenticatedEndpoints(false)
// 		return err
// 	}
// 	if response.Result.Status == "success" {
// 		g.Websocket.SetCanUseAuthenticatedEndpoints(true)
// 		return nil
// 	}
//
// 	return fmt.Errorf("%s cannot authenticate websocket connection: %s",
// 		g.Name,
// 		response.Result.Status)
// }

// wsReadData receives and passes on websocket messages for processing
func (g *Gateio) wsReadData() {
	defer g.Websocket.Wg.Done()

	for {
		resp := g.Websocket.Conn.ReadMessage()
		if resp.Raw == nil {
			return
		}
		err := g.wsHandleData(resp.Raw)
		if err != nil {
			g.Websocket.DataHandler <- err
		}
	}
}

func (g *Gateio) wsHandleData(respRaw []byte) error {
	var result WebsocketResponse
	err := json.Unmarshal(respRaw, &result)
	if err != nil {
		return err
	}

	if result.ID > 0 {
		if g.Websocket.Match.IncomingWithData(result.ID, respRaw) {
			return nil
		}
	}

	if result.Error.Code != 0 {
		if strings.Contains(result.Error.Message, "authentication") {
			g.Websocket.SetCanUseAuthenticatedEndpoints(false)
			return fmt.Errorf("%v - authentication failed: %v", g.Name, err)
		}
		return fmt.Errorf("%v error %s", g.Name, result.Error.Message)
	}

	switch {
	case strings.Contains(result.Channel, "trades"):
		if !g.IsSaveTradeDataEnabled() {
			fmt.Println("ERROR: save trade data not enabled for gateio")
			return nil
		}

		var tradeData WebsocketTrade
		err = json.Unmarshal(result.Result, &tradeData)
		if err != nil {
			return err
		}

		var p currency.Pair
		p, err = currency.NewPairFromString(tradeData.CurrencyPair)
		if err != nil {
			fmt.Println("error new pair from string", result, tradeData.CurrencyPair)
			return err
		}
		var trades []trade.Data
		var tSide order.Side
		tSide, err = order.StringToOrderSide(tradeData.Side)
		if err != nil {
			g.Websocket.DataHandler <- order.ClassificationError{
				Exchange: g.Name,
				Err:      err,
			}
		}
		trades = append(trades, trade.Data{
			Timestamp:    convert.TimeFromUnixTimestampDecimal(tradeData.Time),
			CurrencyPair: p,
			AssetType:    asset.Spot,
			Exchange:     g.Name,
			Price:        tradeData.Price,
			Amount:       tradeData.Amount,
			Side:         tSide,
			TID:          strconv.FormatInt(tradeData.ID, 10),
		})

		return trade.AddTradesToBuffer(g.Name, trades...)
	default:
		g.Websocket.DataHandler <- stream.UnhandledMessageWarning{
			Message: g.Name + stream.UnhandledMessage + string(respRaw),
		}
		return nil
	}
	return nil
}

// GenerateAuthenticatedSubscriptions returns authenticated subscriptions
func (g *Gateio) GenerateAuthenticatedSubscriptions() ([]stream.ChannelSubscription, error) {
	if !g.Websocket.CanUseAuthenticatedEndpoints() {
		return nil, nil
	}
	var channels = []string{"balance.subscribe", "order.subscribe"}
	var subscriptions []stream.ChannelSubscription
	enabledCurrencies, err := g.GetEnabledPairs(asset.Spot)
	if err != nil {
		return nil, err
	}
	for i := range channels {
		for j := range enabledCurrencies {
			subscriptions = append(subscriptions, stream.ChannelSubscription{
				Channel:  channels[i],
				Currency: enabledCurrencies[j],
				Asset:    asset.Spot,
			})
		}
	}
	return subscriptions, nil
}

// GenerateDefaultSubscriptions returns default subscriptions
func (g *Gateio) GenerateDefaultSubscriptions() ([]stream.ChannelSubscription, error) {

	// enabledCurrencies, _ := g.GetEnabledPairs(asset.Spot)
	//
	// pairs := make([]string, 0)
	// for _, c := range enabledCurrencies {
	// 	pairs = append(pairs, c.String())
	// }
	//
	// params := make(map[string]interface{})
	// params["payload"] = pairs
	// fmt.Println("pairs", pairs)
	//
	// return []stream.ChannelSubscription{stream.ChannelSubscription{
	// 	Channel:  "spot.trades",
	// 	Currency: enabledCurrencies[0],
	// 	Params:   params,
	// 	Asset:    asset.Spot,
	// }}, nil

	var channels = []string{"spot.trades"}
	var subscriptions []stream.ChannelSubscription
	enabledCurrencies, err := g.GetEnabledPairs(asset.Spot)
	if err != nil {
		return nil, err
	}
	for i := range channels {
		for j := range enabledCurrencies {
			params := make(map[string]interface{})
			if strings.EqualFold(channels[i], "depth.subscribe") {
				params["limit"] = 30
				params["interval"] = "0.1"
			} else if strings.EqualFold(channels[i], "kline.subscribe") {
				params["interval"] = 60
			}

			fpair, err := g.FormatExchangeCurrency(enabledCurrencies[j],
				asset.Spot)
			if err != nil {
				return nil, err
			}

			subscriptions = append(subscriptions, stream.ChannelSubscription{
				Channel:  channels[i],
				Currency: fpair.Upper(),
				Params:   params,
				Asset:    asset.Spot,
			})
		}
	}
	return subscriptions, nil
}

// Subscribe sends a websocket message to receive data from the channel
func (g *Gateio) Subscribe(channelsToSubscribe []stream.ChannelSubscription) error {
	var pairs []string
	for i := range channelsToSubscribe {
		pairs = append(pairs, channelsToSubscribe[i].Currency.String())
	}

	payload := WebsocketRequest{
		Time:    time.Now().Unix(),
		ID:      g.Websocket.Conn.GenerateMessageID(false),
		Channel: channelsToSubscribe[0].Channel,
		Event:   "subscribe",
		Payload: pairs,
	}

	fmt.Printf(
		"payload time: %s channel: %s event: %s pairs:%v\n",
		payload.Time,
		payload.Channel,
		payload.Event,
		payload.Payload)

	var errs common.Errors
	resp, err := g.Websocket.Conn.SendMessageReturnResponse(payload.ID, payload)
	if err != nil {
		errs = append(errs, err)
	}
	var response WebsocketAuthenticationResponse
	err = json.Unmarshal(resp, &response)
	fmt.Println("response", response)
	if err != nil {
		errs = append(errs, err)
	}
	if response.Result.Status != "success" {
		errs = append(errs, fmt.Errorf("%v could not subscribe to %v",
			g.Name,
			payload))
	}
	g.Websocket.AddSuccessfulSubscriptions(channelsToSubscribe...)
	if errs != nil {
		return errs
	}
	return nil
}

func (g *Gateio) Unsubscribe(channelsToUnsubscribe []stream.ChannelSubscription) error {
	var channelsThusFar []string
	for i := range channelsToUnsubscribe {
		if common.StringDataCompare(channelsThusFar,
			channelsToUnsubscribe[i].Channel) {
			continue
		}

		channelsThusFar = append(channelsThusFar, channelsToUnsubscribe[i].Channel)

		// unsubscribeText := strings.Replace(channelsToUnsubscribe[i].Channel,
		// 	"subscribe",
		// 	"unsubscribe",
		// 	1)

		unsubscribe := WebsocketRequest{
			ID: g.Websocket.Conn.GenerateMessageID(false),
		}

		resp, err := g.Websocket.Conn.SendMessageReturnResponse(unsubscribe.ID,
			unsubscribe)
		if err != nil {
			return err
		}
		var response WebsocketAuthenticationResponse
		err = json.Unmarshal(resp, &response)
		if err != nil {
			return err
		}
		if response.Result.Status != "success" {
			return fmt.Errorf("%v could not subscribe to %v",
				g.Name,
				channelsToUnsubscribe[i].Channel)
		}
	}
	return nil
}

func (g *Gateio) wsGetBalance(currencies []string) (*WsGetBalanceResponse, error) {
	if !g.Websocket.CanUseAuthenticatedEndpoints() {
		return nil, fmt.Errorf("%v not authorised to get balance", g.Name)
	}
	balanceWsRequest := wsGetBalanceRequest{
		ID:     g.Websocket.Conn.GenerateMessageID(false),
		Method: "balance.query",
		Params: currencies,
	}
	resp, err := g.Websocket.Conn.SendMessageReturnResponse(balanceWsRequest.ID, balanceWsRequest)
	if err != nil {
		return nil, err
	}
	var balance WsGetBalanceResponse
	err = json.Unmarshal(resp, &balance)
	if err != nil {
		return &balance, err
	}

	if balance.Error.Message != "" {
		return nil, fmt.Errorf("%s websocket error: %s",
			g.Name,
			balance.Error.Message)
	}

	return &balance, nil
}

func (g *Gateio) wsGetOrderInfo(market string, offset, limit int) (*WebSocketOrderQueryResult, error) {
	// if !g.Websocket.CanUseAuthenticatedEndpoints() {
	return nil, fmt.Errorf("%v not authorised to get order info", g.Name)
	// }
	// ord := WebsocketRequest{
	// 	ID:      g.Websocket.Conn.GenerateMessageID(false),
	// 	Channel: "order.query",
	// 	Payload: []interface{}{
	// 		market,
	// 		offset,
	// 		limit,
	// 	},
	// }
	//
	// resp, err := g.Websocket.Conn.SendMessageReturnResponse(ord.ID, ord)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// var orderQuery WebSocketOrderQueryResult
	// err = json.Unmarshal(resp, &orderQuery)
	// if err != nil {
	// 	return &orderQuery, err
	// }
	//
	// if orderQuery.Error.Message != "" {
	// 	return nil, fmt.Errorf("%s websocket error: %s",
	// 		g.Name,
	// 		orderQuery.Error.Message)
	// }
	//
	// return &orderQuery, nil
}
