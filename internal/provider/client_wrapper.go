// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"strings"

	"github.com/cenkalti/backoff/v5"
)

func ClientBackoff[T any](f func() (*T, error)) func() (*T, error) {

	return func() (*T, error) {
		response, err := f()

		if err != nil && strings.Contains(err.Error(), "EOF") {
			return nil, backoff.RetryAfter(1)
		}

		if err != nil {
			return nil, backoff.Permanent(err)
		}

		return response, nil
	}
}
