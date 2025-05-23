import csv
import json
from typing import Self

from torch.onnx.symbolic_opset9 import linalg_norm

from typing import Any
import numpy as np

MAX_REALISTIC_DISTANCE = 0.42

class Chunk:
    def __init__(self):
        self.frameCount: int = 0
        self.data_fields = {"frame_index", "ball_x", "ball_y"}
        self.data: list[dict[str, Any]] = []

    def __length__(self):
        return self.frameCount

    def add_data(self, data: dict[str, Any]) -> None:
        self.data.append(data)
        self.frameCount += 1
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

        self.frameCount += 1
        self.data.append(row)

    def write_to_file(self, output_target: str):
        """Save this chunk in a file"""
        fieldnames = {"frame_index", "ball_x", "ball_y"}
        for player in self.data_fields.difference({"frame_index", "ball_x", "ball_y"}):
            for i, d in enumerate(self.data):
                try:
                    coords = d.pop(player)
                except KeyError:
                    continue
                d[player + "_x"] = coords[0]
                d[player + "_y"] = coords[1]
                fieldnames = fieldnames.union({f"{player}_x", f"{player}_y"})
                self.data[i] = d

        csv_file = open(output_target, "w", newline="")
        writer = csv.DictWriter(csv_file, fieldnames=fieldnames)
        writer.writeheader()
        writer.writerows(self.data)
        csv_file.close()

    def filter(self) -> list[Self]:
        """
        Create a list of new chunks for all realistic parts of the original chunk
        """
        last_row: dict[str, Any] = None
        sub_chunks : list[Chunk] = []
        current_chunk = Chunk()

        for row in self.data:
            if last_row is None:
                last_row = row
                continue

            #Mark as dirty if player numbers currently on field changed
            dirty = set(row.keys()) != set(last_row.keys()) #any(x is None for x in row.items()

            for player in set(row).difference({"frame_index", "ball_x", "ball_y"}):
                if dirty:
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
            else:
                current_chunk.add_data(row)

        return sub_chunks


    def add_player_data(self, player_data, row, player_prefix: str):
        player_coords, player_number = player_data["xyz"], player_data["number"]
        key = player_prefix + str(player_number)
        coords = [player_coords[0], player_coords[1]]
        key_x = key + "_x"
        key_y = key + "_y"
        row[key] = coords

        if key_x not in self.data_fields:
            self.data_fields.add(key_x)
        if key_y not in self.data_fields:
            self.data_fields.add(key_y)
