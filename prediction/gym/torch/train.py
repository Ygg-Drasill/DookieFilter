import math
import os.path

import torch
from torch import nn
from torch.utils.data import DataLoader, ConcatDataset
from torch.utils.tensorboard import SummaryWriter

from epoch import train_epoch, validate_epoch
from gym.board_logger import BoardLogger
from model.player_predictor import PlayerPredictor
from player_dataset import PlayerDataset


def format_model_name(n_nearest_players, stack_size, hidden_size, lr, epochs, batch_size, parameters):
    return f'{n_nearest_players}-{stack_size}-{hidden_size}-{lr}-{epochs}-{batch_size}-{parameters}'

if __name__ == '__main__':
    hyper_parameters = {
        'n_nearest_players': range(1, 5+1),
        'stack_size': range(1, 16+1),
        'hidden_size': [8,16,32,64,128,256,512],
        'sequence_length': [10,20,30,40],
        'batch_size': [16,32,64,128,256],
        'lr': [0.0001,0.00001, 0.1, 0.001],
    }

    target_directory = os.path.abspath(f'../runs')
    summary_writer = SummaryWriter(f'{target_directory}/board')

    torch.random.seed()
    dataset_split_ratio = 0.8
    device = torch.device("cuda:0" if torch.cuda.is_available() else "cpu")
    chunk_path = os.path.abspath("../data/chunk_0.csv")
    player_numbers = ["h_10"]#["h_10", "h_13", "h_14", "h_2", "h_8", "a_1", "a_3", "a_19", "a_26", "a_44"]

    for n_nearest in hyper_parameters['n_nearest_players']:
        for stack_size in hyper_parameters['stack_size']:
            for hidden_size in hyper_parameters['hidden_size']:
                for sequence_length in hyper_parameters['sequence_length']:
                    for batch_size in hyper_parameters['batch_size']:
                        for lr in hyper_parameters['lr']:
                            print(f'n:{n_nearest} stack_size:{stack_size} hidden_size:{hidden_size} seq:{sequence_length} lr:{lr} batch_size:{batch_size}')
                            player_dataset = ConcatDataset(
                                [PlayerDataset(chunk_path, sequence_length, num, n_nearest) for num in player_numbers])
                            train_size = math.floor(len(player_dataset) * dataset_split_ratio)
                            validation_size = len(player_dataset) - train_size
                            train_set, validation_set = torch.utils.data.random_split(player_dataset,
                                                                                      [train_size, validation_size])
                            train_dataloader = DataLoader(train_set, batch_size=batch_size, shuffle=True, num_workers=8)
                            validation_dataloader = DataLoader(validation_set, batch_size=batch_size, shuffle=True,
                                                               num_workers=8)

                            model = PlayerPredictor(device, n_nearest, hidden_size, stack_size)
                            model.to(device)
                            n_parameters = sum(p.numel() for p in model.parameters() if p.requires_grad)
                            epochs = 10
                            loss_function = nn.MSELoss()
                            optimizer = torch.optim.Adam(model.parameters(), lr=lr)

                            train_losses = []
                            validation_losses = []

                            training_step = 0
                            validation_step = 0
                            writer = SummaryWriter(f'{target_directory}/player/{format_model_name(n_nearest, stack_size, hidden_size, lr, epochs, batch_size, n_parameters)}')
                            training_logger = BoardLogger(writer)
                            validation_logger = BoardLogger(writer)
                            train_loss, validation_loss = 0, 0
                            for epoch in range(epochs):
                                train_loss += train_epoch(epoch, epochs, model, train_dataloader, loss_function, optimizer, device, training_logger)
                                validation_loss += validate_epoch(epoch, epochs, model, validation_dataloader, loss_function, device, validation_logger)

                            writer.add_hparams({
                                'batch_size': batch_size,
                                'n_nearest_players': n_nearest,
                                'hidden_size': hidden_size,
                                'stack_size': stack_size,
                                'sequence_length': sequence_length,
                                'lr': lr,
                            },{
                               'train_loss': train_loss / epochs,
                               'test_loss': validation_loss / epochs
                            })
                            writer.flush()

                            torch.save(model.state_dict(),
                                       f'{target_directory}/models/{format_model_name(n_nearest, stack_size, hidden_size, lr, epochs, batch_size, n_parameters)}.pt')
