import math
import os
import pickle
import random
import typing

import pandas as pd
import torch
from torch.utils.data import Dataset, ConcatDataset
from tqdm import tqdm

from gym.utils.data import *


class PlayerDataset(Dataset):
    def __init__(self, path: str, sequence_length: int, n_nearest_players: int):
        self.raw_dataframe = pd.read_csv(path)
        self.frame_index = self.raw_dataframe.pop("frame_index")
        self.partition_size = len(self.frame_index)
        self.sequence_length = sequence_length
        self.n_nearest_players = n_nearest_players

        ball = np.array([[x, y] for x, y in zip(self.raw_dataframe.pop("ball_x"),
                                                self.raw_dataframe.pop("ball_y"))])

        self.player_numbers = list(set([k[:-2] for k in self.raw_dataframe.keys()]))

        self.data = {}
        pos: tuple[int, int]
        ball: np.ndarray[tuple[int, int]]

        for player_number in self.player_numbers:
            player = np.array([[x, y] for x, y in zip(self.raw_dataframe.get(player_number + "_x"),
                                                  self.raw_dataframe.get(player_number + "_y"))])
            other_players = {}
            for num in self.player_numbers:
                if num is player_number:
                    continue
                other_players[num] = np.column_stack([self.raw_dataframe.get(player_number + "_x"),
                                                                self.raw_dataframe.get(player_number + "_y")])

            player_data = []
            for i, pos in enumerate(player):
                home_distances, away_distances = [], []
                for key, coords in other_players.items():
                    other_coords = np.array(coords[i])
                    if np.isnan(other_coords).any():
                        continue
                    distance = np.linalg.norm(pos - other_coords)
                    prefix = key[0]
                    if prefix == "h":
                        home_distances.append((distance, key))
                    elif prefix == "a":
                        away_distances.append((distance, key))

                home_closest_keys = [k for _, k in heapq.nsmallest(n_nearest_players, home_distances)]
                away_closest_keys = [k for _, k in heapq.nsmallest(n_nearest_players, away_distances)]
                data_row = [normalize_x(pos[0]), normalize_y(pos[1]),
                            normalize_x(ball[i][0]), normalize_y(ball[i][1])]

                for k in home_closest_keys:
                    coords = other_players[k][i]
                    data_row.append(normalize_x(coords[0]))
                    data_row.append(normalize_y(coords[1]))

                for k in away_closest_keys:
                    coords = other_players[k][i]
                    data_row.append(normalize_x(coords[0]))
                    data_row.append(normalize_y(coords[1]))

                player_data.append(data_row)

            self.data[player_number] = np.array(player_data)

    def __len__(self):
        l = (len(self.data) - self.sequence_length) * len(self.player_numbers)
        if not (l >= 0):
            print("Warning: dataset length is broken")
            return 0
        return (len(self.data) - self.sequence_length) * len(self.player_numbers)

    def __getitem__(self, idx):
        partition_idx = idx // (self.partition_size - self.sequence_length)
        player_number = self.player_numbers[partition_idx]
        player = np.array([[x, y] for x, y in zip(self.raw_dataframe.get(player_number + "_x"),
                                                  self.raw_dataframe.get(player_number + "_y"))])

        #add offset based on what player partition of the dataset we are indexing
        no_padding_idx = (partition_idx*self.sequence_length + idx) % self.partition_size

        seq = self.data[player_number][no_padding_idx:no_padding_idx + self.sequence_length]
        target = player[no_padding_idx + self.sequence_length]
        target_normalized = [target[0], target[1]]
        return (torch.from_numpy(np.array(seq)).to(dtype=torch.float32),
                torch.from_numpy(np.array(target_normalized)).to(dtype=torch.float32))

    @staticmethod
    def from_dir(chunk_dir:str, n_nearest: int, sequence_length: int, split_ratio: float = 0.8):
        cache_file = os.path.abspath(os.path.join("../runs/cache", PlayerDataset.format_cache(n_nearest, sequence_length)))
        if os.path.exists(cache_file):
            dataset = PlayerDataset.load_cache(cache_file)
            return dataset["train"], dataset["validation"]

        chunk_files = []
        match_files = os.listdir(chunk_dir)
        for match_dir in match_files:
            md = os.listdir(chunk_dir + '/' + match_dir)
            for chunk_file in md:
                chunk_files.append(os.path.join(match_dir, chunk_file))

        random.shuffle(chunk_files)
        split_index = int(len(chunk_files) * split_ratio)
        train_chunks = chunk_files[:split_index]
        validation_chunks = chunk_files[split_index:]

        train_set = ConcatDataset(
            [PlayerDataset(f"{chunk_dir}/{path}", sequence_length, n_nearest) for path in tqdm(train_chunks,
                                                                                                unit='chunk',
                                                                                                desc='building training dataset from chunks',
                                                                                                ncols=200)])
        validation_set = ConcatDataset(
            [PlayerDataset(f"{chunk_dir}/{path}", sequence_length, n_nearest) for path in tqdm(validation_chunks,
                                                                                               unit='chunk',
                                                                                               desc='building training dataset from chunks',
                                                                                               ncols=200)])

        split_set = {"train": train_set, "validation": validation_set}
        PlayerDataset.dump_cache(cache_file, split_set)

        return train_set, validation_set

    @staticmethod
    def format_cache(n_nearest: int, sequence_length: int):
        return f"{n_nearest}-{sequence_length}.pickle"

    @staticmethod
    def load_cache(path:str):
        with open(path, 'rb') as f:
            data = pickle.load(f)
        return data

    @staticmethod
    def dump_cache(path:str, dataset:dict[str, Dataset[typing.Any]]):
        with open(path, 'wb') as f:
            pickle.dump(dataset, f, pickle.HIGHEST_PROTOCOL)
