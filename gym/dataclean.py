import csv
import json
import os
import sys
import numpy as np

data_path = sys.argv[1]
output_target = sys.argv[2]

file = open(data_path)
data = []
data_fields = ["frame_index", "ball_x", "ball_y"]
frame_idx = []
chunk_index = 0

def save_game_chunk(chunk_file_index, output_target_dir, fields, data_chunk):
    file_name = "chunk_" + str(chunk_file_index) + ".csv"
    csv_file = open(output_target_dir + file_name, "w")
    writer = csv.DictWriter(csv_file, fieldnames=fields)
    writer.writeheader()
    writer.writerows(data_chunk)
    csv_file.close()

def coords_to_string(xy):
    return str(xy[0]) + ";" + str(xy[1])

def coords_from_string(xy):
    parts = xy[0].split(";")
    x = float(parts[0])
    y = float(parts[1])
    return np.array([x, y])

def add_player_data(player_data, frame_data, frame_fields, player_prefix: str):
    player_coords, player_number = player_data["xyz"], player_data["number"]
    key_x = player_prefix + str(player_number) + "_x"
    key_y = player_prefix + str(player_number) + "_y"
    frame_data[key_x] = player_coords[0]
    frame_data[key_y] = player_coords[1]
    if key_x not in frame_fields:
        frame_fields.append(key_x)
    if key_y not in frame_fields:
        frame_fields.append(key_y)

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
                save_game_chunk(chunk_index, output_target + "/", data_fields, data)
                chunk_index += 1
                data = []
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
            "frame_index": idx,
            "ball_x": seperated_frame["ball"]["xyz"][0],
            "ball_y": seperated_frame["ball"]["xyz"][1],
        }

        for player in home_players:
            add_player_data(player, frame, data_fields, "h_")
        for player in away_players:
            add_player_data(player, frame, data_fields, "a_")

        data.append(frame)
