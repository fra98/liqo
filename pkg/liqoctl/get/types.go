// Copyright 2019-2025 The Liqo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package get

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/liqotech/liqo/pkg/liqoctl/factory"
	"github.com/liqotech/liqo/pkg/liqoctl/rest"
)

// NewGetCommand returns the cobra command for the get subcommand.
func NewGetCommand(ctx context.Context, liqoResources []rest.APIProvider, f *factory.Factory) *cobra.Command {
	options := &rest.GetOptions{
		Factory: f,
	}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get Liqo resources",
		Long:  "Get Liqo resources.",
		Args:  cobra.NoArgs,
	}

	f.AddNamespaceFlag(cmd.PersistentFlags())

	for _, r := range liqoResources {
		api := r()

		apiOptions := api.APIOptions()
		if apiOptions.EnableGet {
			cmd.AddCommand(api.Get(ctx, options))
		}
	}

	return cmd
}
