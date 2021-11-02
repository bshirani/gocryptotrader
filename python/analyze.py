from sklearn.metrics import mean_squared_error
from metrics import sortino_ratio
import pandas as pd
from analyze_trades import analyze_trades

ACCOUNT_BALANCE = 1000.0


class ModelAnalysis:
    def __init__(self, name, preds):
        self.name = name
        self.mpreds = preds

    def analyze(self):
        df = self.mpreds.copy()
        df['net_profit'] = df.profit_loss_quote

        # import pdbr
        # pdbr.set_trace()
        df['time'] = pd.to_datetime(df['time'], unit='s')
        df.set_index('time', inplace=True)
        adf = analyze_trades(df)
        adf['name'] = self.name
        return adf

        # self.ret = (self.mpreds.profit_loss_quote / ACCOUNT_BALANCE).sum() * 100
        # self.net_profit = self.mpreds.profit_loss_quote.sum()
        # self.sortino_ratio_trades = sortino_ratio(
        #     self.mpreds.profit_loss_quote, len(self.mpreds), 0.01)
        # self.sortino_ratio_180 = sortino_ratio(
        #     self.mpreds.profit_loss_quote, 180, 0.01)
        # self.num_trades = len(self.mpreds)
        #
        # df = self.mpreds.copy()
        # df['Date'] = pd.to_datetime(
        #     df['time']) - pd.to_timedelta(7, unit='d')
        #
        # df = df.groupby([pd.Grouper(key='Date', freq='W-MON')]
        #                 )['profit_loss_quote'].sum().reset_index().sort_values('Date')
        # self.percentage_of_weeks_win = len(
        #     df[df.profit_loss_quote > 0]) / len(df)
    # def to_series(self):
    #     return pd.Series([
    #         self.name,
    #         # self.ret,
    #         self.net_profit,
    #         self.percentage_wins,
    #         self.sortino_ratio_trades,
    #         self.sortino_ratio_180,
    #         self.num_trades,
    #         self.percentage_of_weeks_win,
    #     ])


class ModelAnalyzer:
    rows = [
        # 'sum_return',
        'net_profit',
        'perc_winners',
        'profit_factor',
        'expectancy',
        'num_trades',
        'gross_profit',
        'gross_losses',
        'trades_per_year',
        'perc_months_profitable',
        'perc_weeks_profitable',
        'max_consequtive_losers',
        'max_consequtive_winners',
        'max_drawdown',
        # 'sortino_trades',
        # 'sortino_180',
        # 'num_trades',
        # '%weekswon'
    ]

    def __init__(self, preds):
        self.preds = preds
        pap = preds.copy()
        # import pdbr
        # pdbr.set_trace()
        pap.profit_loss_quote = pap.prediction * pap.profit_loss_quote
        self.analyses = [
            ModelAnalysis('m1', preds),
            ModelAnalysis('no_model', preds),
            ModelAnalysis('over_0', preds[preds.prediction > 0]),
            ModelAnalysis('over_0_adj_profit', pap[pap.prediction > 0]),
            ModelAnalysis('over_.5', preds[preds.prediction > 0.5]),
            ModelAnalysis('over_.1', preds[preds.prediction > 0.1]),
            ModelAnalysis('over_.5_adj_profit', pap[pap.prediction > 0.5]),
        ]

    def analyze(self):
        print("prediction range", self.preds.prediction.max(),
              self.preds.prediction.min())

        cols = [x.name for x in self.analyses]
        df = [x.analyze() for x in self.analyses]
        df = pd.concat(df, axis=0)
        df.set_index('name', inplace=True)
        # df.columns = df.iloc[0].values
        # df = df.iloc[1:]
        # df.index = ModelAnalyzer.rows
        return df[ModelAnalyzer.rows].T
