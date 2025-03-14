from typing import List, Union
from dataclasses import dataclass

@dataclass
class Player:
    number: str
    opta_id: str
    player_id: str
    speed: float
    xyz: List[float]

@dataclass
class Ball:
    speed: float
    xyz: List[float]

@dataclass
class Frame:
    away_players: List[Player]
    ball: Ball
    frame_idx: int
    game_clock: float
    home_players: List[Player]
    period: int
    wall_clock: int

@dataclass
class Signal:
    end_frame_idx: int
    end_wall_clock: int
    number: int
    start_frame_idx: int
    start_wall_clock: int

@dataclass
class GamePacket:
    league: str
    game_id: str
    feed_name: str
    message_id: str
    data: List[Union[Frame, Signal]]
