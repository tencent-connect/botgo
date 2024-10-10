package main

import (
	"context"
	"log"

	"github.com/tencent-connect/botgo/dto"
)

func (p Processor) setEmoji(ctx context.Context, channelID string, messageID string) {
	err := p.api.CreateMessageReaction(
		ctx, channelID, messageID, dto.Emoji{
			ID:   "307",
			Type: 1,
		},
	)
	if err != nil {
		log.Println(err)
	}
}

func (p Processor) setPins(ctx context.Context, channelID, msgID string) {
	_, err := p.api.AddPins(ctx, channelID, msgID)
	if err != nil {
		log.Println(err)
	}
}

func (p Processor) setAnnounces(ctx context.Context, data *dto.WSATMessageData) {
	if _, err := p.api.CreateChannelAnnounces(
		ctx, data.ChannelID,
		&dto.ChannelAnnouncesToCreate{MessageID: data.ID},
	); err != nil {
		log.Println(err)
	}
}

func (p Processor) sendChannelReply(ctx context.Context, channelID string, toCreate *dto.MessageToCreate) error {
	if _, err := p.api.PostMessage(ctx, channelID, toCreate); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (p Processor) sendGroupReply(ctx context.Context, groupID string, toCreate dto.APIMessage) error {
	log.Printf("EVENT ID:%v", toCreate.GetEventID())
	if _, err := p.api.PostGroupMessage(ctx, groupID, toCreate); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (p Processor) sendC2CReply(ctx context.Context, userID string, toCreate dto.APIMessage) error {
	log.Printf("EVENT ID:%v", toCreate.GetEventID())
	if _, err := p.api.PostC2CMessage(ctx, userID, toCreate); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
