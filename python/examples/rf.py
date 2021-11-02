import numpy as np
import pandas as pd
from sklearn.ensemble import RandomForestRegressor
from sklearn import datasets
from sklearn.model_selection import train_test_split
from sklearn.feature_selection import SelectFromModel
from sklearn.metrics import accuracy_score
from fastapi import FastAPI

app = FastAPI()

@app.get("/")
async def root():

    df = pd.read_csv('data/btc.csv', header=0)
    df.set_index('time',inplace=True)
    feat_labels = ['n10_close', 'n10_low', 'n10_high']
    X = df[feat_labels].values
    y = df['profit_loss_quote'].values
    X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.4, random_state=0)
    clf = RandomForestRegressor(n_estimators=2000, random_state=0, n_jobs=-1)
    clf.fit(X_train, y_train)
    for feature in zip(feat_labels, clf.feature_importances_):
        print(feature)
    return {"message": "Hello World"}
