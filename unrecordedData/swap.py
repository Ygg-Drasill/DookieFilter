import json
import matplotlib.pyplot as plt
from datetime import datetime

def read_jsonl(file_path):
    with open(file_path, "r", encoding="utf-8") as file:
        for line in file:
            yield json.loads(line)

if __name__ == "__main__":
    file_path = "../raw.jsonl"
    selected_number_1 = "9"   # The first player you want to focus on
    selected_number_2 = "19"  # The second player you want to focus on

    start_time = 1726507508000  # Start time in wallClock milliseconds
    end_time = start_time + 10000  # End time (10 seconds later)

    print(f"Start time (wallClock): {start_time}")  # Print the start wallClock time

    player_1_positions = []  # To store player 1's positions along with timestamps
    player_2_positions = []  # To store player 2's positions along with timestamps

    for entry in read_jsonl(file_path):
        for data_entry in entry.get("data", []):
            wall_clock = data_entry.get("wallClock", 0)
            if start_time <= wall_clock <= end_time:
                if "awayPlayers" in data_entry:  # Check only 'awayPlayers'
                    for player in data_entry["awayPlayers"]:
                        if player.get("number") == selected_number_1:
                            try:
                                _, y, _ = player["xyz"]
                                readable_time = datetime.utcfromtimestamp(wall_clock / 1000).strftime('%H:%M:%S')
                                player_1_positions.append((readable_time, y))  # Store formatted time and y-coordinate
                            except KeyError:
                                print(f"Missing 'xyz' data for player {selected_number_1}")
                        elif player.get("number") == selected_number_2:
                            try:
                                _, y, _ = player["xyz"]
                                readable_time = datetime.utcfromtimestamp(wall_clock / 1000).strftime('%H:%M:%S')
                                player_2_positions.append((readable_time, y))  # Store formatted time and y-coordinate
                            except KeyError:
                                print(f"Missing 'xyz' data for player {selected_number_2}")

    # Plotting the players' positions over time
    if player_1_positions or player_2_positions:
        # Prepare data for player 1
        if player_1_positions:
            times_1, y_positions_1 = zip(*player_1_positions)  # Unzip the list into time and y coordinates
            plt.plot(times_1, y_positions_1, label=f"Player {selected_number_1}")

        # Prepare data for player 2
        if player_2_positions:
            times_2, y_positions_2 = zip(*player_2_positions)  # Unzip the list into time and y coordinates
            plt.plot(times_2, y_positions_2, label=f"Player {selected_number_2}")

        plt.xlabel("Time")
        plt.ylabel("Y Coordinate")
        plt.title("Player Positions Over Time")
        plt.xticks(rotation=45)  # Rotate x-axis labels for better readability
        plt.legend()
        plt.show()
    else:
        print(f"No data found for players {selected_number_1} and {selected_number_2} in the specified time range.")
