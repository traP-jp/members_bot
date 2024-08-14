package handler

import (
	"context"
	"fmt"
	"log"
	"slices"

	"github.com/traPtitech/traq-ws-bot/payload"
)

// スタンプが押されたとき、招待を承認するか却下するか判定する
// :kan:が押されていたら何もしない
// スタンプを時系列で前から見ていき、指定されたスタンプが閾値以上押されていたら招待を送信する
func (h *BotHandler) AcceptOrReject(p *payload.BotMessageStampsUpdated) {
	ctx := context.Background()

	for _, stamp := range p.Stamps {
		if stamp.StampID == h.inactiveStampID && stamp.UserID == h.botUserID {
			return // :kan:押されてたら、何もしない
		}
	}

	adminIDs, err := h.traqClient.GetGroupMemberIDs(ctx, h.adminGroupID)
	if err != nil {
		log.Printf("failed to get group member IDs: %v", err)
		return
	}

	slices.SortFunc(p.Stamps, func(i, j payload.MessageStamp) int { return int(i.CreatedAt.Sub(j.CreatedAt)) })

	acceptStampCount := 0
	rejectStampCount := 0
	accept, reject := false, false
	for _, stamp := range p.Stamps {
		if !slices.Contains(adminIDs, stamp.UserID) {
			continue
		}
		if stamp.StampID == h.acceptStampID {
			acceptStampCount++
		} else if stamp.StampID == h.rejectStampID {
			rejectStampCount++
		}

		if acceptStampCount >= h.acceptStampThreshold {
			accept = true
			break
		}
		if rejectStampCount >= h.rejectStampThreshold {
			reject = true
			break
		}
	}

	if !accept && !reject {
		return
	}

	err = h.traqClient.AddStamp(ctx, p.MessageID, h.inactiveStampID, 1)
	if err != nil {
		log.Printf("failed to add stamp: %v", err)
		return
	}

	if reject {
		return
	}

	invitations, err := h.ir.GetInvitations(ctx, p.MessageID)
	if err != nil {
		log.Printf("failed to get invitations: %v", err)
		return
	}

	err = h.githubClient.SendInvitations(ctx, invitations)
	if err != nil {
		log.Printf("failed to send invitations: %v", err)
		return
	}

	message := "招待を送信しました。確認してください\n"
	for _, inv := range invitations {
		message += fmt.Sprintf("@%s (%s)\n", inv.TraqID(), inv.GitHubID())
	}

	_, err = h.traqClient.PostMessage(ctx, h.botChannelID, message)
	if err != nil {
		log.Printf("failed to post message: %v", err)
	}
}
