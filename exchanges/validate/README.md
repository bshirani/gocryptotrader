# GoCryptoTrader package Validate

<img src="/common/gctlogo.png?raw=true" width="350px" height="350px" hspace="70">



This validate package is part of the GoCryptoTrader codebase.

## This is still in active development

You can track ideas, planned features and what's in progress on this Trello board: [https://trello.com/b/ZAhMhpOy/gocryptotrader](https://trello.com/b/ZAhMhpOy/gocryptotrader).

Join our slack to discuss all things related to GoCryptoTrader! [GoCryptoTrader Slack](https://join.slack.com/t/gocryptotrader/shared_invite/enQtNTQ5NDAxMjA2Mjc5LTc5ZDE1ZTNiOGM3ZGMyMmY1NTAxYWZhODE0MWM5N2JlZDk1NDU0YTViYzk4NTk3OTRiMDQzNGQ1YTc4YmRlMTk)

## Current Features for validate

+ This package allows for validation options to occur exchange side e.g.
	- Checking for ID in an order cancellation struct.
	- Determining the correct withdrawal bank details for a specific exchange.

+ Example Usage below:

```go 
// import package
"github.com/thrasher-corp/exchanges/validate"

// define your data structure across potential exchanges
type Critical struct {
	ID string
	Person string
	Banks string
	MoneysUSD float64
}

// define validation and add a variadic param
func (supercritcalinfo *Critical) Validate(opt ...validate.Checker) error {
	// define base level validation
	if supercritcalinfo != nil {
			// oh no this is nil, could panic program!
	}

	// range over potential checks coming from individual packages
	var errs common.Errors
	for _, o := range opt {
		err := o.Check()
		if err != nil {
			errs = append(errs, err)
		}
	}

	if errs != nil {
		return errs
	}
	return nil
}

// define an exchange or package level check that returns a validate.Checker 
// interface
func (supercritcalinfo *Critical) PleaseDontSendMoneyToParents() validate.Checker {
	return validate.Check(func() error {
		if supercritcalinfo.Person == "Mother Dearest" ||
			supercritcalinfo.Person == "Father Dearest" {
			return errors.New("nope")
		}
	return nil
	})
}


// Now in the package all you have to do is add in your options or not...
d := Critical{Person: "Mother Dearest", MoneysUSD: 1337.30}

// This should not error 
err := d.Validate()
if err != nil {
	return err
}

// This should error 
err := d.Validate(d.PleaseDontSendMoneyToParents())
if err != nil {
	return err
}

```


