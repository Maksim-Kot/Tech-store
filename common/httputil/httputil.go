package httputil

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/Maksim-Kot/Commons/discovery"
)

func ServiceAddr(ctx context.Context, serviceName string, registry discovery.Registry) (string, error) {
	addrs, err := registry.ServiceAddresses(ctx, serviceName)
	if err != nil {
		if errors.Is(err, discovery.ErrNotFound) {
			return "", fmt.Errorf("no addresses found for service %q", serviceName)
		}
		return "", err
	}

	return addrs[rand.Intn(len(addrs))], nil
}
