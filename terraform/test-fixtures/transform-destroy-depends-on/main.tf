resource "aws_instance" "foo" {}

resource "aws_instance" "bar" {
    depends_on = ["aws_instance.foo"]
}
