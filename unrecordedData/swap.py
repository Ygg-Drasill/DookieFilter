import json
import matplotlib.pyplot as plt
from scipy.interpolate import make_interp_spline

def read_jsonl(file_path):
    with open(file_path, "r", encoding="utf-8") as file:
        for line in file:
            yield json.loads(line)

def print_player_location(player_number, wall_clock1, y):
    readable_time = wall_clock1
    print(f"Player {player_number} at {readable_time}: Y = {y}")

def find_swap_time(times_1, y_positions_1, times_2, y_positions_2):
    for i in range(min(len(times_1), len(times_2)) - 1):
        if (y_positions_1[i] - y_positions_2[i]) * (y_positions_1[i+1] - y_positions_2[i+1]) < 0:
            return times_1[i]
    return None

def plot_smoothed_curve_with_dots(times, positions, label):
    if len(times) > 2:
        spline = make_interp_spline(times, positions, k=3)
        # Increase the step size to create more space between dots
        step_size = 50  # Adjust this value to control the spacing between dots
        smoothed_times = range(min(times), max(times), step_size)
        smoothed_positions = spline(smoothed_times)
        plt.plot(smoothed_times, smoothed_positions, label=label, marker='o', linestyle='--', markersize=5)

if __name__ == "__main__":
    file_path = "../raw.jsonl"
    selected_number_1 = "9"
    selected_number_2 = "19"

    initial_start_time = 1726507508000
    initial_end_time = initial_start_time + 10000

    print(f"Initial start time (wallClock): {initial_start_time}")

    player_1_positions = []
    player_2_positions = []

    for entry in read_jsonl(file_path):
        for data_entry in entry.get("data", []):
            wall_clock = data_entry.get("wallClock", 0)
            if initial_start_time <= wall_clock <= initial_end_time:
                if "awayPlayers" in data_entry:
                    for player in data_entry["awayPlayers"]:
                        if player.get("number") in [selected_number_1, selected_number_2]:
                            try:
                                _, y, _ = player["xyz"]
                                print_player_location(player.get("number"), wall_clock, y)
                                if player.get("number") == selected_number_1:
                                    player_1_positions.append((wall_clock, y))
                                else:
                                    player_2_positions.append((wall_clock, y))
                            except KeyError:
                                print(f"Missing 'xyz' data for player {player.get('number')}")

    if player_1_positions and player_2_positions:
        player_1_positions.sort()
        player_2_positions.sort()
        times_1, y_positions_1 = zip(*player_1_positions)
        times_2, y_positions_2 = zip(*player_2_positions)

        swap_time = find_swap_time(times_1, y_positions_1, times_2, y_positions_2)
        if swap_time:
            # Adjust the time range to focus on the swap time
            window_size = 500  # Adjust this value to control the time window around the swap
            start_time = swap_time - window_size
            end_time = swap_time + window_size

            # Filter positions within the new time range
            player_1_positions_focused = [(t, y) for t, y in player_1_positions if start_time <= t <= end_time]
            player_2_positions_focused = [(t, y) for t, y in player_2_positions if start_time <= t <= end_time]

            if player_1_positions_focused and player_2_positions_focused:
                times_1_focused, y_positions_1_focused = zip(*player_1_positions_focused)
                times_2_focused, y_positions_2_focused = zip(*player_2_positions_focused)

                plt.figure(figsize=(10, 6))
                plot_smoothed_curve_with_dots(times_1_focused, y_positions_1_focused, f"Player {selected_number_1}")
                plot_smoothed_curve_with_dots(times_2_focused, y_positions_2_focused, f"Player {selected_number_2}")

                plt.xlabel("Time (wallClock)")
                plt.ylabel("Y Coordinate")
                plt.title(f"Swapping Player Positions", fontsize=28)
                plt.legend()
                plt.grid(True)
                plt.show()
            else:
                print("No data found in the focused time range around the swap time.")
        else:
            print("No swap time found.")
    else:
        print(f"No data found for players {selected_number_1} and {selected_number_2} in the specified time range.")
