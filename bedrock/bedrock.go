package bedrock

type BedrockMessage struct {
    Message string
    Model string
}

func (b *BedrockMessage) NewMessage() string {
    // TODO call bedrock with converse API
    return ""
}
