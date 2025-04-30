import heapq

import numpy as np
from numpy import ndarray

field_width = 105
field_height = 68
field_x_offset = field_width / 2
field_y_offset = field_height / 2

def normalize_x(x): return (x + field_x_offset) / field_width
def normalize_y(y): return (y + field_y_offset) / field_height


def get_features(
        dataframe_row:dict[str, ndarray[int]],
        n_nearest_players:int,
        player_key:str) -> ndarray:
    home_distances, away_distances = [], []
    dataframe_row.pop("frame_index")

    ball = np.array([dataframe_row.pop("ball_x"),
                     dataframe_row.pop("ball_y")])
    player = np.array([dataframe_row.pop(player_key + "_x"),
                       dataframe_row.pop(player_key + "_y")])

    player_numbers = set([k[:-2] for k in dataframe_row.keys()])

    other_players = {key:[dataframe_row.pop(key+ "_x"), dataframe_row.pop(key+ "_y")] for key in player_numbers}
    for key, coords in other_players.items():
        if np.isnan(coords).any():
            continue
        distance = np.linalg.norm(player - coords)
        prefix = key[0]
        if prefix == "h":
            home_distances.append((distance, key))
        elif prefix == "a":
            away_distances.append((distance, key))

    home_closest_keys = [k for _, k in heapq.nsmallest(n_nearest_players, home_distances)]
    away_closest_keys = [k for _, k in heapq.nsmallest(n_nearest_players, away_distances)]
    features = [normalize_x(player[0]), normalize_y(player[1]),
                normalize_x(ball[0]), normalize_y(ball[1])]

    for k in home_closest_keys:
        coords = other_players[k]
        features.append(normalize_x(coords[0]))
        features.append(normalize_y(coords[1]))

    for k in away_closest_keys:
        coords = other_players[k]
        features.append(normalize_x(coords[0]))
        features.append(normalize_y(coords[1]))

    return np.array(features, dtype=float)
