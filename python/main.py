import numpy as np
from pydantic import BaseModel, EmailStr
import urllib
import os
import pandas as pd
import sys
import asyncio
import supervised
from sklearn.ensemble import RandomForestRegressor
from sklearn import datasets
from sklearn.model_selection import train_test_split
from sklearn.feature_selection import SelectFromModel
from sklearn.metrics import accuracy_score
from fastapi import FastAPI, Request, Response, HTTPException
from ex_mljar import MlJarExperiment
from utils import last_file_in_dir
import json

app = FastAPI()
MODE = "Explain"


@app.get("/")
async def root():
    return {"hello": "world"}


@app.get("/select_features")
async def select_features(file: str):
    # if file == "":
    #     file = last_file_in_dir(
    df = pd.read_csv(file, header=0)
    df.set_index('time', inplace=True)

    # TODO get this from the csv header
    feat_labels = ['n10_close', 'n10_low', 'n10_high']
    X = df[feat_labels].values

    # this is whatever the last column of the header row is
    y = df['profit_loss_quote'].values
    X_train, X_test, y_train, y_test = train_test_split(
        X, y, test_size=0.4, random_state=0)
    clf = RandomForestRegressor(n_estimators=1, random_state=0, n_jobs=-1)
    clf.fit(X_train, y_train)
    arr = []
    res = dict()
    for feature in zip(feat_labels, clf.feature_importances_):
        res[feature[0]] = feature[1]
    return res

# receives a list of trades and returns the predictions for each
# returns an array of floats


@app.get("/drop_features")
async def drop_features(model: str):
    return MlJarExperiment.drop_features(model)


@app.post("/learn")
async def learn(model: str, request: Request):
    j = await request.json()
    df = pd.json_normalize(j)
    # rd = dict(req.query_params)
    # del rd['model']
    df.set_index('Time', inplace=True)
    df.drop(columns='Date', inplace=True)
    keep_cols = [x for x in df.columns if "time" not in x]
    df = df[keep_cols]
    df = df.astype(float)
    X_train, X_test, y_train, y_test = train_test_split(
        df[keep_cols], df['ProfitLossQuote'], test_size=0.25
    )
    print(df.columns)

    # import pdbr
    # pdbr.set_trace()

    try:
        MlJarExperiment.learn(
            X_train,
            X_test,
            y_train,
            y_test,
            MODE,
            f"{model}",
            retrain=True,
        )
        return "ok"
    except supervised.exceptions.AutoMLException:
        raise HTTPException(status_code=500, detail="AutoML Exception")


@app.post("/predict")
async def predict(model: str, request: Request):
    j = await request.json()
    df = pd.json_normalize(j)
    # rd = dict(req.query_params)
    # del rd['model']
    df.set_index('Time', inplace=True)
    df.drop(columns='Date', inplace=True)
    keep_cols = [x for x in df.columns if "time" not in x]
    df = df[keep_cols]
    df = df.astype(float)
    print(df.columns)

    # odf = pd.read_csv(filename, header=0)
    # odf.set_index('time', inplace=True)
    # oX_test = odf[-1:].iloc[0]
    # X_test = pd.DataFrame(rd, index=[0]).astype(float)
    # print(X_test.to_dict())
    # print(urllib.parse.urlencode(rd, doseq=False))
    # X_test = pd.DataFrame(rd, index=[0])
    # import pdbr
    # pdbr.set_trace()

    try:
        res = MlJarExperiment.predict(
            df, f"{model}", mode=MODE).tolist()[0]
        print("PREDICTION:", res)
        if res == float('inf'):
            print("INFINITYYYYYYYYYY")
            res = 1
        return res
    except supervised.exceptions.AutoMLException:
        raise HTTPException(status_code=500, detail="AutoML Exception")
        # return f"model not found {model}"


if __name__ == "__main__":
    print(asyncio.run(predict(model=sys.argv[1])))
