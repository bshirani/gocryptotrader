import numpy as np
import json
from datetime import datetime
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


def run_all(mode=MODE, retrain=False, dirname=None, version=None):
    """
    runs full machine learning backtest
    runs training during the test
    models and predicts each strategy individually
    predictions are stored back in the trade results
    """

    if dirname is None:
        dirname = last_file_in_dir('../results/fcsv')
    dirname = max(glob.glob('../results/fcsv/*/'), key=os.path.getctime)
    files = glob.glob(dirname+'*')
    [run_one(filename, retrain) for filename in files]


def run_one(filename: str, retrain: bool):
    """
    generates predictions and saves it to the trades

    tries various machine learning methods on the data
    returns the best result and breakdown of the various methods

    backtests a single strategy

    train on first 30
    retrieve predictions for next 30
    retrain and continue until finished
    """

    df = pd.read_csv(filename, header=0)
    df.time = pd.to_datetime(df.time, unit='s')
    target = df.drop(columns=[TARGET_NAME])
    df.profit_loss_quote = df.profit_loss_quote * 0.001
    input_cols = [
        item for item in df.columns.values if item not in IGNORE_COLS
    ]

    split = int(len(df) / 2)
    X_train = df[input_cols].iloc[:split]
    X_test = df[input_cols].iloc[split:]
    y_train = df[TARGET_NAME].iloc[:split]
    y_test = df[TARGET_NAME].iloc[split:]
    version = os.path.basename(filename).split('-')[0]

    if retrain:
        print("learning", len(X_train), "samples")
        preds = MlJarExperiment.learn(
            X_train,
            X_test,
            y_train,
            y_test,
            mode=MODE,
            version=version,
            retrain=retrain)

    # print(X_test.id)
    print("predicting", len(X_test), "samples")
    preds = MlJarExperiment.predict(
        X_test,
        version=version,
        mode=MODE,
    )
    X_test['prediction'] = preds

    x = f'../results/backtest/{version}*'

    tdirname = last_file_in_dir('../results/backtest/*')
    tfilename = last_file_in_dir(tdirname+"/"+version+"*")
    js = json.load(open(tfilename))
    X_test.set_index('id', inplace=True)

    for jr in js:
        # print("checking", jr['ID'])
        if jr['ID'] in X_test.index.values:
            # print("setting prediction for", jr['ID'])
            pred = X_test.loc[jr['ID']].prediction
            jr['Prediction'] = pred

    print("writing filename", tfilename)
    with open(tfilename, 'w') as f:
        json.dump(js, f)

# def run_one_river(mode: str, filename: str):
#     """
#     tries various machine learning methods on the data
#     returns the best result and breakdown of the various methods
#
#     backtests a single strategy
#
#     method 1:
#     train on first 30
#     retrieve predictions for next 30
#     retrain and continue until finished
#
#     method 2:
#     use river and stream the entire thing
#     """
#
#     from river import linear_model
#     from river import datasets
#     from river import compose
#     from river import preprocessing
#     dataset = datasets.Phishing()
#     # print(dataset)
#
#     model = compose.Pipeline(
#         preprocessing.StandardScaler(),
#         linear_model.LogisticRegression()
#     )
#
#     df = pd.read_csv(filename, header=0)
#     # df['pl_cheat'] = df[TARGET_NAME]
#     df.time = pd.to_datetime(df.time, unit='s')
#     target = df.drop(columns=[TARGET_NAME])
#     df.profit_loss_quote = df.profit_loss_quote * 0.001
#     input_cols = [
#         item for item in df.columns.values if item not in IGNORE_COLS
#     ]
#     #
#     # X_train, X_test, y_train, y_test = train_test_split(
#     #     df[input_cols], df[TARGET_NAME], test_size=0.25
#     # )
#     # version = os.path.basename(filename).split('-')[0]
#     #
#     df.drop(columns=['time', 'id'], inplace=True)
#     # print(df.columns)
#     # for x, y in df.iterrows():
#     #     x = y.to_dict()
#     #     y = x['profit_loss_quote']
#     #     del x['profit_loss_quote']
#     #     # print("learning", type(x), type(y))
#     #     y_pred = model.predict_proba_one(x)
#     #     print("learning", x)
#     #     model.learn_one(x, y)
#     # print("y pred", y_pred)
#     preds = MlJarExperiment.learn(
#         X_train,
#         X_test,
#         y_train,
#         y_test,
#         mode=mode,
#         version=version,
#         retrain=True)
#     #
#     # for i in range(100):
#     #     t1 = datetime.now()
#     #     p = MlJarExperiment.predict(X_test.iloc[:1], version=version)
#
#     # import pdbr
#     # pdbr.set_trace()
#     # analyze_model_performance(df, preds, X_test, y_test)


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
    fire.Fire(run_all)
