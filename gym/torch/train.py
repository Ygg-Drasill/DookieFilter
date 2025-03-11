import os.path

import torch
from tqdm import tqdm
from torch import nn

from torch.utils.data import DataLoader
from dataloader import MatchDataset
from model import PlayerPredictor

if __name__ == '__main__':
    device = torch.device("cuda:0" if torch.cuda.is_available() else "cpu")
    dataset = MatchDataset(os.path.abspath("../data/chunk_0.csv"), device, 2)
    dataloader = DataLoader(dataset, batch_size=1, shuffle=True, num_workers=4)

    m = PlayerPredictor(device, 3, 16, 4)
    m.to(device)

    learning_rate = 1e-4
    epochs = 10
    loss_function = nn.MSELoss()
    optimizer = torch.optim.Adam(m.parameters(), lr=learning_rate)

    for epoch in range(epochs):
        m.train(True)
        running_loss = 0.0
        for batch_index, batch in enumerate(dataloader):# enumerate(tqdm(dataloader)):
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
                print('Batch {0}, Loss: {1:.3f}'.format(batch_index + 1,
                                                        avg_loss_across_batches))
                running_loss = 0.0
