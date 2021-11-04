import pandas as pd
from sklearn.model_selection import train_test_split
from supervised.automl import AutoML
from utils import last_file_in_dir
import json


class MlJarExperiment:
    automl_predictors = dict()

    @classmethod
    def learn(cls, X_train, X_test, y_train, y_test, version=None,
              mode="Explain"):
        if version is not None:
            version = f'../models/{version}'
            automl = AutoML(results_path=version, mode=mode)
        else:
            automl = AutoML(mode=mode)
        automl.fit(X_train, y_train)
        return automl.predict(X_test)

    @classmethod
    def predict(cls, X_test, version):
        if version not in cls.automl_predictors:
            cls.automl_predictors[version] = AutoML(
                results_path=version, mode="Perform")
        return cls.automl_predictors[version].predict(X_test)

    @classmethod
    def drop_features(cls, version):
        return json.load(open(f'../models/{version}/drop_features.json'))
