package fish

import (
	"errors"
	"fmt"
)

type LanternFish struct {
	Timer int
}

func (f *LanternFish) ServiceTimer(showTimer bool) error {
	f.Timer -= 1
	
	if f.Timer < 0 {
		f.Timer = 6
		return errors.New("timer underflow reached, plz make sure a fish was spawned")
	}
	if showTimer {
		fmt.Println(f.Timer)
	}
	return nil
}

/*
Returns a new LanternFish
*/
func (f *LanternFish) SpawnNew() (LanternFish, error) {
	newBorn := LanternFish{Timer: 8}
	return newBorn, nil
}