import numpy as np
import re
import json
from datetime import datetime
from analyze_trades import analyze_trades
import fire
import pandas as pd
import functools
from math import sqrt
import glob
import os

from utils import last_file_in_dir
from ex_mljar import MlJarExperiment
from sklearn.model_selection import train_test_split

from analyze import ModelAnalyzer

MODE = "Explain"
TARGET_NAME = 'profit_loss_quote'
# INPUT_COLS = ['n10_high', 'n10_low']
# INPUT_COLS = ['pl_cheat']
# INPUT_COLS = ['pl_cheat', 'n10_high', 'n10_low']
IGNORE_COLS = [
    'profit_loss_quote',
    'risked_quote',
    # 'id',
    'time',
]
# INPUT_COLS = [
#     'n10_highrel',
#     'n10_lowrel',
#     'n10_openrel',
#     'n10_pctchg',
#     'n10_sloperel'
# ]

# Target columns
# percentage of range captured
# profit_loss_quote
# profit_loss / risked


def pfbacktest():
    """
    runs full portfolio backtest
    applies weights
    this is the final report
    """
    tdirname = last_file_in_dir('../results/backtest/*')
    all_trades = None
    for d in os.listdir(tdirname):
        df = pd.read_json(open(os.path.join(tdirname, d)))
        if all_trades is not None:
            all_trades = pd.concat([all_trades, df])
        else:
            all_trades = df
    all_trades.columns = [re.sub(r'(?<!^)(?=[A-Z])', '_', c).lower()
                          for c in all_trades.columns]
    all_trades['entry_time'] = pd.to_datetime(all_trades.entry_time)
    # all_trades.set_index(pd.to_datetime(all_trades['entry_time']), inplace=True)
    df = analyze_trades(all_trades)
    print(df[sorted(df)].T)


if __name__ == '__main__':
    fire.Fire(pfbacktest)
