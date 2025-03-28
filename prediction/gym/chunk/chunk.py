import csv
import json
from typing import Self

from torch.onnx.symbolic_opset9 import linalg_norm

from typing import Any
import numpy as np

MAX_REALISTIC_DISTANCE = 0.42

class Chunk:
    def __init__(self):
        self.count: int = 0
        self.data_fields = {"frame_index", "ball_x", "ball_y"}
        self.data: list[dict[str, Any]] = []

    def add_data(self, data: dict[str, Any]) -> None:
        self.data.append(data)
        self.count += 1
        self.data_fields = self.data_fields.union(data.keys())

    def add_frame(self, frame):
        homeplayers = frame["homePlayers"]
        awayplayers = frame["awayPlayers"]
        row = {
            "frame_index": frame["frameIdx"],
            "ball_x": frame["ball"]["xyz"][0],
            "ball_y": frame["ball"]["xyz"][1],
        }

        for player in homeplayers:
            self.add_player_data(player, row,  "h_")
        for player in awayplayers:
            self.add_player_data(player, row,  "a_")

        self.count += 1
        self.data.append(row)


    def write_to_file(self, output_target : str):
        csv_file = open(output_target, "w")
        writer = csv.DictWriter(csv_file, fieldnames=self.data_fields)
        writer.writeheader()
        writer.writerows(self.data)
        csv_file.close()


    def filter(self) -> list[Self]:
        last_row: dict[str, Any] = None
        sub_chunks : list[Chunk] = []
        current_chunk = Chunk()

        for row in self.data:
            row.pop("frame_index")
            row.pop("ball_x")
            row.pop("ball_y")

            dirty = False
            if last_row is None:
                last_row = row
                continue

            for player in row:
                if len(set(row.keys()).difference(set(last_row.keys()))) > 0:
                    dirty = True
                    break
                last_pos = np.array(last_row[player], dtype=float)
                current_pos = np.array(row[player], dtype=float)
                distance = np.linalg.norm(current_pos - last_pos)
                dirty = dirty or distance > MAX_REALISTIC_DISTANCE

            last_row = row
            if dirty: #finish chunk and begin new one
                sub_chunks.append(current_chunk)
                current_chunk = Chunk()
                continue

            current_chunk.add_data(row)
        return sub_chunks


    def add_player_data(self, player_data, row, player_prefix: str):
        player_coords, player_number = player_data["xyz"], player_data["number"]
        key = player_prefix + str(player_number)
        coords = [player_coords[0], player_coords[1]]
        row[key] = coords
        key_x = key + "_x"
        key_y = key + "_y"

        if key_x not in self.data_fields:
            self.data_fields.add(key_x)
        if key_y not in self.data_fields:
            self.data_fields.add(key_y)
