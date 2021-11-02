from pytorch_tabular.models.tab_transformer.config import TabTransformerConfig
from pytorch_tabular.models.ft_transformer.config import FTTransformerConfig
import torch
import numpy as np
from torch.functional import norm
from sklearn.datasets import fetch_covtype
from pytorch_tabular.models import AutoIntModel, AutoIntConfig
from pytorch_tabular.config import (
    DataConfig,
    ExperimentConfig,
    ExperimentRunManager,
    ModelConfig,
    OptimizerConfig,
    TrainerConfig,
)
from pytorch_tabular.models.node.config import NodeConfig
from pytorch_tabular.models.category_embedding.config import (
    CategoryEmbeddingModelConfig,
)
from pytorch_tabular.models.category_embedding.category_embedding_model import (
    CategoryEmbeddingModel,
)
import pandas as pd
from omegaconf import OmegaConf
from pytorch_tabular.tabular_datamodule import TabularDatamodule
from pytorch_tabular.tabular_model import TabularModel
import pytorch_lightning as pl
from sklearn.preprocessing import PowerTransformer
from pathlib import Path
import wget
from pytorch_tabular.utils import get_balanced_sampler, get_class_weighted_cross_entropy
from sklearn.preprocessing import PowerTransformer


def experiment_pytorch_tabular(df: pd.DataFrame):
    data_config = DataConfig(
        target=target_name,
        continuous_cols=num_col_names,
        categorical_cols=cat_col_names,
        continuous_feature_transform=None,  # "quantile_normal",
        normalize_continuous_features=True,
        num_workers=12,
    )

    # model_config = CategoryEmbeddingModelConfig(task="regression")
    model_config = AutoIntConfig(
        task="regression",
        deep_layers=True,
        embedding_dropout=0.2,
        batch_norm_continuous_input=True,
        attention_pooling=True,
    )

    trainer_config = TrainerConfig(
        checkpoints=None,
        max_epochs=10,
        gpus=1,
        profiler=None,
        fast_dev_run=True,
        auto_lr_find=False,
    )
    optimizer_config = OptimizerConfig()
    tabular_model = TabularModel(
        data_config=data_config,
        model_config=model_config,
        optimizer_config=optimizer_config,
        trainer_config=trainer_config,
    )
    tr = PowerTransformer()

    def fake_metric(y_hat, y):
        return (y_hat - y).mean()

    tabular_model.fit(
        train=train,
        test=test,
        metrics=[fake_metric],
        target_transform=tr,
        loss=torch.nn.L1Loss(),
        optimizer=torch.optim.Adagrad,
        optimizer_params={},
    )
    # result = tabular_model.evaluate(test)
    result = tabular_model.predict(test)
    print(result.head())
    # print(result[['pl2', 'profit_loss_quote_prediction']].head())

# sampler = get_balanced_sampler(train[target_name].values.ravel())
# # cust_loss = get_class_weighted_cross_entropy(train[target_name].values.ravel())
# tabular_model.fit(
#     train=train,
#     validation=val,
#     # loss=cust_loss,
#     train_sampler=sampler)
#
# from pytorch_tabular.categorical_encoders import CategoricalEmbeddingTransformer
# transformer = CategoricalEmbeddingTransformer(tabular_model)
# train_transform = transformer.fit_transform(train)
# # test_transform = transformer.transform(test)
# # ft = tabular_model.model.feature_importance()
# # result = tabular_model.evaluate(test)
# # print(result)
# # test.drop(columns=ta6rget_name, inplace=True)
# # pred_df = tabular_model.predict(test)
# # print(pred_df.head())
# # pred_df.to_csv("output/temp2.csv")
# # tabular_model.save_model("test_save")
# # new_model = TabularModel.load_from_checkpoint("test_save")
# # result = new_model.evaluate(test)
# # print(result)
