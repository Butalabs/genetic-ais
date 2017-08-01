package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/MaxHalford/gago"
)

//0 rock
//1 paper
//2 scissors
//3 well
type Player []float64

func (p Player) String() string {
	return fmt.Sprintf("rock: %.3f, paper: %.3f, scissors: %.3f, well: %.3f", p[0], p[1], p[2], p[3])
}

func (p Player) getChoice() int {
	a := rand.Float64()
	b := 0.0
	for x, y := range p {
		b += y
		if a < b {
			return x
		}
	}
	return -1
}

func getRandomChoice() int {
	p := Player([]float64{0.1651, 0.3213, 0.2303, 0.2833})
	return p.getChoice()
}

func calcSet(a, b int) bool {
	return a == 0 && b == 2 || a == 1 && (b == 0 || b == 3) || a == 2 && b == 1 || a == 3 && (b == 0 || b == 2)
}

func (p Player) Evaluate() float64 {
	var wins float64
	for i := 0; i < 1000; i++ {
		if calcSet(p.getChoice(), getRandomChoice()) {
			wins--
		}
	}
	return wins
}

func (p Player) Mutate(rng *rand.Rand) {
	maxWidth := 1.0
	for i, _ := range p {
		maxWidth = math.Min(maxWidth, (p[i])/2)
	}

	possibilities := []([]float64){
		[]float64{1, 1, -1, -1},
		[]float64{1, -1, 1, -1},
		[]float64{1, -1, -1, 1},
		[]float64{-1, 1, 1, -1},
		[]float64{-1, 1, -1, 1},
		[]float64{-1, -1, 1, 1},
	}
	width := rng.Float64() * maxWidth
	chosen := possibilities[rng.Intn(6)]

	for i, v := range p {
		p[i] = v + chosen[i]*width
		if p[i] < 0 {
			fmt.Println(v, chosen[i], width)
		}
	}
}

func (p Player) Crossover(p2 gago.Genome, rng *rand.Rand) (gago.Genome, gago.Genome) {
	var o1, o2 = gago.CrossUniformFloat64(p, p2.(Player), rng) // Returns two float64 slices
	return Player(o1), Player(o2)
}

func (p Player) Clone() gago.Genome {
	var p2 = make(Player, len(p))
	copy(p2, p)
	return p2
}

func MakePlayer(rng *rand.Rand) gago.Genome {
	floats := make([]float64, 4)
	floats[0] = rng.Float64()
	remaining := 1 - floats[0]
	for i := 1; i < 3; i++ {
		floats[i] = rng.Float64() * remaining
		remaining -= floats[i]
	}
	floats[3] = remaining
	return Player(floats)
}

func main() {
	var ga = gago.Generational(MakePlayer)
	ga.Initialize()

	fmt.Printf("Best fitness at generation 0: %f, with player: %s\n", ga.Best.Fitness, ga.Best.Genome)
	for i := 0; i < 50; i++ {
		ga.Enhance()
		fmt.Printf("Best fitness at generation %d: %f, with player: %s\n", i, ga.Best.Fitness, ga.Best.Genome)
	}
}
