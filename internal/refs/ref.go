package refs

import (
	"fmt"
	"strings"
	"sync"

	"github.com/DQGriffin/labrador/internal/cli/console"
)

var refs = make(map[string]string)
var mu = sync.RWMutex{}

func IsRef(s string) bool {
	return strings.Contains(s, "[")
}

func SetRef(ref, id string) {
	mu.Lock()
	defer mu.Unlock()

	console.Debugf("Setting ref %s to %s", ref, id)
	refs[ref] = id
}

func RemoveRef(ref string) {
	mu.Lock()
	defer mu.Unlock()

	console.Debugf("Removing ref %s", ref)
	delete(refs, ref)
}

func ResolveRef(ref string) (string, error) {
	mu.RLock()
	defer mu.RUnlock()

	console.Debugf("Getting ref %s", ref)

	id, ok := refs[ref]
	if ok {
		console.Debugf("Ref %s resolved to %s", ref, id)
		return id, nil
	}

	adjustedRef := strings.ReplaceAll(ref, "[[", "")
	adjustedRef = strings.ReplaceAll(adjustedRef, "]]", "")
	id2, ok2 := refs[adjustedRef]

	if ok2 {
		console.Debugf("Ref %s resolved to %s", ref, id2)
		return id2, nil
	}

	console.Debugf("Ref %s could not be resolved", ref)
	return "", fmt.Errorf("undefined ref: %s", ref)
}
