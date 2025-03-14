import os

import matplotlib.pyplot as plt
import numpy as np
import pandas as pd
import torch
import sys
from model import PlayerPredictor
from utils.data import get_features

if __name__ == "__main__":
    model_path = os.path.abspath("./model.pth")
    device = torch.device("cuda:0" if torch.cuda.is_available() else "cpu")
    model = PlayerPredictor(device, 3, 32, 4)
    model.load_state_dict(torch.load(model_path, weights_only=True))
    model.to(device)
    model.eval()

    dataframe = pd.read_csv("../data/chunk_0.csv")
    d: pd.DataFrame
    player_key = "h_10"
    start, end = 2500, 3500
    prediction_start = 3000

    player_truth, player_prediction, ball = [], [], []
    sequence = []
    for dataframe_index in range(start, end):
        features = get_features(dataframe.iloc[dataframe_index], 3, player_key)
        sequence.append(features)
        feature_coords = features.reshape(-1, 2)
        player_truth.append(feature_coords[0])
        player_prediction.append(feature_coords[0])
        ball.append(feature_coords[1])

    sequence = np.array(sequence)
    input_length = 20

    for i in range(prediction_start-start, len(player_prediction)):
        input_tensor = torch.from_numpy(sequence[i-input_length:i].reshape(1,input_length, 16)).float().to(device)
        out = model(input_tensor).cpu().detach().numpy()
        player_prediction[i] = out.squeeze()

    plt.plot(*np.array(player_truth).T)
    plt.plot(*np.array(player_prediction).reshape(end-start, 2).T)
    plt.plot(*np.array(ball).T)

    plt.legend(["truth", "prediction", "ball"])
    plt.show()



