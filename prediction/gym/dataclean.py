import csv
import json
import os
import sys
import uuid

import numpy as np
from tqdm import tqdm

from gym.chunk.chunk import Chunk

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



data_path = sys.argv[1]
output_target = sys.argv[2]


def run(path):
    match_output_target = ""
    file = open(path)
    chunks: list[Chunk] = []
    data_fields = ["frame_index", "ball_x", "ball_y"]
    frame_idx = []
    chunk_index = 0

    game_active = False

    current_chunk = Chunk()
    init = False

    for line in file:
        packet = json.loads(line)

        if not init:
            init = True
            if not os.path.exists(output_target):
                os.mkdir(output_target)
            game_id = ""
            if "gameId" in packet:
                game_id = packet["gameId"]
            else:
                game_id = "unknown" + str(uuid.uuid4())
            match_output_target = output_target + "/" + game_id
            if not os.path.exists(match_output_target):
                os.mkdir(match_output_target)

        for seperated_frame in packet["data"]:
            if "frameIdx" not in dict.keys(seperated_frame):
                print("signal")
                chunks.append(current_chunk)
                current_chunk = Chunk()
                continue

            if len(seperated_frame["homePlayers"]) == 0 or len(seperated_frame["awayPlayers"]) == 0:
                continue

            current_chunk.add_frame(seperated_frame)

    sub_chunks = []
    active_chunks = sorted(chunks, key=lambda c: c.count)[-2:]
    print(len(active_chunks))

    for i, chunk in enumerate(active_chunks):
        print(f"filtering chunk {i}")
        for sub in chunk.filter():
            sub_chunks.append(sub)

    sub_chunks = [x for x in sub_chunks if x.count >= 50]

    print(f"found chunks: {len(sub_chunks)}")
    for i, chunk in enumerate(tqdm(sub_chunks, desc="writing chunks to disk")):
        chunk.write_to_file(f"{match_output_target}/chunk_{i}.csv")

if __name__ == "__main__":
    files = os.listdir(data_path)
    for f in files:
        run(f"{data_path}/{f}")

#1307
