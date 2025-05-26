import sys
import torch
import zmq

from model.player_predictor import PlayerPredictor
from service.load_model import load_model
from service.imputation_service import ImputationService

if __name__ == "__main__":
    device = torch.device("cuda:0" if torch.cuda.is_available() else "cpu")
    torch.set_grad_enabled(False)
    model = PlayerPredictor(device, 2, 4, 2) #load_model(sys.argv[1], device)
    service = ImputationService(model)
    service.connect(5557, 5555, 5558)
    service.socket_storage.send_string("position", zmq.SNDMORE)
    service.socket_storage.send_json({
        "frameIdx": 1,
        "number": 2,
        "home": True,
        "x": 42.2,
        "y": 69.9,
    })
    service.run()
