package main

import (
    )
    
func post_page (page string) {
    substitute("pages/"+page+".html", "in_draft", "posted")
}

func unpost_page (page string) {
    substitute("pages/"+page+".html", "posted", "in_draft")
}