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
    sequence_length = 20
    chunk_path = os.path.abspath("../data/chunk_0.csv")
    player_numbers = ["h_10", "h_13", "h_14", "h_2", "h_8", "a_1", "a_3", "a_19", "a_26", "a_44"]
    player_dataset = ConcatDataset([PlayerDataset(chunk_path, sequence_length, num, 3) for num in player_numbers])
    train_size = math.floor(len(player_dataset)*dataset_split_ratio)
    validation_size = len(player_dataset) - train_size
    train_set, validation_set = torch.utils.data.random_split(player_dataset, [train_size, validation_size])


    #dataset = MatchDataset(os.path.abspath("../data/chunk_0.csv"), device, sequence_length, n_nearest_players=n_nearest_players)
    train_dataloader = DataLoader(train_set, batch_size=batch_size, shuffle=True, num_workers=8)
    validation_dataloader = DataLoader(validation_set, batch_size=batch_size, shuffle=True, num_workers=8)

    model = PlayerPredictor(device, n_nearest_players, 32, 4)
    model.to(device)

    learning_rate = 1e-4
    epochs = 5
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
