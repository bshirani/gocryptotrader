from river import metrics
from river import linear_model
from river import datasets

dataset = datasets.Phishing()
model = linear_model.LogisticRegression()

metric = metrics.ROCAUC()

for x, y in dataset:
    y_pred = model.predict_proba_one(x)
    model.learn_one(x, y)
    metric.update(y, y_pred)

print(metric)
