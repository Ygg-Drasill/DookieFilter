import csv
from chunk import Chunk
from typing import Any

class Chunk:
    def __init__(self):
        self.count = 0
        self.data_fields = ["frame_index", "ball_x", "ball_y"]
        self.data: list[dict[str, Any]] = []

    def add_data(self, data: dict[str, Any]) -> None:
        self.data.append(data)
        self.count += 1

    def add_frame(self, frame):
        homeplayers = frame["homePlayers"]
        awayplayers = frame["awayPlayers"]
        row = {
            "frame_index": frame["frameIndex"],
            "ball_x": frame["ball"]["xyz"][0],
            "ball_y": frame["ball"]["xyz"][1],
        }

        for player in homeplayers:
            self.add_player_data(player, row,  "h_")
        for player in awayplayers:
            self.add_player_data(player, row,  "a_")

        self.data.append(row)

    def write_to_file(self, output_target : str):
        csv_file = open(output_target, "w")
        writer = csv.DictWriter(csv_file, fieldnames=self.data_fields)
        writer.writeheader()
        writer.writerows(self.data)
        csv_file.close()

    def filter(self) -> list[Chunk]:

        last_row = 0
        sub_chunks : list[Chunk] = []
        dirty = False
        current_chunk = Chunk()

        for row in self.data:

            if dirty:
                continue

            current_chunk.add_data(row)

    def add_player_data(self, player_data, row, player_prefix: str):
        player_coords, player_number = player_data["xyz"], player_data["number"]
        key = player_prefix + str(player_number)
        coords = [player_coords[0], player_coords[1]]
        row[key] = coords
        key_x = key + "_x"
        key_y = key + "_y"

        if key_x not in self.data_fields:
            self.data_fields.append(key_x)
        if key_y not in self.data_fields:
            self.data_fields.append(key_y)
