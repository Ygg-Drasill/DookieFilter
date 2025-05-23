import typing

import numpy as np
import torch
import zmq
from numpy import ndarray

from gym.utils.data import denormalize_x, denormalize_y
from model.player_predictor import PlayerPredictor

PROTOCOL = "tcp"
ADDRESS = "localhost"

class ImputationService:
    def __init__(self, model: PlayerPredictor, sequence_length: int = 20):
        self.model = model
        self.nearest_players = model.n_nearest_players
        self.sequence_length = sequence_length
        self.context = zmq.Context()
        self.socket_storage_api = self.context.socket(zmq.REQ)
        self.socket_storage = self.context.socket(zmq.PUSH)
        self.socket_imputation = self.context.socket(zmq.PULL)
        self.player_series: dict[str, list[ndarray]] = dict[str, list[ndarray]]()
        self.active = False

        model.to(self.model.device)
        model.eval()
        torch.set_grad_enabled(False)

    def connect(self, imputation_port: int, storage_api_port: int, storage_port: int):
        self.socket_imputation.bind(f"{PROTOCOL}://{ADDRESS}:{imputation_port}")
        self.socket_storage_api.connect(f"{PROTOCOL}://{ADDRESS}:{storage_api_port}")
        self.socket_storage.connect(f"{PROTOCOL}://{ADDRESS}:{storage_port}")

    def close(self):
        self.socket_imputation.close()
        self.socket_storage_api.close()
        self.socket_storage.close()

    def run(self):
        self.active = True
        while self.active: self.listen()

    def listen(self):
        imputation_request = self.socket_imputation.recv_json()
        if imputation_request[0] == "kill":
            self.active = False
            return

        for request in imputation_request:
            frame_idx = int(request["frame_idx"])
            player_number = int(request["player_number"])
            sequence = self.get_player_history(frame_idx, player_number)
            target_next = self.predict(sequence)
            self.socket_storage.send_string("player", zmq.SNDMORE)
            self.socket_storage.send_json({
                "frame_idx": frame_idx,
                "player_number": player_number,
                "target_next": target_next
            })


    def get_player_history(self, frame_idx: int, player_number: int):
        self.socket_storage_api.send_string("frameRangeNearest", zmq.SNDMORE)
        self.socket_storage_api.send_json({
            "startIdx": frame_idx - self.sequence_length,
            "endIdx": frame_idx,
            "n": self.nearest_players,
            "target": player_number,
        })

        response = self.socket_storage_api.recv_json()
        return response

    def predict(self, sequence: list[dict[str, typing.Any]]) -> dict[str, float]:
        x = sequence_to_tensor(sequence)
        x = x.reshape(1, len(sequence), self.model.input_size).to(self.model.device)
        out: tuple[torch.Tensor] = self.model(x)
        prediction: torch.Tensor = out[0]
        target_next = prediction.detach().cpu().numpy().squeeze()
        return {"x": denormalize_x(target_next[0].item()), "y": denormalize_y(target_next[1].item())}

    def store_player_frame(self):
        pass

def sequence_to_tensor(sequence: list[dict[str, typing.Any]]) -> torch.Tensor:
    x = []
    for item in sequence:
        xt = [item["target"]["x"], item["target"]["y"], item["ball"]["x"], item["ball"]["y"]]
        for t in item["t"]:
            xt.append(t["x"])
            xt.append(t["y"])
        for o in item["o"]:
            xt.append(o["x"])
            xt.append(o["y"])
        x.append(xt)
    return torch.Tensor(x)
