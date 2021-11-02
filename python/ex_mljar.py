import pandas as pd
from sklearn.model_selection import train_test_split
from supervised.automl import AutoML
from utils import last_file_in_dir


def run_experiment(X_train, X_test, y_train, y_test, version=None,
                   mode="Explain"):
    if version is not None:
        automl = AutoML(results_path=version, mode=mode)
    else:
        automl = AutoML(mode=mode)
    automl.fit(X_train, y_train)
    return automl.predict(X_test)
