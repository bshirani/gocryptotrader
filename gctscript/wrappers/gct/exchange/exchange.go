package exchange

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"time"

	"gocryptotrader/currency"
	"gocryptotrader/engine"
	exchange "gocryptotrader/exchanges"
	"gocryptotrader/exchange/account"
	"gocryptotrader/exchange/asset"
	"gocryptotrader/exchange/kline"
	"gocryptotrader/exchange/order"
	"gocryptotrader/exchange/orderbook"
	"gocryptotrader/exchange/ticker"
	"gocryptotrader/portfolio/banking"
	"gocryptotrader/portfolio/withdraw"
)

// Exchange implements all required methods for Wrapper
type Exchange struct{}

// Exchanges returns slice of all current exchanges
func (e Exchange) Exchanges(enabledOnly bool) []string {
	return engine.Bot.GetExchangeNames(enabledOnly)
}

// GetExchange returns IBotExchange for exchange or error if exchange is not found
func (e Exchange) GetExchange(exch string) (exchange.IBotExchange, error) {
	return engine.Bot.GetExchangeByName(exch)
}

// IsEnabled returns if requested exchange is enabled or disabled
func (e Exchange) IsEnabled(exch string) bool {
	ex, err := e.GetExchange(exch)
	if err != nil {
		return false
	}

	return ex.IsEnabled()
}

// Orderbook returns current orderbook requested exchange, pair and asset
func (e Exchange) Orderbook(ctx context.Context, exch string, pair currency.Pair, item asset.Item) (*orderbook.Base, error) {
	return engine.Bot.GetSpecificOrderbook(ctx, pair, exch, item)
}

// Ticker returns ticker for provided currency pair & asset type
func (e Exchange) Ticker(ctx context.Context, exch string, pair currency.Pair, item asset.Item) (*ticker.Price, error) {
	ex, err := e.GetExchange(exch)
	if err != nil {
		return nil, err
	}

	return ex.FetchTicker(ctx, pair, item)
}

// Pairs returns either all or enabled currency pairs
func (e Exchange) Pairs(exch string, enabledOnly bool, item asset.Item) (*currency.Pairs, error) {
	x, err := engine.Bot.Config.GetExchangeConfig(exch)
	if err != nil {
		return nil, err
	}

	ps, err := x.CurrencyPairs.Get(item)
	if err != nil {
		return nil, err
	}

	if enabledOnly {
		return &ps.Enabled, nil
	}
	return &ps.Available, nil
}

// QueryOrder returns details of a valid exchange order
func (e Exchange) QueryOrder(ctx context.Context, exch, orderID string, pair currency.Pair, assetType asset.Item) (*order.Detail, error) {
	o, err := engine.Bot.OrderManager.GetOrderInfo(ctx, exch, orderID, pair, assetType)
	if err != nil {
		return nil, err
	}

	return &o, nil
}

// SubmitOrder submit new order on exchange
func (e Exchange) SubmitOrder(ctx context.Context, submit *order.Submit) (*order.SubmitResponse, error) {
	r, err := engine.Bot.OrderManager.Submit(ctx, submit)
	if err != nil {
		return nil, err
	}

	return &r.SubmitResponse, nil
}

// CancelOrder wrapper to cancel order on exchange
func (e Exchange) CancelOrder(ctx context.Context, exch, orderID string, cp currency.Pair, a asset.Item) (bool, error) {
	orderDetails, err := e.QueryOrder(ctx, exch, orderID, cp, a)
	if err != nil {
		return false, err
	}

	cancel := &order.Cancel{
		AccountID: orderDetails.AccountID,
		ID:        orderDetails.ID,
		Pair:      orderDetails.Pair,
		Side:      orderDetails.Side,
		AssetType: orderDetails.AssetType,
		Exchange:  exch,
	}

	err = engine.Bot.OrderManager.Cancel(ctx, cancel)
	if err != nil {
		return false, err
	}
	return true, nil
}

