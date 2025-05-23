from threading import Thread

import torch
import zmq

from model.player_predictor import PlayerPredictor
from service.imputation_service import ImputationService
from service.imputation_service import sequence_to_tensor

test_response = [{
    "target": {"x": 0, "y": 0},
    "ball": {"x": 0, "y": 0},
    "t": [{"x": 0, "y": 0}, {"x": 0, "y": 0}],
    "o": [{"x": 0, "y": 0}, {"x": 0, "y": 0}],
},{
    "target": {"x": 1, "y": 1},
    "ball": {"x": 1, "y": 1},
    "t": [{"x": 1, "y": 1}, {"x": 1, "y": 1}],
    "o": [{"x": 1, "y": 1}, {"x": 1, "y": 1}],
},{
    "target": {"x": 2, "y": 2},
    "ball": {"x": 2, "y": 2},
    "t": [{"x": 2, "y": 2}, {"x": 2, "y": 2}],
    "o": [{"x": 2, "y": 2}, {"x": 2, "y": 2}],
}]

test_request = [{
    "frame_idx": 420,
    "player_number": 42,
}]

def test_service():
    device = torch.device("cuda:0" if torch.cuda.is_available() else "cpu")
    model = PlayerPredictor(device, 2, 4, 1)
    service = ImputationService(model, sequence_length=3)
    ctx = service.context
    service.socket_imputation.bind(f"inproc://imputation")
    service.socket_storage_api.connect(f"inproc://storage_api")
    service.socket_storage.connect(f"inproc://storage")

    thread = Thread(target=service.run)
    thread.start()

    socket_imputation = ctx.socket(zmq.PUSH)
    socket_storage_api = ctx.socket(zmq.REP)
    socket_storage = ctx.socket(zmq.PULL)

    socket_imputation.connect(f"inproc://imputation")
    socket_storage_api.bind(f"inproc://storage_api")
    socket_storage.bind(f"inproc://storage")

    socket_imputation.send_json(test_request)
    topic = socket_storage_api.recv_string(zmq.SNDMORE)
    assert topic == "frameRangeNearest"
    message = socket_storage_api.recv_json()
    socket_storage_api.send_json(test_response)

    topic = socket_storage.recv_string(zmq.SNDMORE)
    assert topic == "player"
    prediction = socket_storage.recv_json()
    assert prediction["frame_idx"] == test_request[0]["frame_idx"]
    assert prediction["player_number"] == test_request[0]["player_number"]

    socket_imputation.send_json(["kill"])

def test_sequence_to_tensor():
    t = sequence_to_tensor(test_response)
    assert t.shape == (len(test_response), 4 + len(test_response[0]["t"]*2*2))
