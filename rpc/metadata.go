package rpc

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)


type MetaData struct {
	UserAgent string
	ClientIP  string
}

const (
	userAgentHeader            = "user-agent"
	xForwardedForHeader        = "x-forwarded-for"
)

func metaFromCtx(ctx context.Context) *MetaData {
	meta := new(MetaData)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgents := md.Get(userAgentHeader); len(userAgents) > 0 {
			meta.UserAgent = userAgents[0]
		}
		if clientIPs := md.Get(xForwardedForHeader); len(clientIPs) > 0 {
			meta.ClientIP = clientIPs[0]
		}
	}
	if p, ok := peer.FromContext(ctx); ok {
		meta.ClientIP = p.Addr.String()
	}
	return meta
}
