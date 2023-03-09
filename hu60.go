package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/rkonfj/hu60bot/pkg/hu60"
	"github.com/sirupsen/logrus"
)

func getTextMessage(content []hu60.MsgContent) (text string) {
	for _, c := range content[0].MsgUnit {
		if c.Type == "text" {
			text += *c.Value
		}
	}
	text = strings.Trim(text, " ")
	text = strings.Trim(text, "，")
	text = strings.Trim(text, "\r\n")
	text = strings.Trim(text, "\n")
	return
}

func answerHu60(client *hu60.Client, sid string, msg hu60.Msg, answer string, newConversation bool) {

	tokens := strings.Split(msg.Content[0].URL, ".")

	isTopic := tokens[0] == "bbs" && tokens[1] == "topic"

	isChatroom := tokens[0] == "addin" && tokens[1] == "chat"

	answerIntro := ""
	if newConversation {
		answerIntro = "[新会话] "
	}

	if isTopic {

		topicid, _ := strconv.Atoi(tokens[2])

		resp, err := client.GetTopic(context.Background(), topicid, sid)
		if err != nil {
			logrus.Error("answerHu60 get topic err: ", err.Error())
			return
		}

		_, err = client.ReplyTopic(context.Background(), sid, hu60.ReplyTopicRequest{
			Token:   resp.Token,
			Content: fmt.Sprintf("<!-- markdown -->\n%s@#%d，%s", answerIntro, msg.ByUID, answer),
			TopicID: topicid,
		})

		if err != nil {
			logrus.Error("answerHu60 reply err: ", err.Error())
		}
		return
	}

	if isChatroom {
		chatroomName := tokens[2]
		resp, err := client.GetChatroom(context.Background(), chatroomName, sid)
		if err != nil {
			logrus.Error("answerHu60 get chatroom err: ", err.Error())
			return
		}

		_, err = client.ReplyChatroom(context.Background(), sid, hu60.ReplyChatroomRequest{
			Token:        resp.Token,
			Content:      fmt.Sprintf("<!-- markdown -->\n%s@#%d，%s", answerIntro, msg.ByUID, answer),
			ChatroomName: chatroomName,
		})

		if err != nil {
			logrus.Error("answerHu60 reply chatroom err: ", err.Error())
		}
		return
	}

	logrus.Error("unsupported. discard: ", answer, msg)
}
