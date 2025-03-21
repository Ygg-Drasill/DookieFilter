import json
import matplotlib.pyplot as plt
from datetime import datetime
import pandas as pd
import numpy as np

def read_jsonl(file_path):
    with open(file_path, "r", encoding="utf-8") as file:
        for line in file:
            yield json.loads(line)

if __name__ == "__main__":
    dataframe = pd.read_csv("../gym/data/hole.csv")

    xx = dataframe.pop("a_19_x")
    yy = dataframe.pop("a_19_y")

    series = np.column_stack([xx, yy])
    s = len(series)
    ax = np.linspace(0, s,s)
    print(ax.shape, series.T.shape)
    plt.scatter(*series.T)
    plt.title(f"Player 19 X and Y-Position Over Time")
    plt.grid()
    plt.show()





    exit(0)