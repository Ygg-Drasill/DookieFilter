import heapq
import os.path

import numpy as np
import pandas as pd
import torch
from torch.utils.data import Dataset


class MatchDataset(Dataset):
    data_path: str
    player_numbers: list[str]
    frames_per_player: int
    length: int
    field_width: int = 145
    field_x_offset = 145/2
    field_height: int = 68
    field_y_offset = 68 / 2

    def __init__(
            self,
            data_path: str,
            device: torch.device,
            sequence_length: int,
            n_nearest_players=3):
        if not os.path.exists(data_path) or os.path.isdir(data_path):
            return
        self.device = device
        self.data_path = data_path
        self.n_nearest_players = n_nearest_players
        self.sequence_length = sequence_length
        self.match_dataframe = pd.read_csv(self.data_path)
        self.match_dataframe.reset_index()
        self.frames_per_player = len(self.match_dataframe)
        self.n_features = 2 + (n_nearest_players*2)
        self.player_numbers = []
        for key in self.match_dataframe.columns:
            key_split = key.split("_")
            player_key = key_split[0] + "_" + key_split[1]
            if player_key[0] != "h" and player_key[0] != "a":
                continue
            if player_key not in self.player_numbers:
                self.player_numbers.append(player_key)

        self.length = len(self.player_numbers) * (len(self.match_dataframe)-self.sequence_length)

    def __len__(self):
        return self.length

    def __getitem__(self, idx):
        player_number_index = idx // self.frames_per_player
        player_frame_index = ((idx + player_number_index*self.sequence_length) % (self.frames_per_player - 1))
        player_number = self.player_numbers[player_number_index]

        sequence = []
        for sequenced_idx in range(player_frame_index-self.sequence_length, player_frame_index):
            ball, player, home, away = self.get_player_ball__n_nearest(sequenced_idx, player_number, self.n_nearest_players)
            sample = [player, ball]
            for k in dict.keys(home):
                sample.append(home[k])
            for k in dict.keys(away):
                sample.append(away[k])
            for i in range(len(sample)):
                sample[i] = [self.normalize_x(sample[i][0]),
                             self.normalize_y(sample[i][1])]
            sequence.append(sample)

        next_frame = self.match_dataframe.loc[player_frame_index + 1]
        player_next = np.array([self.normalize_x(next_frame[player_number + "_x"]),
                                self.normalize_y(next_frame[player_number + "_y"])])
        return (torch.from_numpy(np.array(sequence).reshape(self.sequence_length, self.n_features, 2)).to(torch.float32),
                torch.from_numpy(np.array(player_next)).to(torch.float32))

    def normalize_x(self, x):
        return (x + self.field_x_offset) / self.field_width

    def normalize_y(self, y):
        return (y + self.field_y_offset) / self.field_height

    def get_player_ball__n_nearest(self, idx: int, player_number: str, n: int):
        """
        Find nearest home- and away players around specified player,
        and return ball, player and two lists with n nearest players
        """
        frame_coords = self.collect_coords_at(idx)
        if player_number not in dict.keys(frame_coords):
            return []

        ball_coords = frame_coords.pop("ball")
        frame_coords_keys = dict.keys(frame_coords.copy())
        other_keys = frame_coords.copy()
        player = np.array(other_keys.pop(player_number))
        home_distances, away_distances = [], []

        for key, coord in other_keys.items():
            distance = np.linalg.norm(player - np.array(coord))
            prefix = key[0]
            if prefix == "h":
                home_distances.append((distance, key))
            elif prefix == "a":
                away_distances.append((distance, key))

        home_closest_keys = [k for _, k in heapq.nsmallest(n, home_distances)]
        away_closest_keys = [k for _, k in heapq.nsmallest(n, away_distances)]

        home, away = {}, {}
        for key in home_closest_keys:
            home[key] = frame_coords[key]
        for key in away_closest_keys:
            away[key] = frame_coords[key]
        return ball_coords, player, home, away

    def collect_coords_at(self, idx: int) -> dict:
        """
        Return a row at index (dict) containing the ball position and
        all player positions, where all coordinates collected into 2D vectors
        """
        frame = self.match_dataframe.loc[idx]
        ball = [frame["ball_x"], frame["ball_y"]]
        frame_coords = {"ball": ball}
        for player in self.player_numbers:
            frame_coords[player] = [frame[player + "_x"], frame[player + "_y"]]
        return frame_coords
