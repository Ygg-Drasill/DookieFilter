import json

import numpy
import torch
from matplotlib import pyplot as plt
import numpy as np
from numpy.ma.core import reshape

dataPath = "../raw.jsonl"

f = open(dataPath)
data = {
    "ball": []
}

frame_idx = []

frame_start =   100000
frame_end =     101000

for line in f:
    packet = json.loads(line)
    for f in packet["data"]:
        if "frameIdx" not in dict.keys(f):
            continue
        if len(f["homePlayers"]) == 0 or len(f["awayPlayers"]) == 0:
            continue
        idx = f["frameIdx"]
        if frame_start <= idx <= frame_end:
            frame_idx.append(f["frameIdx"])
            data["ball"].append(np.array(f["ball"]["xyz"], dtype=np.float32))

xplot = frame_idx
yplot = numpy.array(data["ball"])[:, :2]


fig, axs = plt.subplots(2)
axs[0].plot(*yplot.T, 'o-', linestyle='--', markersize=4)
axs[0].title.set_text("Path")
axs[1].plot(xplot, yplot, 'o-', linestyle='--', markersize=4)
axs[1].title.set_text("Coordinates")

plt.show()
print(len(data["ball"]))
