import math
import os.path

import numpy as np
import torch
from matplotlib import pyplot as plt
from tqdm import tqdm
from torch import nn

from torch.utils.data import DataLoader
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
    sequence_length = 20
    player_dataset = PlayerDataset(os.path.abspath("../data/chunk_0.csv"), sequence_length, "h_10", 3)
    train_size = math.floor(len(player_dataset)*dataset_split_ratio)
    validation_size = len(player_dataset) - train_size
    train_set, validation_set = torch.utils.data.random_split(player_dataset, [train_size, validation_size])

    #dataset = MatchDataset(os.path.abspath("../data/chunk_0.csv"), device, sequence_length, n_nearest_players=n_nearest_players)
    train_dataloader = DataLoader(train_set, batch_size=batch_size, shuffle=True, num_workers=8)
    validation_dataloader = DataLoader(validation_set, batch_size=batch_size, shuffle=True, num_workers=8)

    m = PlayerPredictor(device, n_nearest_players, 32, 4)
    m.to(device)

    learning_rate = 1e-4
    epochs = 2
    loss_function = nn.MSELoss()
    optimizer = torch.optim.Adam(m.parameters(), lr=learning_rate)

    train_losses = []
    validation_losses = []
    for epoch in range(epochs):
        vl = validate_epoch(epoch, epochs, m, validation_dataloader, loss_function, device)
        tl = train_epoch(epoch, epochs, m, train_dataloader, loss_function, optimizer, device)
        train_losses.append(tl)
        validation_losses.append(vl)

    plt.plot(np.linspace(1, epochs, epochs), train_losses, label="train")
    plt.show()
