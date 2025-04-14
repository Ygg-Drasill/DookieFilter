import json
import matplotlib.pyplot as plt
from datetime import datetime

def read_jsonl(file_path):
    with open(file_path, "r", encoding="utf-8") as file:
        for line in file:
            yield json.loads(line)

if __name__ == "__main__":
    file_path = "../visunator/raw.jsonl"
    selected_number = "9"  # The player you want to focus on

    start_time = 1726507990000  # Start time in wallClock milliseconds
    end_time = start_time + 10000  # End time (10 seconds later)

    print(f"Start time (wallClock): {start_time}")  # Print the start wallClock time

    player_positions = []  # To store player positions

    for entry in read_jsonl(file_path):
        for data_entry in entry.get("data", []):
            wall_clock = data_entry.get("wallClock", 0)
            if start_time <= wall_clock <= end_time:
                for key in ["homePlayers", "awayPlayers"]:
                    if key in data_entry:
                        for player in data_entry[key]:
                            if player.get("number") == selected_number:
                                try:
                                    x, y, _ = player["xyz"]
                                    player_positions.append((x, y))  # Store x and y coordinates
                                except KeyError:
                                    print(f"Missing 'xyz' data for player {selected_number}")

    # Plotting the player's path
    if player_positions:
        x_positions, y_positions = zip(*player_positions)  # Unzip the list into x and y coordinates
        step = 5  # Show a dot every 5 points
        plt.plot(x_positions, y_positions, '--', color='blue', alpha=0.7)  # Dashed line for the path
        plt.plot(x_positions[::step], y_positions[::step], 'o', color='blue', markersize=4, label=f"Player {selected_number}")  # Dots every 5 points
        plt.xlabel("X Position")
        plt.ylabel("Y Position")
        plt.title(f"Player {selected_number} Path", fontsize=20)
        plt.legend()
        plt.grid(True)
        plt.show()
    else:
        print(f"No data found for player {selected_number} in the specified time range.")
