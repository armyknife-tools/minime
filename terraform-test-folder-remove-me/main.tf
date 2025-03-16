provider "local" {
  version = ">1.0"
}

resource "local_file" "hello" {
  content  = "Hello world!"
  filename = "hello.txt"
}
