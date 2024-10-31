package botgo

import (
	"context"
	"fmt"
	openai "github.com/sashabaranov/go-openai"
	"github.com/tencent-connect/botgo/openapi"
	"testing"
)

func TestName(t *testing.T) {
	client := openai.NewClient("sk-0dOOjZGkBkmAw1nv0eE6171e31D34eA5Bd58D9744cAbF8Cc")
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Hello!",
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
}

func TestUseOpenAPIVersion(t *testing.T) {
	type args struct {
		version openapi.APIVersion
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"not found", args{version: 0}, true,
		},
		{
			"v1 found", args{version: 1}, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SelectOpenAPIVersion(tt.args.version); (err != nil) != tt.wantErr {
				t.Errorf("SelectOpenAPIVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
