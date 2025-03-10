import os.path
from tqdm import tqdm

from sklearn.preprocessing import MinMaxScaler

from torch.utils.data import DataLoader
from dataloader import MatchDataset

if __name__ == '__main__':
    device = "cpu"# torch.device("cuda:0" if torch.cuda.is_available() else "cpu")
    dataset = MatchDataset(os.path.abspath("../data/chunk_0.csv"), device, 2)
    dataloader = DataLoader(dataset, batch_size=1, shuffle=True, num_workers=12)
    for batch, y in tqdm(dataloader):
        print(batch)
        break
        pass
