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
package addCmd

import (
	"errors"
	"fmt"

	"github.com/donovanmods/projectdaedalus-db-tool/lib/firestore"
	"github.com/donovanmods/projectdaedalus-db-tool/lib/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// addCmd represents the add command
var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add entries to the database",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func doMetaAdd(cmd *cobra.Command, args []string, collection func() (firestore.DBList[string], error)) {
	_ = cmd // Unused for now

	if len(args) == 0 {
		logger.Fatal(errors.New("no item given to add"))
	}

	meta, err := collection()
	if err != nil {
		logger.Fatal(err)
	}

	if meta == nil {
		logger.Warn("Collection not found")
		return
	}

	for _, item := range args {
		if err = meta.Add(item); err != nil {
			if errors.Is(err, firestore.ErrDuplicate) {
				continue
			}
			logger.Fatal(err)
		}
	}

	if _, err := meta.Commit(); err != nil {
		logger.Fatal(err)
	}

	if viper.GetInt("verbosity") > 0 {
		fmt.Println(meta.String())
	}
}
