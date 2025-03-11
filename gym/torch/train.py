import os.path

import torch
from tqdm import tqdm
from torch import nn

from torch.utils.data import DataLoader
from dataloader import MatchDataset
from model import PlayerPredictor

if __name__ == '__main__':
    n_nearest_players = 3
    device = torch.device("cuda:0" if torch.cuda.is_available() else "cpu")
    sequence_length = 20
    dataset = MatchDataset(os.path.abspath("../data/chunk_0.csv"), device, sequence_length, n_nearest_players=n_nearest_players)
    dataloader = DataLoader(dataset, batch_size=1, shuffle=True, num_workers=8)

    m = PlayerPredictor(device, n_nearest_players, 32, 4)
    m.to(device)

    learning_rate = 1e-4
    epochs = 10
    loss_function = nn.MSELoss()
    optimizer = torch.optim.Adam(m.parameters(), lr=learning_rate)

    for epoch in range(epochs):
        m.train(True)
        running_loss = 0.0
        progress_bar = tqdm(enumerate(dataloader), total=len(dataloader))
        for batch_index, batch in progress_bar:
            batch_x, batch_y = batch[0].to(device), batch[1].to(device)
            if torch.isnan(batch_x).any() or torch.isnan(batch_y).any():
                continue

            output = m(batch_x)
            loss = loss_function(output, batch_y)
            running_loss += loss.item()

            optimizer.zero_grad()
            loss.backward()
            torch.nn.utils.clip_grad_norm_(m.parameters(), max_norm=1.0)
            optimizer.step()
            if batch_index % 100 == 99:  # print every 100 batches

                avg_loss_across_batches = running_loss / 100
                progress_bar.set_postfix({'loss': avg_loss_across_batches})
                running_loss = 0.0
