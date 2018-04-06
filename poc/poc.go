package poc

// POC just need to impletements Run().
type POC interface {
	Run(target string)
	Report() interface{}
}

// POCMap is a poc map.
var POCMap = map[string]POC{
	"SeaCMSv654": NewSeaCMSv654(),
}
