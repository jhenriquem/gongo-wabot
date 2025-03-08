package sessions

import (
	"sync"

	"github.com/jhenriquem/gongo-wabot/internal/types"
)

type Session struct {
	UserID       string
	currentOrder types.Order
	stageCounter int
}

var (
	List []*Session = []*Session{}
	mu   sync.Mutex
)

func SetSession(userID string) *Session {
	List = append(List,
		&Session{
			UserID: userID,
		})

	return List[len(List)-1]
}

func RemoveSession(userID string) {
	// Impede que varias conversas acesse a lista e causem um erro de index
	mu.Lock()
	defer mu.Unlock()

	for i, s := range List {
		if s.UserID == userID {
			List = append(List[:i], List[i+1:]...)
			break
		}
	}
}
