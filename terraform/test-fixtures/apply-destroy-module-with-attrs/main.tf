resource "aws_instance" "vpc"   { }

module "child" {
  source = "./child"
  vpc_id = "${aws_instance.vpc.id}"
}
