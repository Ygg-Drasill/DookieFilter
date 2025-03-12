import torch
from torch.utils.data import DataLoader
from tqdm import tqdm
from torch import nn

progress_bar_columns = 200

def train_epoch(epoch: int,
                epoch_total: int,
                model: nn.Module,
                dataloader: DataLoader,
                loss_function,
                optimizer:torch.optim.Optimizer,
                device:torch.device):

    model.train(True)
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

        output = model(batch_x)
        loss = loss_function(output, batch_y);
        running_loss += loss.item()
        total += 1
        optimizer.zero_grad()
        loss.backward()
        torch.nn.utils.clip_grad_norm_(model.parameters(), max_norm=1.0)
        optimizer.step()

        avg_loss_across_batches = running_loss / total
        progress_dataloader.set_postfix({'loss': avg_loss_across_batches})
    return running_loss / total

def validate_epoch(epoch: int,
                   epoch_total: int,
                   model: nn.Module,
                   dataloader: DataLoader,
                   loss_function,
                   device:torch.device):

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
            output = model(batch_x)
            if torch.isnan(output).any() or torch.isnan(batch_y).any():
                continue
            loss = loss_function(output, batch_y);
            running_loss += loss.item()
            total += 1

        avg_loss_across_batches = running_loss / total
        progress_dataloader.set_postfix({'loss': avg_loss_across_batches})
    return running_loss / total
