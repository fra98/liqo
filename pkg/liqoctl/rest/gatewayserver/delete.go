// Copyright 2019-2023 The Liqo Authors
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

package gatewayserver

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	networkingv1alpha1 "github.com/liqotech/liqo/apis/networking/v1alpha1"
	"github.com/liqotech/liqo/pkg/liqoctl/completion"
	"github.com/liqotech/liqo/pkg/liqoctl/output"
	"github.com/liqotech/liqo/pkg/liqoctl/rest"
)

const liqoctlDeleteGatewayServerLongHelp = `Delete a GatewayServer.

Examples:
  $ {{ .Executable }} delete gatewayserver my-gateway-server`

// Delete deletes a GatewayServer.
func (o *Options) Delete(ctx context.Context, options *rest.DeleteOptions) *cobra.Command {
	o.deleteOptions = options

	cmd := &cobra.Command{
		Use:     "gatewayserver",
		Aliases: []string{"gatewayservers", "server", "servers", "gws"},
		Short:   "Delete a GatewayServer",
		Long:    liqoctlDeleteGatewayServerLongHelp,

		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completion.GatewayServers(ctx, o.deleteOptions.Factory, 1),

		PreRun: func(cmd *cobra.Command, args []string) {
			options.Name = args[0]
			o.deleteOptions = options
		},

		Run: func(cmd *cobra.Command, args []string) {
			output.ExitOnErr(o.handleDelete(ctx))
		},
	}

	return cmd
}

func (o *Options) handleDelete(ctx context.Context) error {
	opts := o.deleteOptions
	s := opts.Printer.StartSpinner("Deleting GatewayServer")

	gatewayServer := &networkingv1alpha1.GatewayServer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      opts.Name,
			Namespace: opts.Namespace,
		},
	}
	if err := o.deleteOptions.CRClient.Delete(ctx, gatewayServer); err != nil {
		err = fmt.Errorf("unable to delete GatewayServer: %w", err)
		s.Fail(err)
		return err
	}

	s.Success("GatewayServer deleted")
	return nil
}