import csv
import json
import os
import sys
import numpy as np
import matplotlib.pyplot as plt
from numpy.ma.core import append

data_path = sys.argv[1]
output_target = sys.argv[2]

file = open(data_path)
data = []
data_fields = ["frameIdx", "ball"]
frame_idx = []
chunk_index = 0

def save_game_chunk(chunk_index=None):
    file_name = "chunk_" + str(chunk_index) + ".csv"
    chunk_index += 1
    csv_file = open(output_target + "/" + file_name, "w")
    writer = csv.DictWriter(csv_file, fieldnames=data_fields)
    writer.writeheader()
    writer.writerows(data)

def coords_to_string(xy):
    return str(xy[0]) + ";" + str(xy[1])

def coords_from_string(xy):
    parts = xy[0].split(";")
    x = float(parts[0])
    y = float(parts[1])
    return np.array([x, y])

init = False
game_active = False
for line in file:
    packet = json.loads(line)
    if not init and not os.path.exists(output_target):
        init = True
        output_target = output_target + "/" + packet["gameId"]
        os.mkdir(output_target)

    for seperated_frame in packet["data"]:
        if "frameIdx" not in dict.keys(seperated_frame):
            print("signal")
            if game_active:
                save_game_chunk(chunk_index)
            game_active = not game_active
            continue
        if len(seperated_frame["homePlayers"]) == 0 or len(seperated_frame["awayPlayers"]) == 0:
            continue
        if not game_active:
            continue

        idx = seperated_frame["frameIdx"]
        frame_idx.append(idx)
        home_players = seperated_frame["homePlayers"]
        away_players = seperated_frame["awayPlayers"]
        frame = {
            "frameIdx": idx,
            "ball": coords_to_string(seperated_frame["ball"]["xyz"][:2])
        }

        for player in home_players:
            coords = coords_to_string(player["xyz"][:2])
            number = player["number"]
            player_key = "h" + str(number)
            frame[player_key] = coords
            if player_key not in data_fields:
                data_fields.append(player_key)

        for player in away_players:
            coords = coords_to_string(player["xyz"][:2])
            number = player["number"]
            player_key = "a" + str(number)
            frame[player_key] = coords
            if player_key not in data_fields:
                data_fields.append(player_key)

        data.append(frame)

# plt.plot(*np.array(data["a9"]).T)
# plt.show()
