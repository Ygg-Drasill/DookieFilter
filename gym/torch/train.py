import os.path

import torch
from torch.utils.data import DataLoader

from dataloader import MatchDataset

if __name__ == '__main__':
    device = torch.device("cuda:0" if torch.cuda.is_available() else "cpu")
    dataset = MatchDataset(os.path.abspath("../data/chunk_0.csv"), device)
    dataloader = DataLoader(dataset, batch_size=64, shuffle=True)
    for batch in dataloader:
        print(batch.shape)
        break
