package msgraph

type Messages []Message

func (m Messages) setGraphClient(graphClient *GraphClient) {
	for _, message := range m {
		message.setGraphClient(graphClient)
	}
}
