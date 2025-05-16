import zmq
from numpy import ndarray

from gym.utils.data import get_nearest_players
from model.player_predictor import PlayerPredictor

PROTOCOL = "tcp"
ADDRESS = "localhost"

class ImputationService:
    def __init__(self, model: PlayerPredictor):
        self.model = model
        self.context = zmq.Context()
        self.socket_storage_api = self.context.socket(zmq.REQ)
        self.socket_storage = self.context.socket(zmq.PUSH)
        self.socket_imputation = self.context.socket(zmq.PULL)
        self.player_series: dict[str, list[ndarray]] = dict[str, list[ndarray]]()

    def connect(self, imputation_port: int, storage_api_port: int, storage_port: int):
        self.socket_imputation.bind(f"{PROTOCOL}://{ADDRESS}:{imputation_port}")
        self.socket_storage_api.connect(f"{PROTOCOL}://{ADDRESS}:{storage_api_port}")
        self.socket_storage.connect(f"{PROTOCOL}://{ADDRESS}:{storage_port}")

    def close(self):
        self.socket_imputation.close()
        self.socket_storage_api.close()
        self.socket_storage.close()

    def run(self):
        while True: self.listen()

    def listen(self):
        imputation_request = self.socket_imputation.recv_json()
        for request in imputation_request:
            frame_idx = int(request["frame_idx"])
            player_number = int(request["player_number"])
            self.get_player_history(frame_idx, player_number)

        get_nearest_players()

    def get_player_history(self, frame_idx: int, player_number: int, nearest: int, target_player: int):
        self.socket_storage_api.send("frameRange")
        self.socket_storage_api.send_json({"startIdx": player_number, "endIdx": frame_idx, "n": nearest, "target": target_player})

    def store_player_frame(self):
        pass
