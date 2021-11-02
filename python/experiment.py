import numpy as np
import fire
import pandas as pd
import functools
from math import sqrt

from utils import last_file_in_dir
from ex_mljar import run_experiment as experiment_mljar
from sklearn.model_selection import train_test_split

from analyze import ModelAnalyzer

TARGET_NAME = 'profit_loss_quote'
# INPUT_COLS = ['n10_high', 'n10_low']
# INPUT_COLS = ['pl_cheat']
# INPUT_COLS = ['pl_cheat', 'n10_high', 'n10_low']
IGNORE_COLS = [
    'profit_loss_quote',
    'risk_quote',
    'id',
    'time',
]
# INPUT_COLS = [
#     'n10_highrel',
#     'n10_lowrel',
#     'n10_openrel',
#     'n10_pctchg',
#     'n10_sloperel'
# ]
INPUT_COLS = df.columns - IGNORE_COLS

# Target columns
# percentage of range captured
# profit_loss_quote
# profit_loss / risked


def experiment(test=False):
    """
    tries various machine learning methods on the data
    returns the best result and breakdown of the various methods
    """

    filename = last_file_in_dir('../results/fcsv/*BUY*')
    df = pd.read_csv(filename, header=0)
    df['pl_cheat'] = df[TARGET_NAME]
    df.time = pd.to_datetime(df.time, unit='s')
    target = df.drop(columns=[TARGET_NAME])
    df.profit_loss_quote = df.profit_loss_quote * 0.001

    X_train, X_test, y_train, y_test = train_test_split(
        df[INPUT_COLS], df[TARGET_NAME], test_size=0.25
    )

    print("running experiment on", filename)
    if test:
        preds = np.random.uniform(low=-5, high=5, size=(len(X_test,)))
    else:
        preds = experiment_mljar(X_train, X_test, y_train, y_test)

    analyze_model_performance(df, preds, X_test, y_test)


def analyze_model_performance(df, preds, X_test, y_test):
    preds = pd.DataFrame(preds)
    preds.columns = ['prediction']
    preds.index = y_test.index
    preds = preds.join(df)
    preds['error'] = preds.profit_loss_quote - preds.prediction
    preds['error2'] = preds.error**2
    res = ModelAnalyzer(preds).analyze()
    print(res.round(2))


if __name__ == '__main__':
    fire.Fire(experiment)
