package parser

const (
	StyleMaterial FormStyle = "Material"
	StyleMac                = "Mac"
	StyleWindows            = "Windows"
)

const (
	OutputFlutter FormOutput = "Flutter"
	OutputHtml               = "HTML"
)

type Form struct {
	Header    string
	Style     FormStyle
	Output    FormOutput
	Sections  []FormSection
	InputRows []FormInputRow
}

type FormStyle string

type FormOutput string

type FormSection struct {
	Header    string
	InputRows []FormInputRow
}

type FormInputRow struct {
	Inputs []FormInput
}

type FormInput struct {
	Title      string
	Type       FormInputType
	Attributes map[string]string
	Options    []string
}

type FormInputType int

const (
	text FormInputType = iota
	richText
	integer
	double
	phone
	email
	dropdown
	date
	dateRange
	time
	checkbox
	tristateBox
	files
	images
	slider
	captcha
	color
	creditCard
	address
	searchAndSelect
	progressBar
)
