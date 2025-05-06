import zmq

if __name__ == '__main__':
    context = zmq.Context()
    socket = context.socket(zmq.PUSH)
    socket.connect("tcp://localhost:5555")
    req = [{"num": 11, "idx": 42}]
    socket.send_json(req)
