// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hash

import "unsafe"

var BigEndian = (*[2]uint8)(unsafe.Pointer(&[]uint16{1}))[0] == 0
