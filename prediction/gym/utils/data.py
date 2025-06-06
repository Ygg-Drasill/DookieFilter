import heapq

import numpy as np
from numpy import ndarray

FIELD_WIDTH = 105
FIELD_HEIGHT = 68
FIELD_X_OFFSET = FIELD_WIDTH / 2
FIELD_Y_OFFSET = FIELD_HEIGHT / 2

def normalize_x(x): return (x + FIELD_X_OFFSET) / FIELD_WIDTH
def normalize_y(y): return (y + FIELD_Y_OFFSET) / FIELD_HEIGHT
def denormalize_x(x): return x * FIELD_WIDTH - FIELD_X_OFFSET
def denormalize_y(y): return y * FIELD_HEIGHT - FIELD_Y_OFFSET
def m_to_cm(a: float) -> float: return a*100.0

def get_features(
        dataframe_row:dict[str, ndarray[int]],
        n_nearest_players:int,
        player_key:str) -> ndarray:
    dataframe_row.pop("frame_index")
    ball = np.array([dataframe_row.pop("ball_x"),
                     dataframe_row.pop("ball_y")])
    player = np.array([dataframe_row.pop(player_key + "_x"),
                       dataframe_row.pop(player_key + "_y")])
    features = [normalize_x(player[0]), normalize_y(player[1]),
                normalize_x(ball[0]), normalize_y(ball[1])]

    is_home = player_key.split("_")[0] == "h"
    player_numbers = set([k[:-2] for k in dataframe_row.keys()])
    other_players = {key:[dataframe_row.pop(key+ "_x"), dataframe_row.pop(key+ "_y")] for key in player_numbers}

    home_closest_keys, away_closest_keys = get_nearest_players(player, other_players, n_nearest_players)
    team_keys = home_closest_keys if is_home else away_closest_keys
    opponent_keys = away_closest_keys if is_home else home_closest_keys
    for k in team_keys:
        coords = other_players[k]
        features.append(normalize_x(coords[0]))
        features.append(normalize_y(coords[1]))
    for k in opponent_keys:
        coords = other_players[k]
        features.append(normalize_x(coords[0]))
        features.append(normalize_y(coords[1]))
    return np.array(features, dtype=float)


def get_nearest_players(
        target,
        others: dict[str, list[ndarray[int]]],
        n_nearest_players: int) -> (list[str], list[str]):
    home_distances, away_distances = [], []
    for key, coords in others.items():
        if np.isnan(coords).any():
            continue
        distance = np.linalg.norm(target - coords)
        prefix = key[0]
        if prefix == "h":
            home_distances.append((distance, key))
        elif prefix == "a":
            away_distances.append((distance, key))

    home_closest_keys = [k for _, k in heapq.nsmallest(n_nearest_players, home_distances)]
    away_closest_keys = [k for _, k in heapq.nsmallest(n_nearest_players, away_distances)]
    return home_closest_keys, away_closest_keys





