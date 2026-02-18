package logic

func LoadContextForMessageType(messageType MessageType) ([]Message, error) {
	messages, err := LoadContext()
	if err != nil {
		return nil, err
	}

	var allowedTypes []MessageType

	switch messageType {
	case MessageTypeCommand:
		allowedTypes = []MessageType{MessageTypeFile, MessageTypeCommand, MessageTypeAction}
	case MessageTypeObjective:
		allowedTypes = []MessageType{MessageTypeFile, MessageTypeQuestion, MessageTypeAnswer, MessageTypeObjective, MessageTypePlan, MessageTypeRevision}
	case MessageTypeQuestion:
		allowedTypes = []MessageType{MessageTypeFile, MessageTypeQuestion, MessageTypeAnswer, MessageTypeObjective, MessageTypePlan}
	default:
		return make([]Message, 0), nil
	}

	var filtered []Message
	for _, msg := range messages {
		for _, allowed := range allowedTypes {
			if msg.Type == allowed {
				if messageType == MessageTypeObjective && msg.Type == MessageTypePlan {
					filtered = append(filtered, Message{Type: MessageTypeRevision, Content: msg.Content})
				} else {
					filtered = append(filtered, msg)
				}
				break
			}
		}
	}

	return filtered, nil
}

type MessageType string

const (
	MessageTypeFile      MessageType = "File"
	MessageTypeQuestion  MessageType = "Question"
	MessageTypeAnswer    MessageType = "Answer"
	MessageTypeCommand   MessageType = "Command"
	MessageTypeAction    MessageType = "Action"
	MessageTypeObjective MessageType = "Objective"
	MessageTypePlan      MessageType = "Plan"
	MessageTypeRevision  MessageType = "Revision"
)

func ResponseType(messageType MessageType) MessageType {

	switch messageType {
	case MessageTypeCommand:
		return MessageTypeAction
	case MessageTypeQuestion:
		return MessageTypeAnswer
	case MessageTypeObjective:
		return MessageTypePlan
	default:
		return messageType
	}
}

type Message struct {
	Type    MessageType
	Path    string
	Content string
}
