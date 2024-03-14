package telegram

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	Id      int              `json:"update_id"`
	Message *IncomingMessage `json:"message"`
}

type IncomingMessage struct {
	Text string   `json:"text"`
	From FromUser `json:"from"`
	Chat Chat     `json:"chat"`
}

type FromUser struct {
	Username string `json:"username"`
}

type Chat struct {
	Id int `json:"id"`
}
