import torch
import flash
from flash.core.data.utils import download_data
from flash.tabular import TabularClassificationData, TabularClassifier

train_file="data/btc.csv"
datamodule = TabularClassificationData.from_csv(
    ["time"],
    ["open","low","high","n10_high","n10_low","n10_close"],
    target_fields="profit_loss_quote",
    train_file=train_file,
    val_split=0.1,
)
model = TabularClassifier.from_data(datamodule)
trainer = flash.Trainer(max_epochs=3, gpus=torch.cuda.device_count())
trainer.fit(model, datamodule=datamodule)
predictions = model.predict(train_file)
print(predictions)
trainer.save_checkpoint("tabular_classification_model.pt")
