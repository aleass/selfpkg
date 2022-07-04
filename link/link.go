package main

import "fmt"

type guard struct {
	data *sse
}

type guard2 struct {
	live int
	data []*sse
}

type sse struct {
	data int
	pre  *sse
	next *sse
}

func (g *guard) Add(s *sse) {
	if g.data != nil {
		s.next = g.data
		g.data.pre = s
	}
	g.data = s
	return
}
func (g *guard) Del(s *sse) {
	if s == nil {
		return
	}
	if s.pre == nil { //第一位
		g.data = s.next
		if s.next != nil { //有第二位
			s.next.pre = nil
		}
	} else if s.next == nil { //末尾
		s.pre.next = nil
	} else { //中间
		s.pre.next = s.next
		s.next.pre = s.pre
	}
}

func linkes() {
	var g = guard{}
	for i := 0; i < 100; i++ {
		g.Add(&sse{
			data: i,
		})
	}
	res := g.data
	for ; res != nil; res = res.next {
		fmt.Println(res.data)
		if res.data%2 == 0 {
			g.Del(res)
		}
	}
	fmt.Println()
	res = g.data
	for ; res != nil; res = res.next {
		fmt.Println(res.data)
		g.Del(res)
	}
	fmt.Println(g.data)
	fmt.Println("------------------------------------------")

	var g2 = guard2{}
	for i := 0; i < 100; i++ {
		g2.Add(&sse{
			data: i,
		})
	}
	for id, v := range g2.data {
		fmt.Println(v.data)
		if v.data%2 == 0 {
			g2.Del(id)
		}
	}

	fmt.Println()
	for id, v := range g2.data {
		if v == nil {
			continue
		}
		fmt.Println(v.data)
		g2.Del(id)
	}
	fmt.Println(g2.data)
}

func (g *guard2) Add(s *sse) int {
	if len(g.data) == 0 {
		g.data = []*sse{s}
		return 0
	}
	for id, v := range g.data {
		if v == nil {
			g.data[id] = s
			return id
		}
	}
	g.data = append(g.data, s)
	return len(g.data) + 1
}

func (g *guard2) Del(key int) {
	g.data[key] = nil
}
