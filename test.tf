resource "gatsby_text_bold" "test" {
    text = "Hello world!"
}

output "contents" {
    value = gatsby_text_bold.test.contents
}