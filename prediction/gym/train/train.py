import math
import os.path
from typing import Generator

import numpy as np
import pandas as pd
import torch
from matplotlib import pyplot as plt
from torch import nn
from torch.utils.data import DataLoader, ConcatDataset, Dataset
from torch.utils.tensorboard import SummaryWriter
from tqdm import tqdm

from epoch import train_epoch, validate_epoch
from gym.board_logger import BoardLogger
from model.player_predictor import PlayerPredictor
from player_dataset import PlayerDataset
from test import test_model


def format_model_name(n_nearest_players, stack_size, hidden_size, lr, epochs, batch_size, n_parameters = 0):
    if n_nearest_players > 0:
        return f'{n_nearest_players}-{stack_size}-{hidden_size}-{lr}-{epochs}-{batch_size}-{n_parameters}'
    return f'{n_nearest_players}-{stack_size}-{hidden_size}-{lr}-{epochs}-{batch_size}'


export_directory: str
summary_writer: SummaryWriter
dataset_split_ratio: float
device: str
chunk_path: str
player_numbers: list[str]
hyper_parameters: dict

train_set: PlayerDataset
validation_set: PlayerDataset
train_dataloader: DataLoader
validation_dataloader: DataLoader

def init():
    global export_directory
    global summary_writer
    global dataset_split_ratio
    global device
    global chunk_path
    global player_numbers
    global hyper_parameters

    export_directory = os.path.abspath(f'../runs')

if __name__ == '__main__':
    torch.random.seed()
    dataset_split_ratio = 0.8
    batch_size = 64
    n_nearest_players = 3
    device = torch.device("cuda:0" if torch.cuda.is_available() else "cpu")
    chunk_path = os.path.abspath("../data") #use all matches later

    # hyper_parameters = {
    #     'n_nearest_players': range(1, 5 + 1),
    #     'stack_size': range(1, 16 + 1),
    #     'hidden_size': [8, 16, 32, 64, 128, 256, 512],
    #     'sequence_length': [10, 20, 30, 40],
    #     'batch_size': [16, 32, 64, 128, 256],
    #     'lr': [0.0001, 0.00001, 0.1, 0.001],
    # }

    hyper_parameters = {
        'n_nearest_players': [5], #3-5
        'stack_size': [3], #4-32
        'hidden_size': [128], #32-128
        'sequence_length': [20], #20-40
        'batch_size': [64],
        'lr': [0.001],
    }
    datasets = {}

if __name__ == '__init__':
    init()

def train_model(
    n_nearest: int,
    stack_size: int,
    hidden_size: int,
    sequence_length: int,
    batch_size: int,
    lr: float
):
    print(f'n:{n_nearest} stack_size:{stack_size} hidden_size:{hidden_size} seq:{sequence_length} lr:{lr} batch_size:{batch_size}')
    model = PlayerPredictor(device, n_nearest, hidden_size, stack_size)
    model.to(device)
    n_parameters = sum(p.numel() for p in model.parameters() if p.requires_grad)
    epochs = 10
    loss_function = nn.MSELoss()
    optimizer = torch.optim.Adam(model.parameters(), lr=lr)

    training_step = 0
    validation_step = 0
    writer = SummaryWriter(f'{export_directory}/player/{format_model_name( n_nearest, stack_size, hidden_size, lr, epochs, batch_size, n_parameters)}')
    training_logger = BoardLogger(writer)
    validation_logger = BoardLogger(writer)
    train_loss, validation_loss = 0, 0
    validation_loss_low = float(math.inf)
    train_loss_low = float(math.inf)
    train_losses, validation_losses = [], []
    for epoch in range(epochs):
        tl = train_epoch(epoch, epochs, model, train_dataloader, loss_function, optimizer, device, training_logger)
        vl = validate_epoch(epoch, epochs, model, validation_dataloader, loss_function, device, validation_logger)
        if epoch == 0: continue
        train_losses.append(tl)
        validation_losses.append(vl)

        model_name = format_model_name(n_nearest, stack_size, hidden_size, lr, epochs, batch_size, n_parameters)
        model_path = f'{export_directory}/models/{model_name}.pt'

        figure = test_model(model, "../data/f361a535-4d7e-4470-a187-01074c0046fe/chunk_60.csv")
        writer.add_figure(f"Prediction example {model_name}", figure=figure, global_step=epoch)
        writer.flush()
        if vl < validation_loss_low:
            validation_loss_low = vl
            torch.save(model.state_dict(), model_path)
        if tl < train_loss_low:
            train_loss_low = tl

    writer.add_hparams({
        'batch_size': batch_size,
        'n_nearest_players': n_nearest,
        'hidden_size': hidden_size,
        'stack_size': stack_size,
        'sequence_length': sequence_length,
        'lr': lr
    }, {
        'loss': train_loss_low,
        'deviation': validation_loss_low
    })

    writer.flush()
    writer.close()

def load_dataset(n_nearest: int, sequence_length: int, batch_size: int):
    match_directories, chunk_files = os.listdir(chunk_path), []
    for directory in match_directories:
        match_files = os.listdir(chunk_path + '/' + directory)
        for file_name in match_files:
            chunk_files.append(os.path.join(directory, file_name))

    chunk_sizes = [pd.read_csv(f"{chunk_path}/{p}").shape[0] for p in chunk_files]
    total_samples = sum(chunk_sizes)
    average_chunk_size = total_samples // len(chunk_sizes)
    print(f"found {len(chunk_sizes)} chunks of {total_samples} samples with average chunk size {average_chunk_size}")
    global train_set, validation_set, train_dataloader, validation_dataloader
    train_set, validation_set = PlayerDataset.from_dir(chunk_path, n_nearest, sequence_length)
    train_dataloader = DataLoader(train_set, batch_size=batch_size, shuffle=True, num_workers=2)
    validation_dataloader = DataLoader(validation_set, batch_size=batch_size, num_workers=2)


def parameters() -> Generator[tuple[int, int, int, int, int, float], None, None]:
    """n_nearest, stack_size, hidden_size, sequence_length, batch_size, learning_rate"""
    for n_nearest in hyper_parameters['n_nearest_players']:
        for sequence_length in hyper_parameters['sequence_length']:
            for batch_size in hyper_parameters['batch_size']:
                load_dataset(n_nearest, sequence_length, batch_size)
                for stack_size in hyper_parameters['stack_size']:
                    for hidden_size in hyper_parameters['hidden_size']:
                        for lr in hyper_parameters['lr']:
                            yield n_nearest, stack_size, hidden_size, sequence_length, batch_size, lr


if __name__ == '__main__':
    init()
    for params in parameters():
        train_model(*params)