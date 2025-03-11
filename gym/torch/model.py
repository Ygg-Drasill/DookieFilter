import torch
from torch import nn

class PlayerPredictor(nn.Module):
    def __init__(self, device, n_nearest_players, n_hidden, n_stack):
        super().__init__()
        self.device = device
        self.n_nearest_players = n_nearest_players
        self.n_hidden = n_hidden
        self.n_stack = n_stack
        # input size is n_nearest_players *2 (home and away) + target player and ball
        self.input_size = ((2*n_nearest_players) + 2) * 2
        self.lstm = nn.LSTM(self.input_size, n_hidden, n_stack, batch_first=True)
        self.linear = nn.Linear(n_hidden, 2)

    def forward(self, x: torch.Tensor):
        x = x.flatten(start_dim=-2)
        sq_len = x.size(1)
        n_samples = x.size(0)

        h0 = torch.zeros(self.n_stack, 1, self.n_hidden, dtype=torch.float32).to(self.device)
        c0 = torch.zeros(self.n_stack, 1, self.n_hidden, dtype=torch.float32).to(self.device)

        out, _ = self.lstm(x, (h0, c0))
        out = self.linear(out[:, -1, :])
        return out


def input_to_tensor():
    pass
