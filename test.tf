resource "gatsby_text_bold" "test" {
    text = "Hello world!"

    provisioner "local-exec" {
        command = "echo test"
    }
}

output "contents" {
    value = gatsby_text_bold.test.contents
}
