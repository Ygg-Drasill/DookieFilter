import torch
from torch.utils.data import DataLoader
from tqdm import tqdm
from torch import nn

from gym.board_logger import BoardLogger
from gym.utils.data import denormalize_x, m_to_cm, denormalize_y

progress_bar_columns = 200

LOSS_MOVEMENT_SCALE = 0.01
LOSS_ANGLE_SCALE = 0.01

MIN_MOVEMENT_THRESHOLD = 0.5


def train_epoch(epoch: int,
                epoch_total: int,
                model: nn.Module,
                dataloader: DataLoader,
                loss_function,
                optimizer:torch.optim.Optimizer,
                device:torch.device,
                board_logger: BoardLogger = None):

    model.train(True)
    torch.set_grad_enabled(True)
    progress_dataloader = tqdm(dataloader,
                               ncols=progress_bar_columns,
                               desc=f"Training Epoch {epoch + 1}/{epoch_total}",
                               unit="batch")

    running_loss = 0.0
    total = 0
    for batch_index, batch in enumerate(progress_dataloader):
        batch_x: torch.Tensor
        batch_y: torch.Tensor
        batch_x, batch_y = batch[0].to(device), batch[1].to(device)
        if torch.isnan(batch_x).any() or torch.isnan(batch_y).any():
            continue

        output, delta = model(batch_x)

        low_movement_penalty = torch.clamp(MIN_MOVEMENT_THRESHOLD - delta.abs(), min=0.0).mean()

        prev_delta = batch_x[:, -1, 0: 2] - batch_x[:, -2, 0: 2]
        cos_similarity = torch.nn.functional.cosine_similarity(prev_delta, delta)
        angle_penalty = 1 - cos_similarity.mean()

        out_x = output[:, 0]
        out_y = output[:, 1]
        output = torch.stack((denormalize_x(out_x), denormalize_y(out_y)), dim=1)

        loss = torch.functional.cdist(output, batch_y, p=2)# loss_function(output, batch_y)
        #loss = loss + low_movement_penalty * LOSS_MOVEMENT_SCALE + angle_penalty * LOSS_ANGLE_SCALE

        running_loss += loss.item()
        total += 1
        optimizer.zero_grad()
        loss.backward()
        torch.nn.utils.clip_grad_norm_(model.parameters(), max_norm=1.0)
        optimizer.step()

        avg_loss_across_batches = running_loss / total
        progress_dataloader.set_postfix({'loss': avg_loss_across_batches})
        if board_logger is not None:
            board_logger.log("Loss/train", loss)
    return running_loss / total

def validate_epoch(epoch: int,
                   epoch_total: int,
                   model: nn.Module,
                   dataloader: DataLoader,
                   loss_function,
                   device:torch.device,
                   board_logger: BoardLogger = None):
    progress_dataloader = tqdm(dataloader,
                               ncols=progress_bar_columns,
                               desc=f"Validating Epoch {epoch + 1}/{epoch_total}",
                               unit="batch")

    running_loss = 0.0
    total = 0
    for batch_index, batch in enumerate(progress_dataloader):
        batch_x: torch.Tensor
        batch_y: torch.Tensor
        batch_x, batch_y = batch[0].to(device), batch[1].to(device)

        #if torch.isnan(batch_x).any() or torch.isnan(batch_y).any():
        #    print(batch_index, batch_x, batch_y)
        with torch.no_grad():
            output, _ = model(batch_x)
            out_x = output[:, 0]
            out_y = output[:, 1]
            output = torch.stack((m_to_cm(denormalize_x(out_x)), m_to_cm(denormalize_y(out_y))), dim=1)

            truth_x = batch_y[:, 0]
            truth_y = batch_y[:, 1]
            truth = torch.stack((m_to_cm(truth_x), m_to_cm(truth_y)), dim=1)
            loss = torch.nn.functional.l1_loss(output, truth)
            running_loss += loss.item()
            total += 1
        avg_loss_across_batches = running_loss / total
        progress_dataloader.set_postfix({'avg deviation': avg_loss_across_batches})
        if board_logger is not None:
            board_logger.log("Average Deviation (cm)", loss)
    return running_loss / total
