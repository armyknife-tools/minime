resource "test_object" "A" {
  lifecycle {
    create_before_destroy = true
  }
}

resource "test_object" "B" {
  test_string = "${test_object.A.id}"
}
