import pandas as pd
import numpy as np


def analyze_trades(ts, portfolio=False, initial_balance=10000):
    if "units" not in ts:
        ts.loc[:, "units"] = 1
    if "unit_profit" not in ts:
        ts.loc[:, "unit_profit"] = ts.net_profit

    # remove trades with a 0 unit value
    ts = ts[~(ts.units == 0)].copy()

    if (
        isinstance(ts.index, pd.RangeIndex)
        or isinstance(ts.index, pd.Int64Index)
    ) and "entry_time" in ts.columns:
        ts.set_index("entry_time", inplace=True)

    np = net_profit(ts)
    md = max_drawdown(ts)

    if len(ts) == 0:
        return None

    if md == 0:
        md = 0.01

    sortino_trades = sortino_ratio(ts, len(ts), 0.01)
    sortino_180 = sortino_ratio(ts, 180, 0.01)

    perf = {
        "initial_balance": initial_balance,
        "ending_balance": initial_balance,
        "annual_return": 0 if np == 0 else 100000/np,
        "cum_returns_final": 1,
        "annual_volatility": 1,
        "calmar_ratio": 1,
        "stability_of_timeseries": 1,
        "max_drawdown": md,
        "omega_ratio": 1,
        "sortino_trades": sortino_trades,
        "sortino_180": sortino_180,
        "skew": 1,
        "kurtosis": 1,
        "tail_ratio": 1,
        "common_sense_ratio": 1,
        "value_at_risk": 1,
        "alpha": 1,
        "beta": 1,
        "stability": 1,
        "net_profit": np,
        "recovery_factor": float(np) / md,
        "profit_factor": profit_factor(ts),
        "expectancy": expectancy(ts),
        "num_trades": int(len(ts)),
        "gross_profit": gross_wins(ts),
        "trades_per_year": trades_per_year(ts),
        "max_consequtive_losers": max_consequtive_losers(ts),
        "gross_losses": gross_losses(ts),
        "avg_winner": average_winner(ts),
        "avg_loser": average_loser(ts),
        "perc_winners": perc_profitable(ts),
        "num_wins": int((ts["net_profit"] > 0).sum()),
        "num_losses": int((ts["net_profit"] < 0).sum()),
        "perc_months_profitable": perc_months_profitable(ts),
        "perc_weeks_profitable": perc_weeks_profitable(ts),
        "max_drawdown": md,
        "max_consequtive_winners": max_consequtive_winners(ts),
        "max_win": ts.net_profit.max(),
        "max_loss": ts.net_profit.min(),
        "sharpe_ratio": 1.33333333312312312,
    }

    if portfolio:
        perf["num_trade_rules"] = 5
        perf["num_trade_rule_targets"] = 100
        perf["num_instruments"] = 100

    if "prediction" in ts:
        perf["target_false_negative"] = target_false_negative(ts),
        perf["trade_false_negative"] = trade_false_negative(ts),
        perf["preds_mean"] = ts["prediction"].mean()
        perf["preds_std"] = ts["prediction"].std()
        perf["preds_max"] = ts["prediction"].max()
        perf["preds_min"] = ts["prediction"].min()

    return pd.DataFrame(perf, index=[0])


def net_profit(ts):
    return ts["net_profit"].sum()


def profit_factor(ts):
    wins = ts[ts["net_profit"] >= 0]["net_profit"].sum()
    losses = ts[ts["net_profit"] < 0]["net_profit"].sum()
    if wins == 0 or losses == 0:
        return 0
    return wins / (-1 * losses)


def average_winner(ts):
    winners = ts["net_profit"][ts["net_profit"] >= 0]
    if len(winners) == 0:
        return 0
    if winners.sum() == 0:
        return 0
    return winners.sum() / len(winners)


def average_loser(ts):
    losers = ts["net_profit"][ts["net_profit"] < 0]
    if len(losers) == 0:
        return 0
    if losers.sum() == 0:
        return 0
    return losers.sum() / len(losers)


def perc_profitable(ts):
    profitable = len(ts[ts["net_profit"] >= 0])
    not_profitable = len(ts[ts["net_profit"] < 0])
    if (profitable + not_profitable) < 1:
        return 0
    return profitable / (profitable + not_profitable)


def expectancy(ts):
    win_rate = perc_profitable(ts)
    return (win_rate * average_winner(ts)) - (
        (1 - win_rate) * average_loser(ts) * -1
    )


def gross_wins(ts):
    return ts[ts["net_profit"] >= 0]["net_profit"].sum()


def gross_losses(ts):
    return ts[ts["net_profit"] < 0]["net_profit"].sum()


def target_false_negative(ts):
    return 50


def trade_false_negative(ts):
    return 50


def perc_months_profitable(ts):
    pm = ts.groupby([ts.index.year, ts.index.month])["net_profit"].sum()
    return (pm > 0).sum() / len(pm)


def perc_weeks_profitable(ts):
    ts['Date'] = pd.to_datetime(ts.index) - pd.to_timedelta(7, unit='d')
    ts = ts.groupby([pd.Grouper(key='Date', freq='W-MON')]
                    )['profit_loss_quote'].sum().reset_index().sort_values('Date')
    return len(ts[ts.profit_loss_quote > 0]) / len(ts)


def trades_per_year(ts):
    total_days = (ts.index.max() - ts.index.min()).days + 1
    return round(len(ts) / (total_days / 365))


def max_consequtive_losers(ts):
    y = ts["net_profit"] < 0
    return int((y * y.groupby((y != y.shift()).cumsum()).cumcount()).max())


def max_consequtive_winners(ts):
    y = ts["net_profit"] > 0
    return int((y * y.groupby((y != y.shift()).cumsum()).cumcount()).max())


def max_drawdown(ts):
    ts["cumsum"] = ts["net_profit"].cumsum()
    ts["cummax"] = ts["cumsum"].cummax()
    ts["drawdown"] = ts["cummax"] - ts["cumsum"]
    return ts["drawdown"].max()


def sortino_ratio(ts, N, rf):
    series = ts.net_profit
    mean = series.mean() * N - rf
    std_neg = series[series < 0].std()*np.sqrt(N)
    return mean/std_neg
