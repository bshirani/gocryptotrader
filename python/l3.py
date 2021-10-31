import torch

import flash
from flash.core.data.utils import download_data
from flash.tabular import TabularClassificationData, TabularClassifier

# 1. Create the DataModule
download_data("https://pl-flash-data.s3.amazonaws.com/titanic.zip", "./data")

datamodule = TabularClassificationData.from_csv(
    ["Sex", "Age", "SibSp", "Parch", "Ticket", "Cabin", "Embarked"],
    "Fare",
    target_fields="Survived",
    train_file="data/titanic/titanic.csv",
    val_split=0.1,
)

# 2. Build the task
model = TabularClassifier.from_data(datamodule)

# 3. Create the trainer and train the model
trainer = flash.Trainer(max_epochs=3, gpus=torch.cuda.device_count())
trainer.fit(model, datamodule=datamodule)

# 4. Generate predictions from a CSV
predictions = model.predict("data/titanic/titanic.csv")
print(predictions)

# 5. Save the model!
trainer.save_checkpoint("tabular_classification_model.pt")
