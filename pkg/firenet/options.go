package firenet

import "fire/pkg/fireface"

type Option func(s *Server)

// 只要实现Packet 接口可自由实现数据包解析格式，如果没有则使用默认解析格式
func WithPacket(pack fireface.Packet) Option {
	return func(s *Server) {
		s.packet = pack
	}
}
