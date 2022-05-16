package shellPlugin

import (
	"sync"

	"github.com/AnimeKaizoku/ssg/ssg"
)

var (
	lastId           int
	idGeneratorMutex *sync.Mutex
)

var commandsMap = ssg.NewSafeMap[string, commandContainer]()
