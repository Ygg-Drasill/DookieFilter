import math
import os.path
from os.path import abspath

import numpy as np
import torch
from matplotlib import pyplot as plt
from torch import nn

from torch.utils.data import DataLoader, ConcatDataset
from dataloader import MatchDataset
from model.player_predictor import PlayerPredictor
from player_dataset import PlayerDataset
from epoch import train_epoch, validate_epoch
from torch.utils.tensorboard import SummaryWriter

def format_model_name(n_nearest_players, stack_size, hidden_size, lr, epochs, batch_size, parameters):
    return f'{n_nearest_players}-{stack_size}-{hidden_size}-{lr}-{epochs}-{batch_size}-{parameters}'

if __name__ == '__main__':
    target_directory = os.path.abspath(f'../runs')
    summary_writer = SummaryWriter(f'{target_directory}/board')

    torch.random.seed()
    dataset_split_ratio = 0.8
    batch_size = 128
    n_nearest_players = 4
    device = torch.device("cuda:0" if torch.cuda.is_available() else "cpu")
    sequence_length = 40
    chunk_path = os.path.abspath("../data/chunk_0.csv")
    player_numbers = ["h_10"]#["h_10", "h_13", "h_14", "h_2", "h_8", "a_1", "a_3", "a_19", "a_26", "a_44"]
    player_dataset = ConcatDataset([PlayerDataset(chunk_path, sequence_length, num, n_nearest_players) for num in player_numbers])
    train_size = math.floor(len(player_dataset)*dataset_split_ratio)
    validation_size = len(player_dataset) - train_size
    train_set, validation_set = torch.utils.data.random_split(player_dataset, [train_size, validation_size])


    #dataset = MatchDataset(os.path.abspath("../data/chunk_0.csv"), device, sequence_length, n_nearest_players=n_nearest_players)
    train_dataloader = DataLoader(train_set, batch_size=batch_size, shuffle=True, num_workers=8)
    validation_dataloader = DataLoader(validation_set, batch_size=batch_size, shuffle=True, num_workers=8)

    hidden_size, stack_size = 128, 8
    model = PlayerPredictor(device, n_nearest_players, hidden_size, stack_size)
    model.to(device)

    learning_rate = 1e-4
    epochs = 10
    loss_function = nn.MSELoss()
    optimizer = torch.optim.Adam(model.parameters(), lr=learning_rate)

    train_losses = []
    validation_losses = []

    steps = 0
    for epoch in range(epochs):
        tl = train_epoch(epoch, epochs, model, train_dataloader, loss_function, optimizer, device)
        vl = validate_epoch(epoch, epochs, model, validation_dataloader, loss_function, device)
        if epoch == 0: continue
        train_losses.append(tl)
        validation_losses.append(vl)
        steps += 1

    n_parameters = sum(p.numel() for p in model.parameters() if p.requires_grad)
    torch.save(model.state_dict(),
               f'{target_directory}/out/{format_model_name(n_nearest_players, stack_size, hidden_size, learning_rate, epochs, batch_size, n_parameters)}.pt')
