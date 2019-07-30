package cmd

import (
	"log"
	"os"

	"github.com/alexkaplun/imagediffvisualizer/service"

	"github.com/spf13/cobra"
)

// vars to store provided filenames
var (
	image1 string
	image2 string
	output string
)

var RootCmd = &cobra.Command{
	Use:   "imagediff -1 | --image1 <img1.png> -2 | --image2 <img2.png> -o | --output <output.png>",
	Short: "Provides visualization of greyscale images difference",
	Long:  `Provides visualization of greyscale images difference`,
	Run: func(cmd *cobra.Command, args []string) {
		// instantiante comparer and validate images
		comparer, err := service.NewComparer(image1, image2, output)
		if err != nil {
			log.Println(err)
			return
		}
		// do the comparison
		if err = comparer.Compare(); err != nil {
			log.Println(err)
			return
		}
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(-1)
	}
}

// Initialization to parse flags
func init() {
	RootCmd.Flags().StringVarP(&image1, "image1", "1", "", "first image")
	RootCmd.MarkFlagRequired("image1")
	RootCmd.Flags().StringVarP(&image2, "image2", "2", "", "second image")
	RootCmd.MarkFlagRequired("image2")
	RootCmd.Flags().StringVarP(&output, "output", "o", "", "output file")
	RootCmd.MarkFlagRequired("output")
}
