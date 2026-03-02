package bumble

import (
	"context"
	"fmt"
	"time"

	"github.com/dvcrn/matrix-bumble/bumble/bumbletypes"
)

func (bum *BumbleConnector) cachedClientOpenChat(ctx context.Context, bumbleClient *BumbleClient, bumbleId string) (*bumbletypes.ClientOpenChat, error) {
	// check if we have already a clientOpenChat result for the given bumbleId
	if cached, ok := bum.cache.Get("clientOpenChat_" + bumbleId); ok {
		fmt.Println("cached clientOpenChat result found for bumbleId: ", bumbleId)
		// conver to correct type
		return cached.(*bumbletypes.ClientOpenChat), nil
	}

	fmt.Println("NO cached clientOpenChat result found for bumbleId: ", bumbleId)

	// if not, fetch it
	openChat, err := bumbleClient.ClientOpenChat(bumbleId, 0)
	if err != nil {
		fmt.Printf("error fetching chat instances: %v\n", err)
		return nil, err
	}

	// cache the result for 5 min
	bum.cache.Set("clientOpenChat_"+bumbleId, openChat, 5*time.Minute)
	return openChat, nil
}
