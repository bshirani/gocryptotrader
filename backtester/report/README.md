# GoCryptoTrader Backtester: Report package

<img src="/backtester/common/backtester.png?raw=true" width="350px" height="350px" hspace="70">



This report package is part of the GoCryptoTrader codebase.

## This is still in active development

You can track ideas, planned features and what's in progress on this Trello board: [https://trello.com/b/ZAhMhpOy/gocryptotrader](https://trello.com/b/ZAhMhpOy/gocryptotrader).

Join our slack to discuss all things related to GoCryptoTrader! [GoCryptoTrader Slack](https://join.slack.com/t/gocryptotrader/shared_invite/enQtNTQ5NDAxMjA2Mjc5LTc5ZDE1ZTNiOGM3ZGMyMmY1NTAxYWZhODE0MWM5N2JlZDk1NDU0YTViYzk4NTk3OTRiMDQzNGQ1YTc4YmRlMTk)

## Report package overview

The report package helps generates the output under the `results` folder.

As the application is run, many statistics such as purchase events are tracked. These events are utilised and enhanced in the report package in order to render an HTML report for easy comparison and historical strategy effectiveness.

The report utilises the following sweet technologies:
- go templating ([tpl.gohtml](tpl.gohtml))
- [mdbootstrap](https://mdbootstrap.com/)
- [lightweightcharts](https://github.com/tradingview/lightweight-charts/) by [TradingView](https://www.tradingview.com/)

Output example:
![example](https://user-images.githubusercontent.com/9261323/105283038-c124be00-5c03-11eb-88af-d67e727a8c16.png)



