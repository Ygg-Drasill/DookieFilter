import io
import os

import matplotlib.pyplot as plt
import numpy as np
import pandas as pd
import torch
from matplotlib.figure import Figure

from gym.utils.data import get_features
from model.player_predictor import PlayerPredictor

def load_model(model_path) -> PlayerPredictor:
    device = torch.device("cuda:0" if torch.cuda.is_available() else "cpu")
    model = PlayerPredictor(device, 3, 64, 4)
    model.load_state_dict(torch.load(model_path, weights_only=True, map_location=device))
    model.to(device)
    model.eval()
    return model

def test_model(model: PlayerPredictor, test_data_path: str) -> Figure:
    device = torch.device("cuda:0" if torch.cuda.is_available() else "cpu")
    model.eval()

    dataframe = pd.read_csv(test_data_path, low_memory=True)
    d: pd.DataFrame
    player_key = "h_6"
    start, end = 0, 500
    prediction_start = 50

    player_truth, player_prediction, ball = [], [], []
    sequence = []
    for dataframe_index in range(start, end):
        features = get_features(dataframe.iloc[dataframe_index], model.n_nearest_players, player_key)
        sequence.append(features)
        feature_coords = features.reshape(-1, 2)
        player_truth.append(feature_coords[0])
        player_prediction.append(feature_coords[0])

    sequence = np.array(sequence)
    input_length = 20

    for i in range(prediction_start-start, len(player_prediction)):
        input_tensor = torch.from_numpy(sequence[i-input_length:i].reshape(1,input_length, model.input_size)).float().to(device)
        out = model(input_tensor).cpu().detach().numpy()
        player_prediction[i] = out.squeeze()
        if i < len(sequence)-1:
            sequence[i+1][0] = out.squeeze()[0]
            sequence[i+1][1] = out.squeeze()[1]

    fig = plt.figure()
    plt.plot(*np.array(player_truth).T)
    plt.plot(*np.array(player_prediction).reshape(end-start, 2).T)

    plt.legend(["truth", "prediction", "ball"])
    return fig
