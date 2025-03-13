When
----
Wednesday, March 12th, 2025 @ 17:30 to 20:30

Location
--------
[Creative Fabrica](https://www.creativefabrica.com/)
Westerstraat 187
1015 MA Amsterdam

Talks
-----
- [The powerful toolset of the go-mysql library](go-mysql.pdf), by [Daniël van Eeden](https://github.com/dveeden/)
- [How we built Gopher-Verse](how-we-built-gopher-verse.pdf), by [Herlon Aguiar](https://www.linkedin.com/in/herlonaguiar/) (Creative Fabrica)

Lightning talks
--------------
- [Quiz](quiz.pdf)
- [What's new in Go](synctest.pdf), by [Raiza Claudino](https://www.linkedin.com/in/raizaclaudino/) (Creative Fabrica)
    - [Code used for the presentation](https://gist.github.com/raizacf/513b37029571be6d5f163b9a8204f01f)
- [Evolution of Map in Go](evolution-of-map-in-go.pdf), by [George Sereda](https://www.linkedin.com/in/george-sereda/)

Other slides
------------
* [Meetup intro slides](intro-slides.pdf), by [Ilija Matoski](https://www.linkedin.com/in/ilijamt/)

Video
-----

No video was recorded this meetup.

Notes
-----

There is an error in one of the questions in the quiz.

```go
func main() {
	fmt.Println("start")
	defer fmt.Println("later")
	defer panic("deferred panic")
	panic("normal panic")
}
```

The correct answer is `start, later, panic: normal panic, panic: deferred panic`
