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

    def __init__(self, data_path: str, device: torch.device):
        self.device = device
        if not os.path.exists(data_path) or os.path.isdir(data_path):
            return
        self.data_path = data_path
        self.match_dataframe = pd.read_csv(self.data_path)
        self.match_dataframe.reset_index()
        self.frames_per_player = len(self.match_dataframe)
        self.player_numbers = []
        for key in self.match_dataframe.columns:
            key_split = key.split("_")
            player_key = key_split[0] + "_" + key_split[1]
            if player_key[0] != "h" and player_key[0] != "a":
                continue
            if player_key not in self.player_numbers:
                self.player_numbers.append(player_key)

        self.length = len(self.player_numbers) * len(self.match_dataframe)

    def __len__(self):
        return self.length

    def __getitem__(self, idx):
        player_frame_index = idx % self.frames_per_player
        player_number = self.player_numbers[idx // self.frames_per_player]
        ball, player, home, away = self.get_player_ball__n_nearest(player_frame_index, player_number, 3)
        sample = [ball, player]
        for k in dict.keys(home): sample.append(home[k])
        for k in dict.keys(away): sample.append(away[k])
        return torch.from_numpy(np.array(sample)).to(self.device)

    def get_player_ball__n_nearest(self, idx: int, player_number: str, n: int):
        frame_coords = self.collect_coords_at(idx)
        if player_number not in dict.keys(frame_coords): return []

        ball_coords = frame_coords.pop("ball")
        frame_coords_keys = dict.keys(frame_coords.copy())
        other_keys = frame_coords.copy()
        player = other_keys.pop(player_number)
        home_distances, away_distances = {}, {}
        for other in filter(lambda x: x[0] == "h", other_keys):
            home_distances[other] = np.linalg.norm(np.array(player) - np.array(other_keys[other]))
        for other in filter(lambda x: x[0] == "a", other_keys):
            away_distances[other] = np.linalg.norm(np.array(player) - np.array(other_keys[other]))
        home_closest_keys = sorted(dict.keys(home_distances), key=lambda x: home_distances[x])[:n]
        away_closest_keys = sorted(dict.keys(away_distances), key=lambda x: away_distances[x])[:n]
        home, away = {}, {}
        for key in home_closest_keys:
            home[key] = frame_coords[key]
        for key in away_closest_keys:
            away[key] = frame_coords[key]
        return ball_coords, player, home, away


    def collect_coords_at(self, idx: int):
        frame = self.match_dataframe.loc[idx]
        ball = [frame["ball_x"], frame["ball_y"]]
        frame_coords = {"ball": ball}
        for player in self.player_numbers:
            frame_coords[player] = [frame[player + "_x"], frame[player + "_y"]]
        return frame_coords
