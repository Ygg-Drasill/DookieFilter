import json
import matplotlib.pyplot as plt
from datetime import datetime
import numpy as np
import pandas as pd


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
    ax = np.linspace(0, s, s)
    print(ax.shape, series.T.shape)

    # Scatter plot
    plt.scatter(*series.T, label='Player 19 Position')

    # Dashed line connecting the dots
    plt.plot(*series.T, linestyle='--', color='gray', alpha=0.5, label='Path')

    plt.title(f"Player 19 X and Y-Position Over Time", fontsize=20)
    plt.grid()
    plt.legend()  # Show legend to differentiate between scatter and line
    plt.show()

    exit(0)
