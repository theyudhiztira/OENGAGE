package dto

type TemplateResponseDto struct {
	Data       interface{} `json:"data"`
	Pagination Cursor      `json:"pagination"`
}

type MetaTemplateResponse struct {
	Data   []Template `json:"data"`
	Paging Paging     `json:"paging"`
}

type Template struct {
	Name            string      `json:"name"`
	ParameterFormat string      `json:"parameter_format"`
	Components      []Component `json:"components"`
	Language        string      `json:"language"`
	Status          string      `json:"status"`
	Category        string      `json:"category"`
	ID              string      `json:"id"`
}

type Component struct {
	Type    string      `json:"type"`
	Format  string      `json:"format,omitempty"`
	Text    string      `json:"text,omitempty"`
	Example interface{} `json:"example,omitempty"`
	Buttons []Button    `json:"buttons,omitempty"`
}

type Button struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Paging struct {
	Cursors Cursor `json:"cursors"`
}

type Cursor struct {
	Before string `json:"before"`
	After  string `json:"after"`
}
