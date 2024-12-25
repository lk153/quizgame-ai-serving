package copilotAgent

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	lineApiLib "github.com/lk153/quizgame-ai-serving/lib/copilotAgent/directlinev3"
)

type InputTask struct {
	TaskType        uint8
	TaskRequirement string
	TaskRelatedDoc  string
	CandidateText   string
}

func createInputPromptDemo(input InputTask) string {
	if input.TaskType == 1 {
		return fmt.Sprintf(`
- You are tasked with evaluating and scoring IELTS writing task %d. Your goal is to provide detailed feedback and improvement advice to the student. - Content of task is : ((( %s ))) - Candidate response is: ((( %s ))) - I do not include chart. Please provide your evaluation based on given task. Following below structure for returning: Details: 1) Task Achievement: - Band score: - How to improve: - Strengths: 2) Coherence and Cohesion: - Band score: - How to improve: - Strengths: 3) Lexical Resource: - Band score: - How to improve: - Strengths: 4) Grammatical Range and Accuracy: - Band score: - How to improve: - Strengths: Overall Score: Suggest Essay:`, input.TaskType, input.TaskRequirement, input.CandidateText)
	}

	return fmt.Sprintf(`
- You are tasked with evaluating and scoring IELTS writing task %d. Your goal is to provide detailed feedback and improvement advice to the student. - Content of task is : ((( %s ))) - Candidate response is: ((( %s ))) - I do not include chart. Please provide your evaluation based on given task. Following below structure for returning: Details: 1) Task Response: - Band score: - How to improve: - Strengths: 2) Coherence and Cohesion: - Band score: - How to improve: - Strengths: 3) Lexical Resource: - Band score: - How to improve: - Strengths: 4) Grammatical Range and Accuracy: - Band score: - How to improve: - Strengths: Overall Score: Suggest Essay:`, input.TaskType, input.TaskRequirement, input.CandidateText)

}

func createInputPrompt(input InputTask) string {
	if input.TaskType == 1 {
		return fmt.Sprintf(`
You are IELTS teacher for assessing Writing Task %d. Please help provide assessment IELTS band score.
Given task is : %s
Candidate response is: %s
I do not include chart. Please provide your evaluation based on given task. Following below json structure for returning:
{
  "details": [
    {
      "name": "task-achievement",
      "band_score": "",
      "how_to_improve": "",
      "strengths": ""
    },
    {
      "name": "coherence-and-cohesion",
      "band_score": "",
      "how_to_improve": "",
      "strengths": ""
    },
    {
      "name": "lexical-resource",
      "band_score": "",
      "how_to_improve": "",
      "strengths": ""
    },
    {
      "name": "grammatical-range-and-accuracy",
      "band_score": "",
      "how_to_improve": "",
      "strengths": ""
    }
  ],
  "overall_score": "",
  "suggest_essay": ""
}`, input.TaskType, input.TaskRequirement, input.CandidateText)
	}

	return fmt.Sprintf(`
You are IELTS teacher for assessing Writing Task %d. Please help provide assessment IELTS band score.
Given task is : %s
Candidate response is: %s
Please provide your evaluation based on given task. Following below json structure for returning:
{
  "details": [
    {
      "name": "task-response",
      "band_score": "",
      "how_to_improve": "",
      "strengths": ""
    },
    {
      "name": "coherence-and-cohesion",
      "band_score": "",
      "how_to_improve": "",
      "strengths": ""
    },
    {
      "name": "lexical-resource",
      "band_score": "",
      "how_to_improve": "",
      "strengths": ""
    },
    {
      "name": "grammatical-range-and-accuracy",
      "band_score": "",
      "how_to_improve": "",
      "strengths": ""
    }
  ],
  "overall_score": "",
  "suggest_essay": ""
}`, input.TaskType, input.TaskRequirement, input.CandidateText)

}

func DoAssessmentV1(ctx context.Context, userID string, input InputTask) (result string, err error) {
	api := lineApiLib.New()
	api.Token = os.Getenv("COPILOT_TOKEN")
	conversationId := os.Getenv("COPILOT_CONVERSATION_ID")
	//Send messages
	prompt := createInputPromptDemo(input)
	fmt.Println("PROMPT:", prompt)
	resp, err1 := api.SendMessage(ctx, conversationId, userID, prompt)
	if err1 != nil {
		err = err1
		return
	}

	if strings.TrimSpace(resp.ID) == "" {
		err = errors.New("SendMessage:ERR")
		return
	}

	sentID := resp.ID
	watermark := getWatermark(sentID)

	//Receive messages
ReceiveMessage:
	sleepTime := time.Duration(1)
	for {
		time.Sleep(sleepTime * time.Second)
		resp1, err1 := api.ReceiveMessages(ctx, conversationId, watermark)
		if err1 != nil {
			err = err1
			return
		}

		if len(resp1.Activities) == 0 {
			sleepTime *= 2
			continue
		}

		msgStr := resp1.Activities[len(resp1.Activities)-1]

		//Tricky send message to force it return assessment instead noise us to rephrase again our prompt
		if strings.Contains(msgStr.Text, "I'm sorry, I'm not sure how to help with that. Can you try rephrasing?") {
			resp, err1 := api.SendMessage(ctx, conversationId, userID, "What 's problem?")
			if err1 != nil {
				err = err1
				return
			}
			sentID := resp.ID
			watermark = getWatermark(sentID)
			goto ReceiveMessage
		}

		msgStr.Text = strings.TrimLeft(msgStr.Text, "There was no problem with your previous request. Here is the evaluation based on the given task and candidate response:")
		return strings.TrimRight(msgStr.Text, msgStr.Speak), nil
	}

}

func DoAssessment(ctx context.Context, userID string, input InputTask) (result string, err error) {
	api := lineApiLib.New()
	//Generate token
	token, err := api.GenerateToken(ctx)
	if err != nil {
		return
	}

	if strings.TrimSpace(token.ConversationId) == "" {
		err = errors.New("GenerateToken:ERR")
		return
	}

	//Start conversation
	conversation, err := api.StartConversation(ctx)
	if err != nil {
		return
	}

	if strings.TrimSpace(conversation.Token) == "" {
		err = errors.New("StartConversation:ERR")
		return
	}

	//Send messages
	prompt := createInputPrompt(input)
	fmt.Println("PROMPT:", prompt)
	resp, err1 := api.SendMessage(ctx, conversation.ConversationId, userID, prompt)
	if err1 != nil {
		err = err1
		return
	}

	if strings.TrimSpace(resp.ID) == "" {
		err = errors.New("SendMessage:ERR")
		return
	}

	sentID := resp.ID
	watermark := getWatermark(sentID)

	//Receive messages
	sleepTime := time.Duration(1)
	for {
		time.Sleep(sleepTime * time.Second)
		resp1, err1 := api.ReceiveMessages(ctx, conversation.ConversationId, watermark)
		if err1 != nil {
			err = err1
			return
		}

		if len(resp1.Activities) == 0 {
			sleepTime *= 2
			continue
		}

		msg := resp1.Activities[len(resp1.Activities)-1]
		if msg.ReplyToId == sentID {
			return strings.TrimLeft(msg.Text, msg.Speak), nil
		}
		sleepTime *= 2
	}
}

func getWatermark(sentID string) int {
	sentIDs := strings.Split(sentID, "|")
	watermark, err := strconv.Atoi(sentIDs[1])
	if err != nil {
		log.Println("SendMessage:ConvertWatermarkERR:", err)
		return 0
	}

	return watermark
}
