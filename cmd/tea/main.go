package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type spell int

// TODO use Invoker responses
const (
	Undefined spell = iota
	ColdSnap
	GhostWalk
	IceWall
	Tornado
	EMP
	Alacrity
	SunStrike
	ForgeSpirit
	ChaosMeteor
	DeafeningBlast
)

var spellMap = map[spell]string{
	Undefined:      "",
	ColdSnap:       "Cold Snap",
	GhostWalk:      "Ghost Walk",
	IceWall:        "Ice Wall",
	Tornado:        "Tornado",
	EMP:            "EMP",
	Alacrity:       "Alacrity",
	SunStrike:      "Sun Strike",
	ForgeSpirit:    "Forge Spirit",
	ChaosMeteor:    "Chaos Meteor",
	DeafeningBlast: "Deafening Blast",
}

type incantate map[spell]struct{}

// model??? really??? in 2023???
// think of a better name
// admittedly this is a copy paste from tutorial
type model struct {
	orbs      []string
	invoked   []spell
	cast      bool
	incantate incantate
	point     int
}

func main() {
	m := model{
		orbs:      make([]string, 3, 3),
		invoked:   make([]spell, 2, 2),
		incantate: make(map[spell]struct{}),
	}
	m.cast, m.incantate = gen(nil)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func (m model) Init() tea.Cmd {
	// TODO init spells to incantate
	// can be 2 or 3 spells
	// can require casting or not
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q", "w", "e": // handle upper case?
			m.orbs[2] = m.orbs[1]
			m.orbs[1] = m.orbs[0]
			m.orbs[0] = strings.ToUpper(msg.String())
		case "r":
			i := invoke(m.orbs)
			if i == Undefined {
				break
			}
			// if stage only requires invoking, we can still count?
			if !m.cast {
				if _, ok := m.incantate[i]; ok {
					delete(m.incantate, i)
					m.point++
				}
			}
			if m.invoked[0] == i {
				break
			}
			m.invoked[1] = m.invoked[0]
			m.invoked[0] = i
		case "d":
			if _, ok := m.incantate[m.invoked[0]]; ok {
				delete(m.incantate, m.invoked[0])
				m.point++
			}
		case "f":
			if _, ok := m.incantate[m.invoked[1]]; ok {
				delete(m.incantate, m.invoked[1])
				m.point++
			}
		}
	}
	if len(m.incantate) == 0 {
		if !m.cast {
			m.cast, m.incantate = gen(m.incantate)
		} else {
			m.cast, m.incantate = gen(nil)
		}
	}
	return m, nil
}

// TODO allocate slice
// use strings builder
// nice box
// !!! center div Kappa
func (m model) View() string {
	var s string

	if !m.cast {
		s += "INVOKE\n\n"
	} else {
		s += "INVOKE AND CAST\n\n"
	}

	var incantate []string
	for k := range m.incantate {
		incantate = append(incantate, spellMap[k])
	}
	s += "| "
	s += strings.Join(incantate, " | ")
	s += " |"
	s += "\n--------------------------------------------------\n"

	s += "| "
	s += strings.Join(m.orbs, " | ")
	s += " |"
	s += "\n--------------------------------------------------\n"

	s += "| Q | W | E |"
	for _, i := range m.invoked {
		s += " "
		s += spellMap[i]
		s += " |"
	}
	s += "\n--------------------------------------------------\n"

	s += fmt.Sprintf("%d POINTS\n\n", m.point)
	return s
}

func invoke(orbs []string) spell {
	var count int
	q, w, e := "Q", "W", "E" // cosmetic purpose
	combo := make(map[string]int)
	for _, o := range orbs {
		if o == "" {
			continue
		}
		combo[o]++
		count++
	}
	if count < 3 {
		return Undefined
	}
	if combo[q] == 3 {
		return ColdSnap
	} else if combo[w] == 3 {
		return EMP
	} else if combo[e] == 3 {
		return SunStrike
	} else if combo[q] == 2 {
		if combo[w] == 1 {
			return GhostWalk
		} else if combo[e] == 1 {
			return IceWall
		}
	} else if combo[w] == 2 {
		if combo[q] == 1 {
			return Tornado
		} else if combo[e] == 1 {
			return Alacrity
		}
	} else if combo[e] == 2 {
		if combo[q] == 1 {
			return ForgeSpirit
		} else if combo[w] == 1 {
			return ChaosMeteor
		}
	}
	return DeafeningBlast
}

// TODO performance?
// create random array of spells to incantate
// it can either be 2 or 3 spells
// if previous incantate is nil, create a new one with no overlap
// else overlap is fine
func gen(prev incantate) (bool, incantate) {
	var cast bool
	r := rand.New(rand.NewSource(time.Now().Unix()))
	if c := r.Intn(2); c != 0 {
		cast = true
	}

	next := make(map[spell]struct{})
	length := 2 + r.Intn(2)
	// TODO check how many passes do we need to generate
	for len(next) < length {
		n := spell(1 + r.Intn(10)) // this will rarely gets to deafening?
		if prev != nil {
			if _, ok := prev[n]; ok {
				continue
			}
		}
		if _, ok := next[n]; ok {
			continue
		}
		next[n] = struct{}{}
	}
	return cast, next
}