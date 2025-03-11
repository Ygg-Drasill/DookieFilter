import heapq

import numpy as np
import pandas as pd
import torch
from torch.utils.data import Dataset
from tqdm import tqdm


class PlayerDataset(Dataset):
    def __init__(self, path: str, sequence_length: int, target_player_key: str, n_nearest_players: int):
        self.raw_dataframe = pd.read_csv(path)
        self.frame_index = self.raw_dataframe.pop("frame_index")
        self.sequence_length = sequence_length
        self.n_nearest_players = n_nearest_players
        self.player_numbers = []

        ball = np.array([[x, y] for x, y in zip(self.raw_dataframe.pop("ball_x"),
                                                self.raw_dataframe.pop("ball_y"))])
        self.player = np.array([[x, y] for x, y in zip(self.raw_dataframe.pop(target_player_key + "_x"),
                                                  self.raw_dataframe.pop(target_player_key + "_y"))])

        self.player_numbers = set([k[:-2] for k in self.raw_dataframe.keys()])

        other_players = {}
        for player_number in self.player_numbers:
            other_players[player_number] = np.column_stack([self.raw_dataframe[player_number + "_x"],
                                                            self.raw_dataframe[player_number + "_y"]])

        self.data = []
        for i, pos in enumerate(tqdm(self.player, unit="row", desc="calculating rows")):
            home_distances, away_distances = [], []
            for key, coords in other_players.items():
                other_coords = np.array(coords[i])
                if np.isnan(other_coords).any():
                    continue
                distance = np.linalg.norm(self.player[i] - other_coords)
                prefix = key[0]
                if prefix == "h":
                    home_distances.append((distance, key))
                elif prefix == "a":
                    away_distances.append((distance, key))

            home_closest_keys = [k for _, k in heapq.nsmallest(n_nearest_players, home_distances)]
            away_closest_keys = [k for _, k in heapq.nsmallest(n_nearest_players, away_distances)]
            data_row = [self.normalize_x(pos[0]), self.normalize_y(pos[1]),
                        self.normalize_x(ball[i][0]), self.normalize_y(ball[i][0])]

            for k in home_closest_keys:
                coords = other_players[k][i]
                data_row.append(self.normalize_x(coords[0]))
                data_row.append(self.normalize_y(coords[1]))

            for k in away_closest_keys:
                coords = other_players[k][i]
                data_row.append(self.normalize_x(coords[0]))
                data_row.append(self.normalize_y(coords[1]))

            self.data.append(data_row)

        self.data = np.array(self.data)

    def __len__(self):
        return len(self.data) - self.sequence_length

    def __getitem__(self, idx):
        seq = self.data[idx:idx + self.sequence_length]
        target = self.player[idx + self.sequence_length]
        target_normalized = [self.normalize_x(target[0]), self.normalize_y(target[1])]
        return (torch.from_numpy(np.array(seq)).to(dtype=torch.float32),
                torch.from_numpy(np.array(target_normalized)).to(dtype=torch.float32))

    field_width = 145
    field_height = 68
    field_x_offset = field_width/2
    field_y_offset = field_height/2
    def normalize_x(self, x):
        return (x + self.field_x_offset) / self.field_width
    def normalize_y(self, y):
        return (y + self.field_y_offset) / self.field_height
