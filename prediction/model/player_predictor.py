import torch
from torch import nn

MAX_DELTA = 0.5
FIELD_WIDTH = 105
FIELD_HEIGHT = 68

MAX_NORMALIZED_X_DELTA = MAX_DELTA / FIELD_WIDTH
MAX_NORMALIZED_Y_DELTA = MAX_DELTA / FIELD_HEIGHT

class PlayerPredictor(nn.Module):
    def __init__(self, device, n_nearest_players, n_hidden, n_stack):
        super().__init__()
        self.device = device
        self.n_nearest_players = n_nearest_players
        self.n_hidden = n_hidden
        self.n_stack = n_stack
        # input size is n_nearest_players *2 (home and away) + target player and ball
        self.input_size = ((2*n_nearest_players) + 2) * 2
        self.lstm = nn.LSTM(self.input_size, n_hidden, n_stack, batch_first=True, dropout=0.2)
        self.linear = nn.Linear(n_hidden, 2)

    def forward(self, x: torch.Tensor) -> torch.Tensor:
        batch_size = x.size(0)
        h0 = torch.zeros(self.n_stack, batch_size, self.n_hidden, dtype=torch.float32).to(self.device)
        c0 = torch.zeros(self.n_stack, batch_size, self.n_hidden, dtype=torch.float32).to(self.device)

        out, (hn, cn) = self.lstm(x, (h0, c0))
        out = self.linear(out[:, -1, :])
        out = torch.sigmoid(out)

        decoded_delta_x = (out[:, 0] - 0.5) * (2 * MAX_NORMALIZED_X_DELTA)
        decoded_delta_y = (out[:, 1] - 0.5) * (2 * MAX_NORMALIZED_Y_DELTA)
        decoded_delta = torch.stack([decoded_delta_x, decoded_delta_y], dim=1)

        prev_pos = x[:, -1, 0:2] #pick the first two elements of the input tensor (target player position)
        prediction = prev_pos + decoded_delta

        return prediction, decoded_delta
