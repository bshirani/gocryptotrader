import numpy as np
import pandas as pd
from sklearn.ensemble import RandomForestRegressor
from sklearn import datasets
from sklearn.model_selection import train_test_split
from sklearn.feature_selection import SelectFromModel
from sklearn.metrics import accuracy_score
from fastapi import FastAPI
import json

app = FastAPI()
@app.get("/")
async def root():
    return {"hello": "world"}

@app.get("/select_features")
async def select_features(file: str):
    df = pd.read_csv(file, header=0)
    df.set_index('time',inplace=True)

    # TODO get this from the csv header
    feat_labels = ['n10_close', 'n10_low', 'n10_high']
    X = df[feat_labels].values

    # this is whatever the last column of the header row is
    y = df['profit_loss_quote'].values
    X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.4, random_state=0)
    clf = RandomForestRegressor(n_estimators=1, random_state=0, n_jobs=-1)
    clf.fit(X_train, y_train)
    arr = []
    res = dict()
    for feature in zip(feat_labels, clf.feature_importances_):
        res[feature[0]] = feature[1]
    return res

# receives a list of trades and returns the predictions for each
# returns an array of floats
@app.get("/predict")
async def predict(file: str):
    df = pd.read_csv(file, header=0)
    df.set_index('time',inplace=True)

    res = []
    res.append({"id": 1, "pred": 0.52})
    res.append({"id": 2, "pred": 0.93})
    res.append({"id": 3, "pred": 0.19})
    return res
