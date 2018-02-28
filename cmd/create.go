// Copyright © 2018 Patrick Motard <motard19@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	// "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/gob"
	"encoding/pem"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var name string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		usr, err := user.Current()
		check(err)
		var publicKeyName string
		var privateKeyName string

		// parse name and assign public/private names accordingly
		if strings.HasSuffix(name, ".pub") {
			publicKeyName = name
			privateKeyName = strings.TrimSuffix(name, filepath.Ext(name))
		} else {
			publicKeyName = strings.Join([]string{name, "pub"}, ".")
			privateKeyName = name
		}

		publicKeyFilePath := strings.Join([]string{usr.HomeDir, ".ssh", publicKeyName}, "/")
		privateKeyFilePath := strings.Join([]string{usr.HomeDir, ".ssh", privateKeyName}, "/")

		if fileExists(publicKeyFilePath) {
			fmt.Println(fmt.Sprintf("Error: public key %s already exists.", publicKeyName))
		}

		if fileExists(privateKeyFilePath) {
			fmt.Println(fmt.Sprintf("Error: private key %s already exists.", privateKeyName))
		}

		if fileExists(privateKeyFilePath) || fileExists(publicKeyFilePath) {
			os.Exit(1)
		}

		fmt.Println("public", publicKeyName)
		fmt.Println("private", privateKeyName)

		fmt.Println("create called")
	},
}

func fileExists(filepath string) bool {
	if _, err := os.Stat(filepath); err == nil {
		return true
	} else {
		return false
	}
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&name, "name", "n", "", "name of key to create")
	createCmd.MarkFlagRequired("name")

	// fmt.Println(strings.Join([]string{name, "key"}, "."))

	// reader := rand.Reader
	// bitsize := 4096
	// key, err := rsa.GenerateKey(reader, bitsize)
	// checkError(err)
	// fmt.Println(strings.Join([]string{name, "key"}, "."))

	// publicKey := key.PublicKey

	// saveGobKey("private.key", key)
	// saveGobKey("public.key", publicKey)
	// saveGobKey(strings.Join([]string{name, "key"}, "."), key)
	// savePEMKey("private.pem", key)

	// saveGobKey("public.key", publicKey)
	// savePublicPEMKey("public.pem", publicKey)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func saveGobKey(fileName string, key interface{}) {
	outFile, err := os.Create(fileName)
	checkError(err)
	defer outFile.Close()

	encoder := gob.NewEncoder(outFile)
	err = encoder.Encode(key)
	checkError(err)
}

func savePEMKey(fileName string, key *rsa.PrivateKey) {
	outFile, err := os.Create(fileName)
	checkError(err)
	defer outFile.Close()

	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(outFile, privateKey)
	checkError(err)
}

func savePublicPEMKey(fileName string, pubkey rsa.PublicKey) {
	asn1Bytes, err := asn1.Marshal(pubkey)
	checkError(err)

	var pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	pemfile, err := os.Create(fileName)
	checkError(err)
	defer pemfile.Close()

	err = pem.Encode(pemfile, pemkey)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
