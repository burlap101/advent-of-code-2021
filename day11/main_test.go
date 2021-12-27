package main

import (
	"testing"
	"fmt"
)

func TestMapAfterOneIteration (t *testing.T) {
	octopiMap, err := createOctopiMap("example.txt")
	if err != nil {
		t.Errorf("Received error %s", err)
	}
	
	cav := cavern{octopiMap, uint64(0), make(map[mcoord]struct{})}
	expectedMap, err := createOctopiMap("testdata/step1.txt")
	if err != nil {
		t.Errorf("Received error %s", err)
	}

	iterationLimit := 1
	for i := 0; i < iterationLimit; i++ {
		for m := 0; m < len(cav.octopiMap); m++ {
			for n:=0; n < len(cav.octopiMap[m]); n++ {
				cav.incrementOctopus(mcoord{m, n})
			}
		}
		cav.flashedBuffer = make(map[mcoord]struct{})
	}

	for m := 0; m < len(cav.octopiMap); m++ {
		for n := 0; n < len(cav.octopiMap); n++ {
			if cav.octopiMap[m][n] != expectedMap[m][n] {
				t.Errorf("energy at %+v=%d; expected %d", mcoord{m,n}, cav.octopiMap[m][n], expectedMap[m][n])
			}
		}
	}
}

func TestMapAfterTwoIterations (t *testing.T) {
	octopiMap, err := createOctopiMap("example.txt")
	if err != nil {
		t.Errorf("Received error %s", err)
	}
	
	cav := cavern{octopiMap, uint64(0), make(map[mcoord]struct{})}
	expectedMap, err := createOctopiMap("testdata/step2.txt")
	if err != nil {
		t.Errorf("Received error %s", err)
	}

	iterationLimit := 2
	for i := 0; i < iterationLimit; i++ {
		for m := 0; m < len(cav.octopiMap); m++ {
			for n:=0; n < len(cav.octopiMap[m]); n++ {
				cav.incrementOctopus(mcoord{m, n})
			}
		}
		cav.flashedBuffer = make(map[mcoord]struct{})
	}

	for m := 0; m < len(cav.octopiMap); m++ {
		for n := 0; n < len(cav.octopiMap); n++ {
			if cav.octopiMap[m][n] != expectedMap[m][n] {
				t.Errorf("energy at %+v=%d; expected %d", mcoord{m,n}, cav.octopiMap[m][n], expectedMap[m][n])
			}
		}
	}
	t.Logf("\n%s", cav)

}

func TestMapAfterThreeIterations (t *testing.T) {
	octopiMap, err := createOctopiMap("example.txt")
	if err != nil {
		t.Errorf("Received error %s", err)
	}
	
	cav := cavern{octopiMap, uint64(0), make(map[mcoord]struct{})}
	expectedMap, err := createOctopiMap("testdata/step3.txt")
	if err != nil {
		t.Errorf("Received error %s", err)
	}

	iterationLimit := 3
	for i := 0; i < iterationLimit; i++ {
		for m := 0; m < len(cav.octopiMap); m++ {
			for n:=0; n < len(cav.octopiMap[m]); n++ {
				cav.incrementOctopus(mcoord{m, n})
			}
		}
		cav.flashedBuffer = make(map[mcoord]struct{})
	}

	for m := 0; m < len(cav.octopiMap); m++ {
		for n := 0; n < len(cav.octopiMap); n++ {
			if cav.octopiMap[m][n] != expectedMap[m][n] {
				t.Errorf("energy at %+v=%d; expected %d", mcoord{m,n}, cav.octopiMap[m][n], expectedMap[m][n])
			}
		}
	}
	t.Logf("\n%s", cav)
}

func TestMapAfterEachIteration (t *testing.T) {
	octopiMap, err := createOctopiMap("example.txt")
	if err != nil {
		t.Errorf("Received error %s", err)
	}
	
	cav := cavern{octopiMap, uint64(0), make(map[mcoord]struct{})}

	iterationLimit := 10
	for i := 0; i < iterationLimit; i++ {
		for m := 0; m < len(cav.octopiMap); m++ {
			for n:=0; n < len(cav.octopiMap[m]); n++ {
				cav.incrementOctopus(mcoord{m, n})
			}
		}
		cav.flashedBuffer = make(map[mcoord]struct{})
		testName := fmt.Sprintf("TestMapAfter%dIterations", i+1)
		t.Run(testName, func(t *testing.T) {
			filename := fmt.Sprintf("testdata/step%d.txt", i+1)
			expectedMap, err := createOctopiMap(filename)
			if err != nil {
				t.Errorf("Received error %s", err)
			}
			for m := 0; m < len(cav.octopiMap); m++ {
				for n := 0; n < len(cav.octopiMap); n++ {
					if cav.octopiMap[m][n] != expectedMap[m][n] {
						t.Errorf("energy at %+v=%d; expected %d", mcoord{m,n}, cav.octopiMap[m][n], expectedMap[m][n])
					}
				}
			}
		})
	}
	
	t.Logf("\n%s", cav)
}