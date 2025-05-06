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

    player_positions = []  # To store player positions along with timestamps

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
                                    # Convert wallClock to a readable time format
                                    readable_time = datetime.utcfromtimestamp(wall_clock / 1000).strftime('%H:%M:%S')
                                    player_positions.append((readable_time, y))  # Store formatted time and y-coordinate
                                except KeyError:
                                    print(f"Missing 'xyz' data for player {selected_number}")

    # Plotting the player's positions over time
    if player_positions:
        times, y_positions = zip(*player_positions)  # Unzip the list into time and y coordinates
        plt.plot(times, y_positions, label=f"Player {selected_number}")
        plt.xlabel("Time")
        plt.ylabel("Y Coordinate")
        plt.title(f"Player {selected_number} Y-Position Over Time")
        plt.xticks(rotation=45)  # Rotate x-axis labels for better readability
        plt.legend()
        plt.grid()
        plt.show()
    else:
        print(f"No data found for player {selected_number} in the specified time range.")
