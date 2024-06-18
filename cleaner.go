package main

import "fmt"

func (s *Server) cleanKV() {

	// clean kv
	fmt.Println("cleaning kv")

	fmt.Println(s.kv.List())
	all := s.kv.List()

	for k, v  := range all {
		fmt.Println(k, string(v))

		s.kv.Del([]byte(k))
	}
}
