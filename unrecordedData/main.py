import json
import matplotlib.pyplot as plt

def read_jsonl(file_path):
    with open(file_path, "r", encoding="utf-8") as file:
        for line in file:
            yield json.loads(line)

if __name__ == "__main__":
    file_path = "../visunator/raw.jsonl"
    selected_numbers = {"5"}  # Players to filter by (as strings to match JSON format)
    filtered_players = []

    for entry in read_jsonl(file_path):
        for key in ["homePlayers", "awayPlayers"]:
            if key in entry.get("data", [{}])[0]:
                for player in entry["data"][0][key]:
                    if player.get("number") in selected_numbers:
                        filtered_players.append(player)

    # Plotting
    plt.figure(figsize=(8, 6))
    for player in filtered_players:
        try:
            x, y, _ = player["xyz"]
            plt.scatter(x, y, label=f"Player {player['number']}")
        except KeyError:
            print(f"Missing 'xyz' data for player {player['number']}")

    plt.xlabel("X Coordinate")
    plt.ylabel("Y Coordinate")
    plt.title("Player Positions")
    plt.legend()
    plt.grid()
    plt.show()