// AccountInformation returns account information (balance etc) for requested exchange
func (e Exchange) AccountInformation(ctx context.Context, exch string, assetType asset.Item) (account.Holdings, error) {
	ex, err := e.GetExchange(exch)
	if err != nil {
		return account.Holdings{}, err
	}

	accountInfo, err := ex.FetchAccountInfo(ctx, assetType)
	if err != nil {
		return account.Holdings{}, err
	}

	return accountInfo, nil
}

// DepositAddress gets the address required to deposit funds for currency type
func (e Exchange) DepositAddress(exch string, currencyCode currency.Code) (out string, err error) {
	if currencyCode.IsEmpty() {
		return "", errors.New("currency code is empty")
	}
	return engine.Bot.DepositAddressManager.GetDepositAddressByExchangeAndCurrency(exch, currencyCode)
}

// WithdrawalFiatFunds withdraw funds from exchange to requested fiat source
func (e Exchange) WithdrawalFiatFunds(ctx context.Context, bankAccountID string, request *withdraw.Request) (string, error) {
	ex, err := e.GetExchange(request.Exchange)
	if err != nil {
		return "", err
	}
	var v *banking.Account
	v, err = banking.GetBankAccountByID(bankAccountID)
	if err != nil {
		v, err = ex.GetBase().GetExchangeBankAccounts(bankAccountID, request.Currency.String())
		if err != nil {
			return "", err
		}
	}

	otp, err := engine.Bot.GetExchangeOTPByName(request.Exchange)
	if err == nil {
		otpValue, errParse := strconv.ParseInt(otp, 10, 64)
		if errParse != nil {
			return "", errors.New("failed to generate OTP unable to continue")
		}
		request.OneTimePassword = otpValue
	}
	request.Fiat.Bank.AccountName = v.AccountName
	request.Fiat.Bank.AccountNumber = v.AccountNumber
	request.Fiat.Bank.BankName = v.BankName
	request.Fiat.Bank.BankAddress = v.BankAddress
	request.Fiat.Bank.BankPostalCity = v.BankPostalCity
	request.Fiat.Bank.BankCountry = v.BankCountry
	request.Fiat.Bank.BankPostalCode = v.BankPostalCode
	request.Fiat.Bank.BSBNumber = v.BSBNumber
	request.Fiat.Bank.SWIFTCode = v.SWIFTCode
	request.Fiat.Bank.IBAN = v.IBAN

	resp, err := engine.Bot.WithdrawManager.SubmitWithdrawal(ctx, request)
	if err != nil {
		return "", err
	}
	return resp.Exchange.ID, nil
}

// WithdrawalCryptoFunds withdraw funds from exchange to requested Crypto source
func (e Exchange) WithdrawalCryptoFunds(ctx context.Context, request *withdraw.Request) (string, error) {
	// Checks if exchange is enabled or not so we don't call OTP generation
	_, err := e.GetExchange(request.Exchange)
	if err != nil {
		return "", err
	}
	otp, err := engine.Bot.GetExchangeOTPByName(request.Exchange)
	if err == nil {
		v, errParse := strconv.ParseInt(otp, 10, 64)
		if errParse != nil {
			return "", errors.New("failed to generate OTP unable to continue")
		}
		request.OneTimePassword = v
	}

	resp, err := engine.Bot.WithdrawManager.SubmitWithdrawal(ctx, request)
	if err != nil {
		return "", err
	}
	return resp.Exchange.ID, nil
}

// OHLCV returns open high low close volume candles for requested exchange/pair/asset/start & end time
func (e Exchange) OHLCV(ctx context.Context, exch string, pair currency.Pair, item asset.Item, start, end time.Time, interval kline.Interval) (kline.Item, error) {
	ex, err := e.GetExchange(exch)
	if err != nil {
		return kline.Item{}, err
	}
	ret, err := ex.GetHistoricCandlesExtended(ctx, pair, item, start, end, interval)
	if err != nil {
		return kline.Item{}, err
	}

	sort.Slice(ret.Candles, func(i, j int) bool {
		return ret.Candles[i].Time.Before(ret.Candles[j].Time)
	})

	ret.FormatDates()

	return ret, nil
}
