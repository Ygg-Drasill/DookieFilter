import os.path

import torch

from model.player_predictor import PlayerPredictor


def load_model(model_path, device):
    if not os.path.isfile(model_path): raise FileNotFoundError

    model_name = os.path.basename(model_path)
    model = PlayerPredictor(device, *decode_model_params(model_name))
    model.load_state_dict(torch.load(model_path, weights_only=True, map_location=device))
    model.to(device), model.eval()

    return model


def decode_model_params(model_name: str) -> (int, int, int):
    """returns n_nearest_players, n_hidden, n_stack"""
    params_str = model_name.split("-")
    return int(params_str[0]), int(params_str[2]), int(params_str[1])
