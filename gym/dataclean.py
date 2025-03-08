import json
import os
import sys
import numpy as np
import matplotlib.pyplot as plt

data_path = sys.argv[1]
output_target = sys.argv[2]

file = open(data_path)
data = {
    "ball": []
}

frame_idx = []
game_active = False
done = False

init = False
for line in file:
    packet = json.loads(line)
    if not init:
        init = True
        os.mkdir(output_target + "/" + packet["gameId"])

    for seperated_frame in packet["data"]:
        if "frameIdx" not in dict.keys(seperated_frame):
            print("signal")
            if game_active:
                done = True
                break
            game_active = not game_active
            continue
        if len(seperated_frame["homePlayers"]) == 0 or len(seperated_frame["awayPlayers"]) == 0:
            continue
        if not game_active:
            continue

        idx = seperated_frame["frameIdx"]
        frame_idx.append(idx)
        home_players = seperated_frame["homePlayers"]
        away_players = seperated_frame["awayPlayers"]

        for player in home_players:
            coords = np.array([player["xyz"][0], player["xyz"][1]])
            number = player["number"]
            player_key = "h" + str(number)
            if player_key not in dict.keys(data):
                data[player_key] = [coords]
            else:
                data[player_key].append(coords)

        for player in away_players:
            coords = np.array([player["xyz"][0], player["xyz"][1]])
            number = player["number"]
            player_key = "a" + str(number)
            if player_key not in dict.keys(data):
                data[player_key] = [coords]
            else:
                data[player_key].append(coords)

    if done:
        break

plt.plot(*np.array(data["a9"]).T)
plt.show()
