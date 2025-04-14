import json
import matplotlib.pyplot as plt
import matplotlib.dates as mdates
from matplotlib.ticker import MaxNLocator
from datetime import datetime

def read_jsonl(file_path):
    with open(file_path, "r", encoding="utf-8") as file:
        for line in file:
            yield json.loads(line)

def print_player_location(player_number, wall_clock1, y):
    readable_time = wall_clock1
    print(f"Player {player_number} at {readable_time}: Y = {y}")


if __name__ == "__main__":
    file_path = "../raw.jsonl"
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
                                    print_player_location(player.get("number"), wall_clock, y)
                                    # Convert wallClock to a datetime object for better plotting
                                    time_obj = datetime.utcfromtimestamp(wall_clock / 1000)
                                    player_positions.append((time_obj, y))  # Store datetime object and y-coordinate
                                except KeyError:
                                    print(f"Missing 'xyz' data for player {selected_number}")

    # Plotting the player's positions over time
    if player_positions:
        times, y_positions = zip(*player_positions)  # Unzip the list into time and y coordinates

        # Create the plot
        plt.figure(figsize=(10, 6))  # Adjust figure size for better readability
        plt.scatter(times, y_positions, label=f"Player {selected_number}", color='blue', marker='o')

        # Format the x-axis for better readability
        ax = plt.gca()
        ax.xaxis.set_major_formatter(mdates.DateFormatter('%H:%M:%S'))  # Format time as HH:MM:SS
        ax.xaxis.set_major_locator(MaxNLocator(nbins=6))  # Limit the number of x-axis labels
        plt.xticks(rotation=45)  # Rotate x-axis labels for better readability

        # Add labels and title
        plt.xlabel("Time (HH:MM:SS)")
        plt.ylabel("Y Coordinate")
        plt.title(f"Player {selected_number} Y-Position Over Time")
        plt.legend()
        plt.grid()
        plt.tight_layout()  # Adjust layout to prevent label overlap
        plt.show()
    else:
        print(f"No data found for player {selected_number} in the specified time range.")
