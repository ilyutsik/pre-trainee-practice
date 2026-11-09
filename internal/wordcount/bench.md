go version go1.27.2 darwin/arm64

Оптимизация func Tokens(s string) []string  

before

goos: darwin
goarch: arm64
pkg: github.com/ilyutsik/pre-trainee-practice/internal/wordcount
cpu: Apple M2
BenchmarkTokens-8   	    3280	    365486 ns/op	  180537 B/op	    6354 allocs/op
BenchmarkTokens-8   	    3306	    363078 ns/op	  180536 B/op	    6354 allocs/op
BenchmarkTokens-8   	    3297	    363898 ns/op	  180536 B/op	    6354 allocs/op

after

goos: darwin
goarch: arm64
pkg: github.com/ilyutsik/pre-trainee-practice/internal/wordcount
cpu: Apple M2
BenchmarkTokens-8   	    6556	    180658 ns/op	  126336 B/op	    1084 allocs/op
BenchmarkTokens-8   	    6692	    180869 ns/op	  126336 B/op	    1084 allocs/op
BenchmarkTokens-8   	    6489	    181267 ns/op	  126336 B/op	    1084 allocs/op

Существенная оптимизация allocs/op и B/op выполнена за счет получение слов через (s[start:end]),
а не накопление промежуточного слова через слайс []rune (tempSlice = append(tempSlice, unicode.ToLower(v))) и его 
последующая аллокация string(tempSlice), раньше на каждое слово выделялась новая строка,
теперь аллокации остались только для слов с заглавными буквами и для самого выходного слайса,
если его изначальной емкости не хватило. Строки теперь ссылаются на исходную строку и удерживают ее в памяти.

goos: darwin
goarch: arm64
pkg: github.com/ilyutsik/pre-trainee-practice/internal/wordcount
cpu: Apple M2
        │   old.txt   │               new.txt               │
        │   sec/op    │   sec/op     vs base                │
Tokens-8   363.2µ ± 0%   182.6µ ± 1%  -49.74% (p=0.000 n=10)

         │   old.txt    │               new.txt                │
         │     B/op     │     B/op      vs base                │
Tokens-8   176.3Ki ± 0%   123.4Ki ± 0%  -30.02% (p=0.000 n=10)

         │   old.txt   │               new.txt               │
         │  allocs/op  │  allocs/op   vs base                │
Tokens-8   6.354k ± 0%   1.084k ± 0%  -82.94% (p=0.000 n=10)


wordcount bench test

goos: darwin
goarch: arm64
pkg: github.com/ilyutsik/pre-trainee-practice/internal/wordcount
cpu: Apple M2
BenchmarkTokens-8                   6538            181050 ns/op          126336 B/op       1084 allocs/op
BenchmarkCount-8                    9498            128171 ns/op          218480 B/op         18 allocs/op
BenchmarkJoinPlus-8                   93          12378931 ns/op        119902087 B/op      6360 allocs/op
BenchmarkJoinBuilder-8             31088             37872 ns/op          154360 B/op         20 allocs/op

joinPlus на каждой итерации выделяет новую строку и копирует в неё всё накопленное, так как строки в Go неизменяемы.
Для n кусков это O(n²) копирования: ~120 МБ и ~12 мс на вызов. strings.Builder дописывает в один растущий буфер (амортизированно O(n))
и отдаёт его в String() без копирования: ~154 КБ и ~38 мкс, примерно в 300 раз быстрее и в 800 раз меньше памяти.