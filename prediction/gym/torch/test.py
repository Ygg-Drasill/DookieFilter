import os

import matplotlib.pyplot as plt
import numpy as np
import pandas as pd
import torch
from tqdm import tqdm

from gym.utils.data import get_features
from model.player_predictor import PlayerPredictor


if __name__ == "__main__":
    n_nearest_players = 4
    model_path = os.path.abspath("./model.pth")#"../../model/out/model.pth")
    device = torch.device("cuda:0" if torch.cuda.is_available() else "cpu")
    model = PlayerPredictor(device, n_nearest_players, 128, 8)
    print(sum(p.numel() for p in model.parameters()))
    model.load_state_dict(torch.load(model_path, weights_only=True, map_location=device))
    model.to(device)
    model.eval()

    dataframe = pd.read_csv("../data/chunk_0.csv", low_memory=True)
    d: pd.DataFrame
    player_key = "h_14"
    start, end = 2500, 3500
    prediction_start = 3000

    player_truth, player_prediction, ball = [], [], []
    sequence = []
    for dataframe_index in range(start, end):
        features = get_features(dataframe.iloc[dataframe_index], n_nearest_players, player_key)
        sequence.append(features)
        feature_coords = features.reshape(-1, 2)
        player_truth.append(feature_coords[0])
        player_prediction.append(feature_coords[0])

    sequence = np.array(sequence)
    input_length = 40

    for i in tqdm(range(prediction_start-start, len(player_prediction))):
        input_tensor = torch.from_numpy(sequence[i-input_length:i].reshape(1,input_length, 20)).float().to(device)
        out = model(input_tensor).cpu().detach().numpy()
        player_prediction[i] = out.squeeze()
        if i < len(sequence)-1:
            sequence[i+1][0] = out.squeeze()[0]
            sequence[i+1][1] = out.squeeze()[1]

    plt.plot(*np.array(player_truth).T)
    plt.plot(*np.array(player_prediction).reshape(end-start, 2).T)

    plt.legend(["truth", "prediction"])
    plt.show()



