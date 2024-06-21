package types

type ChannelType string

type Msg struct {
	Text string
}

type RichMsgInterface interface {
	BuildContent(ChannelType) interface{}
}

type RichMsg struct {
	Type    string
	Content interface{}
}

type ChannelRichMsgMap map[ChannelType]RichMsg

type AlertLevel int

const (
	AlertLevelFatal   AlertLevel = 1
	AlertLevelError   AlertLevel = 2
	AlertLevelWarn    AlertLevel = 3
	AlertLevelInfo    AlertLevel = 4
	AlertLevelDebug   AlertLevel = 5
	AlertLevelVerbose AlertLevel = 6
)
