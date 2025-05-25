import sys
import torch

from service.load_model import load_model
from service.imputation_service import ImputationService

if __name__ == "__main__":
    device = torch.device("cuda:0" if torch.cuda.is_available() else "cpu")
    model = load_model(sys.argv[1], device)
    service = ImputationService(model)
    service.connect(5557, 5555, 5558)
    service.run()
