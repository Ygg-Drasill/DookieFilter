import math
import os.path

import numpy as np
import pandas as pd
import torch
from matplotlib import pyplot as plt
from torch import nn

from torch.utils.data import DataLoader, ConcatDataset
from dataloader import MatchDataset
from model import PlayerPredictor
from player_dataset import PlayerDataset
from epoch import train_epoch, validate_epoch

if __name__ == '__main__':
    torch.random.seed()
    dataset_split_ratio = 0.8
    batch_size = 64
    n_nearest_players = 3
    device = torch.device("cuda:0" if torch.cuda.is_available() else "cpu")
    chunk_path = os.path.abspath("../data/test") #use all matches later

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
        'stack_size': [4], #4-32
        'hidden_size': [64], #32-128
        'sequence_length': [20], #20-40
        'batch_size': [64],
        'lr': [0.0001], #0.0001-0.00001
    }

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
    chunk_files = os.listdir(chunk_path)
    chunk_sizes = [pd.read_csv(f"{chunk_path}/{p}").shape[0] for p in chunk_files]
    total_samples = sum(chunk_sizes)
    average_chunk_size = total_samples // len(chunk_sizes)
    print(f"found {len(chunk_sizes)} chunks of {total_samples} samples with average chunk size {average_chunk_size}")

    player_dataset = ConcatDataset(
        [PlayerDataset(f"{chunk_path}/{path}", sequence_length, n_nearest) for path in tqdm(os.listdir(chunk_path),
                                                                                            unit='chunk',
                                                                                            desc='building dataset from chunks',
                                                                                            ncols=200)])
    train_size = math.floor(len(player_dataset) * dataset_split_ratio)
    validation_size = len(player_dataset) - train_size
    train_set, validation_set = torch.utils.data.random_split(player_dataset,[train_size, validation_size])
    train_dataloader = DataLoader(train_set, batch_size=batch_size, shuffle=True, num_workers=1)
    validation_dataloader = DataLoader(validation_set, batch_size=batch_size, shuffle=True, num_workers=1)

    model = PlayerPredictor(device, n_nearest_players, 32, 4)
    model.to(device)
    n_parameters = sum(p.numel() for p in model.parameters() if p.requires_grad)
    epochs = 20
    loss_function = nn.MSELoss()
    optimizer = torch.optim.Adam(model.parameters(), lr=learning_rate)

    train_losses = []
    validation_losses = []
    for epoch in range(epochs):
        tl = train_epoch(epoch, epochs, model, train_dataloader, loss_function, optimizer, device)
        vl = validate_epoch(epoch, epochs, model, validation_dataloader, loss_function, device)
        if epoch == 0: continue
        train_losses.append(tl)
        validation_losses.append(vl)

    torch.save(model.state_dict(), os.path.abspath("./model.pth"))

    plt.plot(np.linspace(1, epochs, epochs-1), np.array([train_losses, validation_losses]).T, label="train")
    plt.show()
