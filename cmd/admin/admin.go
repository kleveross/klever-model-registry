/*
 * Copyright 2017 caicloud authors. All rights reserved.
 */

package main

import (
	"fmt"

	"github.com/caicloud/golang-template-project/pkg/utils/net"

	"github.com/caicloud/go-common/version"
)

func main() {
	fmt.Println("Hello world admin")
	fmt.Printf("go-common build information: %v\n", version.Get().Pretty())

	net.Helper()
}
