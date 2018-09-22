package backend
import "time"

type BackendService interface {
	Save(string) (string, error)
	Load(string) (string, error)
	LoadInfo(string) (*Tuple, error)
	Close() error
}

type Tuple struct {
	URL     		string 		`json:"url"`
	Visited 		bool   		`json:"visited"`
	Count   		int    		`json:"count"`
	EntryTime		time.Time   	`json:"entrytime"`
	LastVisitedTime		time.Time   	`json:"lastvisitedtime"`
}
