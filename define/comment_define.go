package define

import (
	pbGlobal "e.coding.net/mmstudio/blade/server/proto/global"
)

const (
	TopicType_Begin = iota
	TopicType_Hero  = iota - 1 // 英雄评论
	TopicType_Item             // 物品评论
	TopicType_End
)

// 发表元数据
type PublisherMetadata struct {
	PublisherId   int64  `json:"publisher_id" bson:"publisher_id"`     // 发表者玩家id
	PublisherName string `json:"publisher_name" bson:"publisher_name"` // 发表者名字
	ReplyToId     int64  `json:"reply_to_id" bson:"reply_to_id"`       // 回复**玩家id
	ReplyToName   string `json:"reply_to_name" bson:"reply_to_name"`   // 回复**玩家名字
	Content       string `json:"content" bson:"content"`               // 内容
	Thumbs        int32  `json:"thumbs" bson:"thumbs"`                 // 点赞数
	Date          int32  `json:"date" bson:"date"`                     // 发表时间
}

// 评论话题
type CommentTopic struct {
	Type   int32 `json:"type" bson:"type"`       // 评论主体类型 ：英雄、物品等
	TypeId int32 `json:"type_id" bson:"type_id"` // 评论主体type_id：英雄type_id、物品type_id等
}

// 评论元数据
type CommentMetadata struct {
	CommentId int64        `json:"_id" bson:"_id"`     // 评论唯一id
	Topic     CommentTopic `json:"topic" bson:"topic"` // 评论主体，作为compound_index

	PublishMetadata *PublisherMetadata   `json:"publish_metadata" bson:"publish_metadata"` // 发表者元数据
	ReplyMetadatas  []*PublisherMetadata `json:"reply_metadatas" bson:"reply_metadatas"`   // 此条评论回复列表
}

func (c *CommentMetadata) FromPB(pb *pbGlobal.CommentMetadata) {
	c.CommentId = pb.GetCommentId()

	c.Topic.Type = pb.GetTopic().GetTopicType()
	c.Topic.TypeId = pb.GetTopic().GetTopicTypeId()

	c.PublishMetadata = &PublisherMetadata{
		PublisherId:   pb.GetPublisherMetadata().GetPublisherId(),
		PublisherName: pb.GetPublisherMetadata().GetPublisherName(),
		ReplyToId:     pb.GetPublisherMetadata().GetReplyToId(),
		ReplyToName:   pb.GetPublisherMetadata().GetReplyToName(),
		Content:       pb.GetPublisherMetadata().GetContent(),
		Thumbs:        pb.GetPublisherMetadata().GetThumbs(),
	}

	c.ReplyMetadatas = make([]*PublisherMetadata, 0, len(pb.GetReplyMetadatas()))
	for n := 0; n < len(pb.GetReplyMetadatas()); n++ {
		c.ReplyMetadatas = append(c.ReplyMetadatas, &PublisherMetadata{
			PublisherId:   pb.GetReplyMetadatas()[n].GetPublisherId(),
			PublisherName: pb.GetReplyMetadatas()[n].GetPublisherName(),
			ReplyToId:     pb.GetReplyMetadatas()[n].GetReplyToId(),
			ReplyToName:   pb.GetReplyMetadatas()[n].GetReplyToName(),
			Content:       pb.GetReplyMetadatas()[n].GetContent(),
			Thumbs:        pb.GetReplyMetadatas()[n].GetThumbs(),
		})
	}
}

func (c *CommentMetadata) ToPB() *pbGlobal.CommentMetadata {
	pb := &pbGlobal.CommentMetadata{
		CommentId: c.CommentId,

		Topic: &pbGlobal.CommentTopic{
			TopicType:   c.Topic.Type,
			TopicTypeId: c.Topic.TypeId,
		},

		PublisherMetadata: &pbGlobal.PublisherMetadata{
			PublisherId:   c.PublishMetadata.PublisherId,
			PublisherName: c.PublishMetadata.PublisherName,
			ReplyToId:     c.PublishMetadata.ReplyToId,
			ReplyToName:   c.PublishMetadata.ReplyToName,
			Content:       c.PublishMetadata.Content,
			Thumbs:        c.PublishMetadata.Thumbs,
			Date:          c.PublishMetadata.Date,
		},

		ReplyMetadatas: make([]*pbGlobal.PublisherMetadata, 0, len(c.ReplyMetadatas)),
	}

	for n := 0; n < len(c.ReplyMetadatas); n++ {
		pb.ReplyMetadatas = append(pb.ReplyMetadatas, &pbGlobal.PublisherMetadata{
			PublisherId:   c.ReplyMetadatas[n].PublisherId,
			PublisherName: c.ReplyMetadatas[n].PublisherName,
			ReplyToId:     c.ReplyMetadatas[n].ReplyToId,
			ReplyToName:   c.ReplyMetadatas[n].ReplyToName,
			Content:       c.ReplyMetadatas[n].Content,
			Thumbs:        c.ReplyMetadatas[n].Thumbs,
			Date:          c.ReplyMetadatas[n].Date,
		})
	}
	return pb
}
