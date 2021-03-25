package shortener

//Shortener ...
type Shortener interface {
	Encode(uint64) string
	Decode(string) (uint64, error)
}
