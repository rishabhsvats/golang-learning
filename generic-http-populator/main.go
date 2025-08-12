package main

import "flag"

func main() {
	var image string
	var namespace string
	var mode string
	var uri string

	flag.StringVar(&image, "image", "", "Image for populator component")
	flag.StringVar(&namespace, "namespace", "", "Namespace for populator component")
	flag.StringVar(&mode, "mode", "", "Mode to run the application in")
	flag.StringVar(&uri, "uri", "", "URI for the content of the volume")

	flag.Parse()
}
