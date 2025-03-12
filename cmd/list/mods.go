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
	"fmt"

	"github.com/donovanmods/projectdaedalus-db-tool/lib/firestore"
	"github.com/donovanmods/projectdaedalus-db-tool/lib/logger"
	"github.com/donovanmods/projectdaedalus-db-tool/lib/mod"
	"github.com/spf13/cobra"
)

// delCmd represents the del command
var modsCmd = &cobra.Command{
	Use:   "mods",
	Short: "Display a list of all mods",
	Run: func(cmd *cobra.Command, args []string) {
		doModsList(cmd, args, firestore.ModList)
	},
}

func init() {
	ListCmd.AddCommand(modsCmd)
}

func doModsList(cmd *cobra.Command, args []string, collection func() (firestore.DBList[mod.Mod], error)) {
	_ = args // Unused for now -- will be used for Searching later

	mods, err := collection()
	if err != nil {
		logger.Fatal(err)
	}

	if mods == nil {
		return
	}

	if mods.Count() == 0 {
		logger.Warn("No items returned")
		return
	}

	if jsonFlag, _ := cmd.Flags().GetBool("json"); jsonFlag {
		if j, err := mods.MarshalJSON(); err != nil {
			logger.Fatal(err)
		} else {
			fmt.Println(string(j))
		}
	} else {
		fmt.Println(mods.String())
	}
}
