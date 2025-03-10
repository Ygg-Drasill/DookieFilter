import torch
from torch import nn

class PlayerPredictor(nn.Module):
    def __init__(self, device, n_nearest_players, n_hidden):
        super().__init__()
        self.device = device
        self.n_nearest_players = n_nearest_players
        self.n_hidden = n_hidden
        # input size is n_nearest_players *2 (home and away) + target player and ball
        input_size = ((2*n_nearest_players) + 2) * 2
        self.lstm = nn.LSTM(input_size, n_hidden, 4)
        self.linear = nn.Linear(n_hidden, 2)

    def forward(self, x: torch.Tensor, future=0):
        outputs = []
        x = x.flatten(start_dim=-2)
        n_samples = x.size[0]

        h_t1 = torch.zeros(n_samples, self.n_hidden, dtype=torch.float32).to(x.device)
        c_t1 = torch.zeros(n_samples, self.n_hidden, dtype=torch.float32).to(x.device)
        h_t2 = torch.zeros(n_samples, self.n_hidden, dtype=torch.float32).to(x.device)
        c_t2 = torch.zeros(n_samples, self.n_hidden, dtype=torch.float32).to(x.device)

        return x


def input_to_tensor():
    pass
