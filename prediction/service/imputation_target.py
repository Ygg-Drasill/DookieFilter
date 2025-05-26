

class ImputationTarget:
    def __init__(self, player_number: int, home: bool, frame_index: int):
        self.player_number = player_number
        self.frame_index = frame_index
        self.home = home

    def step(self):
        self.frame_index += 1