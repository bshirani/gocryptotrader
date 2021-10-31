import os
import pytorch_lightning as pl
import pandas as pd
from torchvision import datasets, transforms

import torch
from pytorch_lightning import LightningModule, Trainer
from torchmetrics.functional.classification.accuracy import accuracy
from torch import nn
from torch.nn import functional as F
from torch.utils.data import DataLoader, random_split
from torchvision import transforms
import torch.utils.data as data_utils

# targets_data = [random.random() for i in range(10)]
# targets_df = pd.DataFrame(data=targets_data)
# targets_df.columns = ['targets']
# torch_tensor = torch.tensor(targets_df['targets'].values)

PATH_DATASETS = os.environ.get("PATH_DATASETS", ".")
AVAIL_GPUS = min(1, torch.cuda.device_count())
BATCH_SIZE = 256 if AVAIL_GPUS else 64

filename = "../results/fcsv/2021-10-30-19-02-17-trend@BTC_USDT@BUY.csv"

df = pd.read_csv(filename, header=0)
x = df['n10_high'].values
y = df['profit_loss_quote'].values
print(x,y)

target = torch.tensor(df['profit_loss_quote'].values)
features = torch.tensor(df.drop('profit_loss_quote', axis = 1).values)
train = data_utils.TensorDataset(features, target)
train_loader = data_utils.DataLoader(train, batch_size=10, shuffle=True)

# train = data_utils.TensorDataset(x, y)
# train_loader = data_utils.DataLoader(train, batch_size=50, shuffle=True)

# inputs = [[ 1,  2,  3,  4,  5],[ 2,  3,  4,  5,  6]]
# targets = [ 6,7]
# batch_size = 2
# inputs  = torch.tensor(inputs)
# targets = torch.IntTensor(targets)
# dataset =TensorDataset(inputs, targets)
# data_loader = DataLoader(dataset, batch_size, shuffle = True)


class TradeDataModule(pl.LightningDataModule):
    def __init__(self):
        super().__init__()

        self.df = pd.read_csv(filename, header=0)
        self.download_dir = ''
        self.batch_size = 2
        self.transform = transforms.Compose([
            transforms.ToTensor()
        ])

    def __len__(self):
        return len(self.df)

    def __getitem__(self, index):
        print("GET ITEM")
        row = self.df.loc[index]
        return (
            torchvision.transforms.functional.to_tensor(row.values),
            row["profit_loss_quote"],
        )

    def prepare_data(self):
        pass

    def setup(self, stage=None):
        data = datasets.MNIST(self.download_dir,
                              train = True,
                              transform = self.transform)
        print(len(self.df))
        self.train_data, self.valid_data = random_split(self.df.values,
                                                        [10, 4])
        # self.train_data, self.valid_data = random_split(data,
        #                                                 [55000, 5000])
        self.test_data = datasets.MNIST(self.download_dir,
                                        train = False,
                                        transform = self.transform)
        print(len(self.test_data))
        print(type(self.test_data))

    def train_dataloader(self):
        return DataLoader(self.train_data, batch_size=BATCH_SIZE,
                          num_workers=12)

    def val_dataloader(self):
        return DataLoader(self.valid_data, batch_size=BATCH_SIZE, num_workers=12)

    def test_dataloader(self):
        return DataLoader(self.test_data, batch_size=BATCH_SIZE, num_workers=12)

class StrategyModel(LightningModule):
    def __init__(self, data_dir=PATH_DATASETS, hidden_size=64, learning_rate=2e-4):
        super().__init__()

        # Set our init args as class attributes
        self.data_dir = data_dir
        self.hidden_size = hidden_size
        self.learning_rate = learning_rate

        # # Hardcode some dataset specific attributes
        # self.num_classes = 2
        self.dims = (1, 28, 28)
        channels, width, height = self.dims
        # print("channels", channels, "width", width, "height", height)
        # self.transform = transforms.Compose(
        #     [
        #         transforms.ToTensor(),
        #         transforms.Normalize((0.1307,), (0.3081,)),
        #     ]
        # )

        # Define PyTorch model
        self.model = nn.Sequential(
            nn.Flatten(),
            nn.Linear(channels * width * height, hidden_size),
            nn.ReLU(),
            nn.Dropout(0.1),
            nn.Linear(hidden_size, hidden_size),
            nn.ReLU(),
            nn.Dropout(0.1),
            nn.Linear(hidden_size, self.num_classes),
        )

    def forward(self, x):
        x = self.model(x)
        return F.log_softmax(x, dim=1)

    def training_step(self, batch, batch_idx):
        x, y = batch
        logits = self(x)
        loss = F.nll_loss(logits, y)
        return loss

    def validation_step(self, batch, batch_idx):
        x, y = batch
        logits = self(x)
        loss = F.nll_loss(logits, y)
        preds = torch.argmax(logits, dim=1)
        acc = accuracy(preds, y)

        # Calling self.log will surface up scalars for you in TensorBoard
        self.log("val_loss", loss, prog_bar=True)
        self.log("val_acc", acc, prog_bar=True)
        return loss

    def test_step(self, batch, batch_idx):
        # Here we just reuse the validation_step for testing
        return self.validation_step(batch, batch_idx)

    def configure_optimizers(self):
        optimizer = torch.optim.Adam(self.parameters(), lr=self.learning_rate)
        return optimizer


# tradeData = TradeDataModule()
model = StrategyModel()
trainer = Trainer(
    gpus=1,
    max_epochs=3,
    progress_bar_refresh_rate=20,
)
trainer.fit(model, train_loader)

# # file =  "../results/fcsv/2021-10-30-16-07-54-BTC_USDT-trend@BTC_USDT@BUY.csv"
# file = "../results/fcsv/2021-10-30-19-02-17-trend@BTC_USDT@BUY.csv"
#
# def run():
#     df = pd.read_csv(file, header=0)
#     df.set_index('time',inplace=True)
#
#     # create test and train set
#
#
#     res = []
#     res.append({"id": 1, "pred": 0.52})
#     res.append({"id": 2, "pred": 0.93})
#     res.append({"id": 3, "pred": 0.19})
#     print(res)
#     return res
#
# run()
