import numpy as np


def sortino_ratio(series, N, rf):
    mean = series.mean() * N - rf
    std_neg = series[series < 0].std()*np.sqrt(N)
    return mean/std_neg
