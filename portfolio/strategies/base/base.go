package base

import (
	"encoding/json"
	"fmt"
	"gocryptotrader/config"
	"gocryptotrader/currency"
	"gocryptotrader/data"
	"gocryptotrader/eventtypes"
	"gocryptotrader/eventtypes/event"
	"gocryptotrader/eventtypes/signal"
	"gocryptotrader/exchange/order"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"
)

// Strategy is base implementation of the Handler interface
type Strategy struct {
	Name  string
	ID    int
	NumID int
	// currencySettings          *ExchangeAssetPairSettings
	pair                      currency.Pair
	exchange                  string
	weight                    decimal.Decimal
	direction                 order.Side
	useSimultaneousProcessing bool
	usingExchangeLevelFunding bool
	Debug                     bool
	dropFeatures              []string
}

func (s *Strategy) SetDropFeatures() {
	tmpUrl := fmt.Sprintf("http://localhost:8000/drop_features")
	req, err := http.NewRequest("GET", tmpUrl, nil)
	req.URL.RawQuery = fmt.Sprintf("model=%s", s.GetLabel())
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
	drop := make([]string, 0)
	err = json.NewDecoder(resp.Body).Decode(&drop)
	s.dropFeatures = drop
}

func (s *Strategy) GetPrediction(fe FactorEngineHandler) float64 {
	tmpUrl := fmt.Sprintf("http://localhost:8000/predict")
	req, err := http.NewRequest("GET", tmpUrl, nil)

	params := fe.ToQueryParams()
	params["risked_quote"] = 12.0
	rawParams := ""
	for k, v := range params {
		if rawParams == "" {
			rawParams = fmt.Sprintf("%s=%f", k, v)
		} else {
			rawParams = fmt.Sprintf("%s&%s=%f", rawParams, k, v)
		}
	}

	req.URL.RawQuery = fmt.Sprintf("model=%s&%s", s.GetLabel(), rawParams)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	f, _ := strconv.ParseFloat(string(body), 64)
	return f
}

func (s *Strategy) SetName(name string) {
	s.Name = name
}

func (s *Strategy) GetSettings() *config.StrategySetting {
	return &config.StrategySetting{
		Side:    order.Side(s.GetDirection()),
		Capture: s.Name,
		Pair:    s.pair.Upper(),
	}
}

func (s *Strategy) GetLabel() string {
	return fmt.Sprintf("%s@%s@%s", s.Name, s.pair, s.direction)
}

func (s *Strategy) GetNumID() int {
	return s.NumID
}

func (s *Strategy) SetNumID(i int) {
	s.NumID = i
}

func (s *Strategy) GetName() string {
	return s.Name
}

func (s *Strategy) GetID() int {
	return s.ID
}

func (s *Strategy) GetDirection() order.Side {
	return s.direction
}

// GetBaseData returns the non-interface version of the Handler
func GetBaseData(d data.Handler) (signal.Signal, error) {
	if d == nil {
		return signal.Signal{}, eventtypes.ErrNilArguments
	}
	latest := d.Latest()
	if latest == nil {
		return signal.Signal{}, eventtypes.ErrNilEvent
	}
	return signal.Signal{
		Base: event.Base{
			Offset:       latest.GetOffset(),
			Exchange:     latest.GetExchange(),
			Time:         latest.GetTime(),
			CurrencyPair: latest.Pair(),
			AssetType:    latest.GetAssetType(),
			Interval:     latest.GetInterval(),
			Reason:       latest.GetReason(),
		},
		ClosePrice: latest.ClosePrice(),
		HighPrice:  latest.HighPrice(),
		OpenPrice:  latest.OpenPrice(),
		LowPrice:   latest.LowPrice(),
	}, nil
}

func (s *Strategy) SetID(id int) {
	s.ID = id
}

func (s *Strategy) GetPair() currency.Pair {
	return s.pair
}

func (s *Strategy) SetPair(p currency.Pair) {
	s.pair = p
}

func (s *Strategy) SetWeight(d decimal.Decimal) {
	s.weight = d
}

func (s *Strategy) GetWeight() decimal.Decimal {
	return s.weight
}

// UsingSimultaneousProcessing returns whether multiple currencies can be assessed in one go
func (s *Strategy) UsingSimultaneousProcessing() bool {
	return s.useSimultaneousProcessing
}

// SetSimultaneousProcessing sets whether multiple currencies can be assessed in one go
func (s *Strategy) SetSimultaneousProcessing(b bool) {
	s.useSimultaneousProcessing = b
}

// UsingExchangeLevelFunding returns whether funding is based on currency pairs or individual currencies at the exchange level
func (s *Strategy) UsingExchangeLevelFunding() bool {
	return s.usingExchangeLevelFunding
}

// SetExchangeLevelFunding sets whether funding is based on currency pairs or individual currencies at the exchange level
func (s *Strategy) SetExchangeLevelFunding(b bool) {
	s.usingExchangeLevelFunding = b
}

// func (s *Strategy) Direction() order.Side {
// 	return s.direction
// }

func (s *Strategy) SetDirection(direction order.Side) {
	s.direction = direction
}

func (s *Strategy) SelectFeatures() {
	fmt.Println("make api request here")
}

func (s *Strategy) Stop() {
	// fmt.Println("num trades:", len(p.ClosedTrades))
	return
	// for i := range s.indicatorValues {
	// 	x := s.indicatorValues[i]
	// 	fmt.Printf("%d,%s,%s\n", x.Timestamp.Unix(), x.rsiValue, x.maValue)
	// }
}

// func (s *Strategy) GetCurrencySettings() *ExchangeAssetPairSettings {
// 	return s.currencySettings
// }
//
// func (s *Strategy) SetCurrencySettings(e *ExchangeAssetPairSettings) {
// 	s.currencySettings = e
// }
