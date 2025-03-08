from torch import nn

class Model(nn.Module):
    def __init__(self, input_size, output_size):
        super(Model, self).__init__()
        self.fc1 = nn.Linear(input_size, output_size)

    def forward(self, x):
        return x


def input_to_tensor():
    pass
