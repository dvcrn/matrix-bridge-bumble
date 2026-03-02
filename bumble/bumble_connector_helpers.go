package bumble

import (
	"context"
	"fmt"
	"time"

	"github.com/dvcrn/matrix-bumble/repo"
)

func (bum *BumbleConnector) LockUserRoom(ctx context.Context, bumUser *repo.BumbleUser) error {
	if bumUser.RoomLockedAt != nil {
		fmt.Println("room already locked, skipping")
		return nil
	}

	fmt.Println("marking room as locked", bumUser.Name)
	_, err := bum.kit.MarkRoomReadOnly(ctx, bumUser.ToMatrixRoom())
	if err != nil {
		return err
	}

	now := time.Now()
	bumUser.RoomLockedAt = &now
	return repo.UpdateBumbleUser(bum.gdb, bumUser)
}

func (bum *BumbleConnector) UnlockUserRoom(ctx context.Context, bumUser *repo.BumbleUser) error {
	if bumUser.RoomLockedAt == nil {
		fmt.Println("room isn't locked, skipping")
		return nil
	}

	fmt.Println("marking room as unlocked", bumUser.Name)
	_, err := bum.kit.ResetRoomPermission(ctx, bumUser.ToMatrixRoom())
	if err != nil {
		return err
	}

	bumUser.RoomLockedAt = nil
	return repo.UpdateBumbleUser(bum.gdb, bumUser)
}
