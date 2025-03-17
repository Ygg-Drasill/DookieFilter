from torch.utils.tensorboard import SummaryWriter

class BoardLogger:
    def __init__(self, summary_writer: SummaryWriter):
        self.step = 0
        self.summary_writer = summary_writer

    def log(self, tag, value):
        self.summary_writer.add_scalar(tag, value, self.step)
        self.step += 1
