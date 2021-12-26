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

// 回复元数据
type ReplyerMetadata struct {
	CommentId         int64 `json:"_id" bson:"_id"` // 评论唯一id
	PublisherMetadata `json:",inline" bson:"inline"`
}

// 评论话题
type CommentTopic struct {
	Type   int32 `json:"type" bson:"type"`       // 评论主体类型 ：英雄、物品等
	TypeId int32 `json:"type_id" bson:"type_id"` // 评论主体type_id：英雄type_id、物品type_id等
}

func (c *CommentTopic) Valid() bool {
	return c.Type >= TopicType_Begin && c.Type < TopicType_End && c.TypeId != -1
}

func (c *CommentTopic) Equal(cc *CommentTopic) bool {
	return c.Type == cc.Type && c.TypeId == cc.TypeId
}

func (c *CommentTopic) FromPB(pb *pbGlobal.CommentTopic) {
	c.Type = pb.GetTopicType()
	c.TypeId = pb.GetTopicTypeId()
}

func (c *CommentTopic) ToPB() *pbGlobal.CommentTopic {
	pb := &pbGlobal.CommentTopic{
		TopicType:   c.Type,
		TopicTypeId: c.TypeId,
	}
	return pb
}

// 评论元数据
type CommentMetadata struct {
	CommentId int64        `json:"_id" bson:"_id"`     // 评论唯一id
	Topic     CommentTopic `json:"topic" bson:"topic"` // 评论主体，作为compound_index

	PublisherMetadata *PublisherMetadata `json:"publisher_metadata" bson:"publisher_metadata"` // 发表者元数据
	ReplyerMetadatas  []*ReplyerMetadata `json:"replyer_metadatas" bson:"replyer_metadatas"`   // 此条评论回复列表
}

func (c *CommentMetadata) FromPB(pb *pbGlobal.CommentMetadata) {
	c.CommentId = pb.GetCommentId()

	c.Topic.Type = pb.GetTopic().GetTopicType()
	c.Topic.TypeId = pb.GetTopic().GetTopicTypeId()

	c.PublisherMetadata = &PublisherMetadata{
		PublisherId:   pb.GetPublisherMetadata().GetPublisherId(),
		PublisherName: pb.GetPublisherMetadata().GetPublisherName(),
		ReplyToId:     pb.GetPublisherMetadata().GetReplyToId(),
		ReplyToName:   pb.GetPublisherMetadata().GetReplyToName(),
		Content:       pb.GetPublisherMetadata().GetContent(),
		Thumbs:        pb.GetPublisherMetadata().GetThumbs(),
	}

	c.ReplyerMetadatas = make([]*ReplyerMetadata, 0, len(pb.GetReplyMetadatas()))
	for n := 0; n < len(pb.GetReplyMetadatas()); n++ {
		c.ReplyerMetadatas = append(c.ReplyerMetadatas, &ReplyerMetadata{
			CommentId: c.CommentId,
			PublisherMetadata: PublisherMetadata{
				PublisherId:   pb.GetReplyMetadatas()[n].GetPublisherMetadata().GetPublisherId(),
				PublisherName: pb.GetReplyMetadatas()[n].GetPublisherMetadata().GetPublisherName(),
				ReplyToId:     pb.GetReplyMetadatas()[n].GetPublisherMetadata().GetReplyToId(),
				ReplyToName:   pb.GetReplyMetadatas()[n].GetPublisherMetadata().GetReplyToName(),
				Content:       pb.GetReplyMetadatas()[n].GetPublisherMetadata().GetContent(),
				Thumbs:        pb.GetReplyMetadatas()[n].GetPublisherMetadata().GetThumbs(),
			},
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
			PublisherId:   c.PublisherMetadata.PublisherId,
			PublisherName: c.PublisherMetadata.PublisherName,
			ReplyToId:     c.PublisherMetadata.ReplyToId,
			ReplyToName:   c.PublisherMetadata.ReplyToName,
			Content:       c.PublisherMetadata.Content,
			Thumbs:        c.PublisherMetadata.Thumbs,
			Date:          c.PublisherMetadata.Date,
		},

		ReplyMetadatas: make([]*pbGlobal.ReplyerMetadata, 0, len(c.ReplyerMetadatas)),
	}

	for n := 0; n < len(c.ReplyerMetadatas); n++ {
		pb.ReplyMetadatas = append(pb.ReplyMetadatas, &pbGlobal.ReplyerMetadata{
			CommentId: c.CommentId,
			PublisherMetadata: &pbGlobal.PublisherMetadata{
				PublisherId:   c.ReplyerMetadatas[n].PublisherId,
				PublisherName: c.ReplyerMetadatas[n].PublisherName,
				ReplyToId:     c.ReplyerMetadatas[n].ReplyToId,
				ReplyToName:   c.ReplyerMetadatas[n].ReplyToName,
				Content:       c.ReplyerMetadatas[n].Content,
				Thumbs:        c.ReplyerMetadatas[n].Thumbs,
				Date:          c.ReplyerMetadatas[n].Date,
			},
		})
	}
	return pb
}
