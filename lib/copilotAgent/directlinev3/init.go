package directlinev3

import (
	"context"
)

const (
	HTTP_POST string = "POST"
	HTTP_GET         = "GET"

	secret             = "rPa8k8RaboA.OmS6ffObQmYGS81qYOOH32Ovp38jkQU0u0pOuZaNsU4"
	ApiPath            = "https://directline.botframework.com/v3/directline/"
	DefaultLocale      = "en-US"
	DefaultMessageType = "message"
)

type (
	IDirectLineAPI interface {
		GenerateToken(context.Context) (GenerateTokenResp, error)
		StartConversation(context.Context) (CreatedConversation, error)
		SendMessage(context.Context, string, string, string) (SendMessageResp, error)
		ReceiveMessages(context.Context, string, int) (ReceiveMessages, error)
	}

	directLine struct {
		Token string
	}
)

var (
	_ IDirectLineAPI = &directLine{}
)

func init() {
	//TODO:
}

func New() *directLine {
	return &directLine{}
}

type (
	GenerateTokenResp struct {
		ConversationId string `json:"conversationId"`
		Token          string `json:"token"`
		ExpiresIn      uint32 `json:"expires_in"`
	}

	CreatedConversation struct {
		ConversationId     string `json:"conversationId"`
		Token              string `json:"token"`
		ExpiresIn          uint32 `json:"expires_in"`
		StreamUrl          string `json:"streamUrl"`
		ReferenceGrammarId string `json:"referenceGrammarId"`
	}

	From struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Role string `json:"role"`
	}

	Conversation struct {
		ID string `json:"id"`
	}

	SendMessageReq struct {
		Locale string `json:"locale"`
		Type   string `json:"type"`
		From   From   `json:"from"`
		Text   string `json:"text"`
	}

	SendMessageResp struct {
		ID string `json:"id"`
	}

	Activity struct {
		Type         string       `json:"type"`
		ID           string       `json:"id"`
		Timestamp    string       `json:"timestamp"`
		From         From         `json:"from"`
		Conversation Conversation `json:"conversation"`
		Text         string       `json:"text"`
		Speak        string       `json:"speak"`
		ReplyToId    string       `json:"replyToId"`
	}

	ReceiveMessages struct {
		Activities []Activity `json:"activities"`
	}

	WritingTaskResp struct {
		Details      []Criteria `json:"details"`
		OverallScore string     `json:"overall_score"`
		SuggestEssay string     `json:"suggest_essay"`
	}

	Criteria struct {
		Name         string `json:"name"`
		BandScore    string `json:"band_score"`
		HowToImprove string `json:"how_to_improve"`
		Strengths    string `json:"strengths"`
	}
)
