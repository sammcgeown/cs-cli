/*
Copyright © 2021 Sam McGeown <smcgeown@vmware.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// getEndpointCmd represents the endpoint command
var getEndpointCmd = &cobra.Command{
	Use:   "endpoint",
	Short: "Get Code Stream Endpoint Configurations",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ensureTargetConnection()

		response, err := getEndpoint(id, name, project, typename, export, exportPath)
		if err != nil {
			fmt.Print("Unable to get endpoints: ", err)
		}
		var resultCount = len(response)
		if resultCount == 0 {
			// No results
			fmt.Println("No results found")
		} else if resultCount == 1 {
			PrettyPrint(response[0])
		} else {
			// Print result table
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "Project", "Type", "Description"})
			for _, c := range response {
				table.Append([]string{c.Name, c.Project, c.Type, c.Description})
			}
			table.Render()
		}

	},
}

// createEndpointCmd represents the endpoint create command
var createEndpointCmd = &cobra.Command{
	Use:   "endpoint",
	Short: "Create an Endpoint",
	Long: `Create an Endpoint by importing a YAML specification.
	
	Create from YAML
	  cs-cli create endpoint --importPath "/Users/sammcgeown/Desktop/endpoint.yaml"
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ensureTargetConnection()
		if importPath != "" {
			if importYaml(importPath, "create") {
				fmt.Println("Imported successfully, Endpoint created.")
			}
		}
	},
}

// updateEndpointCmd represents the endpoint update command
var updateEndpointCmd = &cobra.Command{
	Use:   "endpoint",
	Short: "Update an Endpoint",
	Long: `Update an Endpoint by importing the YAML specification

	Update from YAML
	cs-cli update endpoint --importPath "/Users/sammcgeown/Desktop/updated-endpoint.yaml"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		ensureTargetConnection()
		if importPath != "" {
			if importYaml(importPath, "apply") {
				fmt.Println("Imported successfully, pipeline updated.")
			}
		}
	},
}

func init() {
	getCmd.AddCommand(getEndpointCmd)
	getEndpointCmd.Flags().StringVarP(&name, "name", "n", "", "Get Endpoint by Name")
	getEndpointCmd.Flags().StringVarP(&id, "id", "i", "", "Get Endpoint by ID")
	getEndpointCmd.Flags().StringVarP(&project, "project", "p", "", "Filter Endpoint by Project")
	getEndpointCmd.Flags().StringVarP(&typename, "type", "t", "", "Filter Endpoint by Type")
	getEndpointCmd.Flags().StringVarP(&exportPath, "exportPath", "", "", "Path to export objects - relative or absolute location")
	getEndpointCmd.Flags().BoolVarP(&export, "export", "e", false, "Export Endpoint")
	// Create
	createCmd.AddCommand(createEndpointCmd)
	createEndpointCmd.Flags().StringVarP(&importPath, "importPath", "c", "", "YAML configuration file to import")
	createEndpointCmd.MarkFlagRequired("importPath")
	// Update
	updateCmd.AddCommand(updateEndpointCmd)
	updateEndpointCmd.Flags().StringVarP(&importPath, "importPath", "c", "", "YAML configuration file to import")
	updateEndpointCmd.MarkFlagRequired("importPath")

}
