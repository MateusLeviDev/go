#### how do we run the code in our project?

```
go run main.go
```

- go run: compiles and executes one or two files
- go build: compiles a bunch of go source code files
- go fmt: formats all the code in each file in the current directory
- go install: compiles and installs a package
- go get: downloads the raw source code of someone else package
- go test: test the current project

#### what does 'package main' mean?

voce pode pensar um pacote sendo como um projeto ou um workspace.

- the package is a collection of common source code files
- um pacote pode ter muitos arquivos relacionados dentro dele. cada arquivo que terminar com .go
- o unico requisito esta na primeira linha de cada file deve declarar o pacote no qual pertence

- executable package: package main -> defines a package that can be compiled and then executed must have a func called main
- reusable package: package calculator, package uploader: defines a package that can be used as a dependency (helper code)

#### what does 'import "fmt"' means?

fmt: standard lib package that is included with the go by default. (format)

- used to print out info, specifically to the terminal, just to give a better sense of debugging

- usamos o import to form um link from our package to these other ones

<br>

`golang.org/pkg`

<br>

#### basic go types

- bool, string, int, float64
- many others...

<br>

card := "ace of spades"

- we only use this colon-equals syntax when we are defining a new variable
- if we are reassigning a existing variable a new value, we do not have to use this anymore

- array and slice: array fixed length of thing. slice grow or shrink. deve ser o mesmo tipo

```
package main

import "fmt"

func main() {
	cards := []string{"Ace of Diamonds", newCard()}
	cards = append(cards, "Six of Spades")

	for i, card := range cards {
		fmt.Println(i, card)
	}
}

func newCard() string {
	return "Five of Diamonds"
}

```

- in this case. dont we really just have to declare i and card one time? ppr que usar novamente o :=?

with for loops, every single time that we step through this list of cards, we are really throwing away the previous index and card that had been declared. por isso estamos re-declaring the variables i and card here by using :=