package main

import "github.com/nint8835/terraform-gatsby-service/service"

func main() {
	service.GetRouter().Run(":9000")
}
