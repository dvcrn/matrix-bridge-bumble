package bumble

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dvcrn/matrix-bridgekit/matrix"
	"github.com/dvcrn/matrix-bumble/repo"
	"maunium.net/go/mautrix/format"
)

func (bum *BumbleConnector) GenerateAIWelcomeMessage(ctx context.Context, user *matrix.User, bumbleUser *repo.BumbleUser) error {
	// TODO: make this less costly
	return nil

	userState, ok := bum.userState[user.MXID]
	if !ok {
		fmt.Println("could not find userstate for user ", user.DisplayName)
		return errors.New("could not find userstate")
	}

	aiWelcomeMessageEventId := bumbleUser.AIWelcomeMessageEventID
	fmt.Println("--- aiWelcomeMessageEventId: ", aiWelcomeMessageEventId)

	if aiWelcomeMessageEventId == "" && bumbleUser.AIWelcomeMessageSentAt != nil {
		fmt.Println("aiWelcomeMessageEventId is empty")
		return nil
	}

	// fetch user info
	getUserResponse, err := userState.Bumble.GetUser(bumbleUser.BumbleID)
	if err != nil {
		fmt.Printf("error fetching user: %v\n", err)
		return err
	}

	buser := getUserResponse
	output, err := bum.ai.GenWelcomePrompt(ctx, userState.UserProfile.DumpProfile(), buser.DumpProfile())
	if err != nil {
		fmt.Println("error generating output: ", err)
		return err
	}

	fmt.Println("generated output: ", output)

	content := format.RenderMarkdown(output, true, true)

	_, err = bum.kit.SendBotMessageInRoom(ctx, bumbleUser.ToMatrixRoom(), &content)
	if err != nil {
		fmt.Println("error sending bot message: ", err)
		return err
	}

	now := time.Now()
	bumbleUser.AIWelcomeMessageSentAt = &now
	return repo.UpdateBumbleUser(bum.gdb, bumbleUser)
}

func (bum *BumbleConnector) ReconcileAIWelcomeMessages(ctx context.Context, u *matrix.User) error {
	users, err := repo.GetBumbleUsersWithUnprocessedAIMessage(bum.gdb)
	if err != nil {
		return err
	}

	for _, user := range users {
		bum.queue.Enqueue("gen_ai_welcome_msg"+user.Name, func() {
			if err := bum.GenerateAIWelcomeMessage(ctx, u, user); err != nil {
				fmt.Println("error generating ai welcome message from reconcile: ", err)
			}
		})
	}

	return nil
}
