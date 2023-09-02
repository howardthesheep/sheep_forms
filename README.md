# sheep_forms
An experiment of designing a custom file format and parser whose purpose is to streamline the creation
of user input forms. This is not meant to be a "final solution" to creating custom input forms, rather
this should be used for generating the foundational code. All "extra" business logic, such as changing
a dropdown option values based on previous field input, should be done in your language of choice
AFTER you've generated your boilerplate forms using the parser. I'm not adding overly complex features
to make this a swiss army knife of form creation. Rather, this is the fast-food equivalent of the
form creation world. Its fast, its hot, and it's the bare minimum. Is it production ready? No. 
Should it be? No. This is meant to get you to the MVP of your needs as fast as possible with no
extra fluff or bullshit. Everything else is on you, the programmer. Do your job and quit crawling
GitHub to automate your entire career.

## TODO
- Implement attribute=value parsing
- Implement Flutter template strings
- Implement HTML template strings

## Using sheep_parser

### Windows
The sheep parser is built from source code using the command:
`./build.bat build`

To run the parser, use the command:
`./build.bat run`

To install the parser for use in the CLI, use command:
`./build.bat install`

### Linux
The sheep parser is built from source code using the command:
`./build.sh build`

To run the parser, use the command:
`./build.sh run`

To install the parser for use in the CLI, use command:
`./build.sh install`

## Form Specifications

Here is an example .sheepform, which specifies our form in a quick to type, human-readable, manner.
I specifically did not want to use a verbose syntax like XML or even JSON. The main goal is to make
specifying the rough layout & details of forms as fast as possible with a minimalist syntax.

_Note: should I allow creation of multiple forms in a file? Nice convenience, but not sure if too much 
voodoo_

Example Input Form
```
Style=Material
Output=Flutter
User Input Form

Personal Details
	<First Name, text> <Last Name, text, attribute=value>
	<Age, int>
	<Address, text>
	<Phone Number, phone>
	
Another Section
	<Are you my dad?, dropdown>
		Yes
		Maybe
	<What is today?, date>
	<Do we have your consent?, checkbox>
	<How much money do you have?, double>
```
### Form Style
Users can specify the "style" of form they would like output. This influences the final look & feel
of the generated code. Valid style values include:
- Material
- Mac
- Windows

### Form Output
Users can specify the final output type of form within the .sheepform via using `Output=` at the beginning
of the file. Valid output values include:
- Flutter
- HTML/CSS

### Form Header (Optional)
The first line of text within a file that is not immediately proceeded by another line containing `\t<`.
In our example form above, the header would be `User Input Form`

### Form Section (Optional)
The subheader shown around certain sections of a form. This is used to group "like" inputs visually.
Specified by `Some Text` proceeded by another line containing `\t<` which is the beginning of
creating a form input.
In the example form above, the section would be `Personal Details` and `Another Section`

### Form Inputs
Form inputs are specified by the `<>` characters. You can have one or many form inputs per line. Form inputs
on the same line will be displayed in the inline fashion in the final form output produced by the parser.
Form inputs on the next line will not be shown inline. Rather they will be displayed on the next line in
the final form output produced by the parser.

### Form Input Types
The array of unique form input types which influence the formatting of the final form output produced by 
the .sheepform parser. The allowed types are:
- text (default)
- int
- double
- phone
- email
- dropdown
- date
- date range
- checkbox
- tri-state box
- files
- images
- time
- rich text (emojis, bold, italics, etc.)
- slider
- captcha
- color
- credit card
- address
- search and select
- progress bar _Note: This isn't really a form input, but it is reusable. Is within scope?_

### Form Input Attributes
There is only one required form attribute and that is the `type`. The form type is specified by the text
proceeding the first `, (comma)` within a form input. Users can specify as many additional form attributes
as they would like by continuing the theme of comma-separated values. In our example above:
```
	<First Name, text> <Last Name, text, attribute=value>
	<Age, int>
```

### Form Input Options
This really only applies to the form inputs with the `dropdown` type at the moment. These form inputs
should be proceeded by the form input options that the user is allowed to choose from. By default, the
first element listed will be the default choice. Additionally, form input options should be preceded
with a `\t` to nest it under the dropdown it is associated with. In our example above:
```
<Are you my dad?, dropdown>
	Yes
	Maybe
```



