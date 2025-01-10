package bedrock

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	bedrocktypes "github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
)

type BedrockConverse struct {
    Message string
    Model string
    messages []bedrocktypes.Message
}

type Response struct {
    Output string `json:"output"` 
    Latency int64 `json:"latency"` 
}

func (b *BedrockConverse) NewMessage(ctx context.Context, client *bedrockruntime.Client) (Response, error) {
    // TODO call bedrock with converse API
    fmt.Println("using model:", b.Model)

    userMessage := bedrocktypes.Message{
        Content: []bedrocktypes.ContentBlock{
            &bedrocktypes.ContentBlockMemberText{
                Value: b.Message,
            },
        },
        Role: bedrocktypes.ConversationRoleUser,
    }

    messages := []bedrocktypes.Message{}
    messages = append(messages, userMessage)

    converseInput := bedrockruntime.ConverseInput{
        ModelId: &b.Model,
        Messages: messages,
    }

    bedrockOutput, err := client.Converse(ctx, &converseInput)

    if err != nil {
        return Response{}, err
    }

    responseBedrock := bedrockOutput.Output.(*bedrocktypes.ConverseOutputMemberMessage)
    latency := bedrockOutput.Metrics.LatencyMs
    bedrockMessage := responseBedrock.Value.Content[0].(*bedrocktypes.ContentBlockMemberText)

    response := Response{
        Output: bedrockMessage.Value,
        Latency: *latency,
    }


    return response, nil
}
