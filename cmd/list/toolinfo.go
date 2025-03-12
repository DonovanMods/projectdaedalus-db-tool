/*
Copyright © 2025 Donovan C. Young

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package listCmd

import (
	"github.com/donovanmods/projectdaedalus-db-tool/lib/firestore"
	"github.com/spf13/cobra"
)

// delCmd represents the del command
var toolinfoCmd = &cobra.Command{
	Use:   "toolinfo",
	Short: "Display a list of all toolinfo entries",
	Run: func(cmd *cobra.Command, args []string) {
		doMetaList(cmd, args, firestore.ToolInfo)
	},
}

func init() {
	ListCmd.AddCommand(toolinfoCmd)
}
