import pickle
import zmq


# dictionary to store cached file content
hole_cache = []


class CacheHoles():
