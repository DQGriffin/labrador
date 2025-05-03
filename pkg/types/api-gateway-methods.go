package types

import (
	"fmt"

	"github.com/DQGriffin/labrador/internal/cli/styles"
	"github.com/DQGriffin/labrador/pkg/helpers"
	"github.com/charmbracelet/lipgloss/tree"
)

func (config ApiGatewayConfig) ToTreeNodes(verbose bool) []*tree.Tree {
	var nodes []*tree.Tree
	for _, gateway := range config.Gateways {
		childNodes := gateway.ToTreeNodes(verbose)
		nodes = append(nodes, childNodes...)
	}

	return nodes
}

func (gateway ApiGatewaySettings) ToTreeNodes(verbose bool) []*tree.Tree {
	var nodes []*tree.Tree
	node := tree.New().Root(styles.Tertiary.Bold(true).Render(*gateway.Name))

	if verbose {
		node.Child(styles.Primary.Render("Region:       ") + styles.Secondary.Render(helpers.PtrOrDefault(gateway.Region, "[region not set]")))
		node.Child(styles.Primary.Render("Protocol:     ") + styles.Secondary.Render(helpers.PtrOrDefault(gateway.Protocol, "[protocol not set]")))
		node.Child(styles.Primary.Render("Description:  ") + styles.Secondary.Render(helpers.PtrOrDefault(gateway.Description, "[description not set]")))
		stagesNode := tree.New().Root(styles.Primary.Render("Stages"))

		for _, stage := range *gateway.Stages {
			childNodes := generateApiGatewayStageNodes(&stage)
			for i := range childNodes {
				stagesNode.Child(childNodes[i])
			}
		}

		integrationsNode := tree.New().Root(styles.Primary.Render("Integrations"))
		for integationIndex, integration := range gateway.Integrations {
			childNodes := generateApiGatewayIntegrationNodes(&integration, integationIndex)
			for i := range childNodes {
				integrationsNode.Child(childNodes[i])
			}
		}

		routesNode := tree.New().Root(styles.Primary.Render("Routes"))
		for _, route := range gateway.Routes {
			childNodes := generateApiGatewayRouteNodes(&route)
			for i := range childNodes {
				routesNode.Child(childNodes[i])
			}
		}

		node.Child(stagesNode)
		node.Child(integrationsNode)
		node.Child(routesNode)
	}

	nodes = append(nodes, node)
	return nodes
}

func generateApiGatewayStageNodes(stage *ApiGatewayStage) []*tree.Tree {
	var nodes []*tree.Tree
	node := tree.New().Root(styles.Primary.Render(stage.Name))
	node.Child(styles.Primary.Render("Description:  ") + styles.Secondary.Render(helpers.PtrOrDefault(&stage.Description, "[description not set]")))
	node.Child(styles.Primary.Render("Auto-Deploy:  ") + styles.Secondary.Render(fmt.Sprintf("%t", stage.AutoDeploy)))

	nodes = append(nodes, node)
	return nodes
}

func generateApiGatewayIntegrationNodes(integration *ApiGatewayIntegration, iteration int) []*tree.Tree {
	var nodes []*tree.Tree
	node := tree.New().Root(styles.Primary.Render(fmt.Sprintf("Integration %d", iteration+1)))

	node.Child(styles.Primary.Render("Type:                ") + styles.Secondary.Render(integration.Type))
	node.Child(styles.Primary.Render("Target:              ") + styles.Secondary.Render(integration.Target.External.Dynamic.Name))
	node.Child(styles.Primary.Render("Payload Version:     ") + styles.Secondary.Render(integration.PayloadVersion))
	node.Child(styles.Primary.Render("Integration Method:  ") + styles.Secondary.Render(integration.IntegrationMethod))

	nodes = append(nodes, node)
	return nodes
}

func generateApiGatewayRouteNodes(route *ApiGatewayRoute) []*tree.Tree {
	var nodes []*tree.Tree
	node := tree.New().Root(styles.Primary.Render(route.Method + " " + route.Route))
	node.Child(styles.Primary.Render("Target:  ") + styles.Secondary.Render(*route.Target.Ref))

	nodes = append(nodes, node)
	return nodes
}
