import pandas as pd
import matplotlib.pyplot as plt
import numpy as np

if __name__ == "__main__":
    # Load your data
    dataframe = pd.read_csv("../gym/data/hole.csv")
    xx = dataframe["a_19_x"]
    yy = dataframe["a_19_y"]

    # Create a mask for valid (non-null) points
    mask = ~(xx.isnull() | yy.isnull())
    valid_indices = np.where(mask)[0]

    # Calculate distances between consecutive valid points
    distances = []
    for i in range(len(valid_indices) - 1):
        start = valid_indices[i]
        end = valid_indices[i + 1]
        x1, y1 = xx[start], yy[start]
        x2, y2 = xx[end], yy[end]
        distance = np.sqrt((x2 - x1) ** 2 + (y2 - y1) ** 2)
        distances.append(distance)

    # Calculate average distance and threshold for "too much space"
    if distances:  # only if we have valid distances
        avg_distance = np.mean(distances)
        threshold_distance = 2 * avg_distance  # adjust multiplier as needed
    else:
        threshold_distance = float('inf')  # no valid distances

    # Create figure
    plt.figure(figsize=(12, 6))

    # Plot all points first (we'll handle connections separately)
    plt.scatter(xx[mask], yy[mask], color='blue', label='Position data points', zorder=3)

    # Plot connections between points with appropriate styles
    for i in range(len(valid_indices) - 1):
        start = valid_indices[i]
        end = valid_indices[i + 1]

        # Calculate distance between these points
        x1, y1 = xx[start], yy[start]
        x2, y2 = xx[end], yy[end]
        current_distance = np.sqrt((x2 - x1) ** 2 + (y2 - y1) ** 2)

        # Determine line style based on conditions
        if end == start + 1 and current_distance <= threshold_distance:
            # Solid line for normal consecutive points
            plt.plot([x1, x2], [y1, y2], 'b-', alpha=0.7, linewidth=1.5, zorder=2)
        else:
            # Dashed line for gaps or large distances
            plt.plot([x1, x2], [y1, y2], 'b--', alpha=0.3,
                     dashes=(5, 2), linewidth=1,
                     label='Large gap/missing data' if i == 0 else "", zorder=1)

    # Customize the plot
    plt.title("Player 19 X and Y-Position with missing data", fontsize=20)
    plt.xlabel("X Position")
    plt.ylabel("Y Position")
    plt.grid(True, linestyle='--', alpha=0.5)

    # Simplify legend
    handles, labels = plt.gca().get_legend_handles_labels()
    by_label = dict(zip(labels, handles))  # remove duplicates
    plt.legend(by_label.values(), by_label.keys())

    plt.tight_layout()
    plt.show()