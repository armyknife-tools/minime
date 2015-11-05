resource "test_instance" "foo" {
    ami = "bar"

    network_interface {
      device_index = 0
      name = test
      description = "Main network interface"
    }
}
