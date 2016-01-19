package main

import (
    "flag"
)

func main() {
    /*create_site("testsite1");*/
    
    createSitePtr := flag.String("create", "default_site", "Use this command to create a new site structure.")
    flag.Parse()
    
    create_site(*createSitePtr)
}