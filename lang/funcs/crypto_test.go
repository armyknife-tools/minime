package funcs

import (
	"fmt"
	"testing"

	"github.com/zclconf/go-cty/cty"
	"golang.org/x/crypto/bcrypt"
)

func TestUUID(t *testing.T) {
	result, err := UUID()
	if err != nil {
		t.Fatal(err)
	}

	resultStr := result.AsString()
	if got, want := len(resultStr), 36; got != want {
		t.Errorf("wrong result length %d; want %d", got, want)
	}
}

func TestBase64Sha256(t *testing.T) {
	tests := []struct {
		String cty.Value
		Want   cty.Value
		Err    bool
	}{
		{
			cty.StringVal("test"),
			cty.StringVal("n4bQgYhMfWWaL+qgxVrQFaO/TxsrC4Is0V1sFbDwCgg="),
			false,
		},
		// This would differ because we're base64-encoding hex represantiation, not raw bytes.
		// base64encode(sha256("test")) =
		// "OWY4NmQwODE4ODRjN2Q2NTlhMmZlYWEwYzU1YWQwMTVhM2JmNGYxYjJiMGI4MjJjZDE1ZDZjMTViMGYwMGEwOA=="
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("base64sha256(%#v)", test.String), func(t *testing.T) {
			got, err := Base64Sha256(test.String)

			if test.Err {
				if err == nil {
					t.Fatal("succeeded; want error")
				}
				return
			} else if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if !got.RawEquals(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}

func TestBase64Sha512(t *testing.T) {
	tests := []struct {
		String cty.Value
		Want   cty.Value
		Err    bool
	}{
		{
			cty.StringVal("test"),
			cty.StringVal("7iaw3Ur350mqGo7jwQrpkj9hiYB3Lkc/iBml1JQODbJ6wYX4oOHV+E+IvIh/1nsUNzLDBMxfqa2Ob1f1ACio/w=="),
			false,
		},
		// This would differ because we're base64-encoding hex represantiation, not raw bytes
		// base64encode(sha512("test")) =
		// "OZWUyNmIwZGQ0YWY3ZTc0OWFhMWE4ZWUzYzEwYWU5OTIzZjYxODk4MDc3MmU0NzNmODgxOWE1ZDQ5NDBlMGRiMjdhYzE4NWY4YTBlMWQ1Zjg0Zjg4YmM4ODdmZDY3YjE0MzczMmMzMDRjYzVmYTlhZDhlNmY1N2Y1MDAyOGE4ZmY="
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("base64sha512(%#v)", test.String), func(t *testing.T) {
			got, err := Base64Sha512(test.String)

			if test.Err {
				if err == nil {
					t.Fatal("succeeded; want error")
				}
				return
			} else if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if !got.RawEquals(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}

func TestBcrypt(t *testing.T) {
	// single variable test
	p, err := Bcrypt(cty.StringVal("test"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(p.AsString()), []byte("test"))
	if err != nil {
		t.Fatalf("Error comparing hash and password: %s", err)
	}

	// testing with two parameters
	p, err = Bcrypt(cty.StringVal("test"), cty.NumberIntVal(5))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(p.AsString()), []byte("test"))
	if err != nil {
		t.Fatalf("Error comparing hash and password: %s", err)
	}

	// Negative test for more than two parameters
	_, err = Bcrypt(cty.StringVal("test"), cty.NumberIntVal(10), cty.NumberIntVal(11))
	if err == nil {
		t.Fatal("succeeded; want error")
	}
}

func TestMd5(t *testing.T) {
	tests := []struct {
		String cty.Value
		Want   cty.Value
		Err    bool
	}{
		{
			cty.StringVal("tada"),
			cty.StringVal("ce47d07243bb6eaf5e1322c81baf9bbf"),
			false,
		},
		{ // Confirm that we're not trimming any whitespaces
			cty.StringVal(" tada "),
			cty.StringVal("aadf191a583e53062de2d02c008141c4"),
			false,
		},
		{ // We accept empty string too
			cty.StringVal(""),
			cty.StringVal("d41d8cd98f00b204e9800998ecf8427e"),
			false,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("md5(%#v)", test.String), func(t *testing.T) {
			got, err := Md5(test.String)

			if test.Err {
				if err == nil {
					t.Fatal("succeeded; want error")
				}
				return
			} else if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if !got.RawEquals(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}

func TestRsaDecrypt(t *testing.T) {
	tests := []struct {
		Ciphertext cty.Value
		Privatekey cty.Value
		Want       cty.Value
		Err        bool
	}{
		// Base-64 encoded cipher decrypts correctly
		{
			cty.StringVal(CipherBase64),
			cty.StringVal(PrivateKey),
			cty.StringVal("message"),
			false,
		},
		// Wrong key
		{
			cty.StringVal(CipherBase64),
			cty.StringVal(WrongPrivateKey),
			cty.UnknownVal(cty.String),
			true,
		},
		// Bad key
		{
			cty.StringVal(CipherBase64),
			cty.StringVal("bad key"),
			cty.UnknownVal(cty.String),
			true,
		},
		// Empty key
		{
			cty.StringVal(CipherBase64),
			cty.StringVal(""),
			cty.UnknownVal(cty.String),
			true,
		},
		// Bad cipher
		{
			cty.StringVal("bad cipher"),
			cty.StringVal(PrivateKey),
			cty.UnknownVal(cty.String),
			true,
		},
		// Empty cipher
		{
			cty.StringVal(""),
			cty.StringVal(PrivateKey),
			cty.UnknownVal(cty.String),
			true,
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("RsaDecrypt(%#v, %#v)", test.Ciphertext, test.Privatekey), func(t *testing.T) {
			got, err := RsaDecrypt(test.Ciphertext, test.Privatekey)

			if test.Err {
				if err == nil {
					t.Fatal("succeeded; want error")
				}
				return
			} else if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if !got.RawEquals(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}

func TestSha1(t *testing.T) {
	tests := []struct {
		String cty.Value
		Want   cty.Value
		Err    bool
	}{
		{
			cty.StringVal("test"),
			cty.StringVal("a94a8fe5ccb19ba61c4c0873d391e987982fbbd3"),
			false,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("sha1(%#v)", test.String), func(t *testing.T) {
			got, err := Sha1(test.String)

			if test.Err {
				if err == nil {
					t.Fatal("succeeded; want error")
				}
				return
			} else if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if !got.RawEquals(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}

func TestSha256(t *testing.T) {
	tests := []struct {
		String cty.Value
		Want   cty.Value
		Err    bool
	}{
		{
			cty.StringVal("test"),
			cty.StringVal("9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"),
			false,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("sha256(%#v)", test.String), func(t *testing.T) {
			got, err := Sha256(test.String)

			if test.Err {
				if err == nil {
					t.Fatal("succeeded; want error")
				}
				return
			} else if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if !got.RawEquals(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}

func TestSha512(t *testing.T) {
	tests := []struct {
		String cty.Value
		Want   cty.Value
		Err    bool
	}{
		{
			cty.StringVal("test"),
			cty.StringVal("ee26b0dd4af7e749aa1a8ee3c10ae9923f618980772e473f8819a5d4940e0db27ac185f8a0e1d5f84f88bc887fd67b143732c304cc5fa9ad8e6f57f50028a8ff"),
			false,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("sha512(%#v)", test.String), func(t *testing.T) {
			got, err := Sha512(test.String)

			if test.Err {
				if err == nil {
					t.Fatal("succeeded; want error")
				}
				return
			} else if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if !got.RawEquals(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}
}

const (
	CipherBase64 = "eczGaDhXDbOFRZGhjx2etVzWbRqWDlmq0bvNt284JHVbwCgObiuyX9uV0LSAMY707IEgMkExJqXmsB4OWKxvB7epRB9G/3+F+pcrQpODlDuL9oDUAsa65zEpYF0Wbn7Oh7nrMQncyUPpyr9WUlALl0gRWytOA23S+y5joa4M34KFpawFgoqTu/2EEH4Xl1zo+0fy73fEto+nfkUY+meuyGZ1nUx/+DljP7ZqxHBFSlLODmtuTMdswUbHbXbWneW51D7Jm7xB8nSdiA2JQNK5+Sg5x8aNfgvFTt/m2w2+qpsyFa5Wjeu6fZmXSl840CA07aXbk9vN4I81WmJyblD/ZA=="
	PrivateKey   = `
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAgUElV5mwqkloIrM8ZNZ72gSCcnSJt7+/Usa5G+D15YQUAdf9
c1zEekTfHgDP+04nw/uFNFaE5v1RbHaPxhZYVg5ZErNCa/hzn+x10xzcepeS3KPV
Xcxae4MR0BEegvqZqJzN9loXsNL/c3H/B+2Gle3hTxjlWFb3F5qLgR+4Mf4ruhER
1v6eHQa/nchi03MBpT4UeJ7MrL92hTJYLdpSyCqmr8yjxkKJDVC2uRrr+sTSxfh7
r6v24u/vp/QTmBIAlNPgadVAZw17iNNb7vjV7Gwl/5gHXonCUKURaV++dBNLrHIZ
pqcAM8wHRph8mD1EfL9hsz77pHewxolBATV+7QIDAQABAoIBAC1rK+kFW3vrAYm3
+8/fQnQQw5nec4o6+crng6JVQXLeH32qXShNf8kLLG/Jj0vaYcTPPDZw9JCKkTMQ
0mKj9XR/5DLbBMsV6eNXXuvJJ3x4iKW5eD9WkLD4FKlNarBRyO7j8sfPTqXW7uat
NxWdFH7YsSRvNh/9pyQHLWA5OituidMrYbc3EUx8B1GPNyJ9W8Q8znNYLfwYOjU4
Wv1SLE6qGQQH9Q0WzA2WUf8jklCYyMYTIywAjGb8kbAJlKhmj2t2Igjmqtwt1PYc
pGlqbtQBDUiWXt5S4YX/1maIQ/49yeNUajjpbJiH3DbhJbHwFTzP3pZ9P9GHOzlG
kYR+wSECgYEAw/Xida8kSv8n86V3qSY/I+fYQ5V+jDtXIE+JhRnS8xzbOzz3v0WS
Oo5H+o4nJx5eL3Ghb3Gcm0Jn46dHrxinHbm+3RjXv/X6tlbxIYjRSQfHOTSMCTvd
qcliF5vC6RCLXuc7R+IWR1Ky6eDEZGtrvt3DyeYABsp9fRUFR/6NluUCgYEAqNsw
1aSl7WJa27F0DoJdlU9LWerpXcazlJcIdOz/S9QDmSK3RDQTdqfTxRmrxiYI9LEs
mkOkvzlnnOBMpnZ3ZOU5qIRfprecRIi37KDAOHWGnlC0EWGgl46YLb7/jXiWf0AG
Y+DfJJNd9i6TbIDWu8254/erAS6bKMhW/3q7f2kCgYAZ7Id/BiKJAWRpqTRBXlvw
BhXoKvjI2HjYP21z/EyZ+PFPzur/lNaZhIUlMnUfibbwE9pFggQzzf8scM7c7Sf+
mLoVSdoQ/Rujz7CqvQzi2nKSsM7t0curUIb3lJWee5/UeEaxZcmIufoNUrzohAWH
BJOIPDM4ssUTLRq7wYM9uQKBgHCBau5OP8gE6mjKuXsZXWUoahpFLKwwwmJUp2vQ
pOFPJ/6WZOlqkTVT6QPAcPUbTohKrF80hsZqZyDdSfT3peFx4ZLocBrS56m6NmHR
UYHMvJ8rQm76T1fryHVidz85g3zRmfBeWg8yqT5oFg4LYgfLsPm1gRjOhs8LfPvI
OLlRAoGBAIZ5Uv4Z3s8O7WKXXUe/lq6j7vfiVkR1NW/Z/WLKXZpnmvJ7FgxN4e56
RXT7GwNQHIY8eDjDnsHxzrxd+raOxOZeKcMHj3XyjCX3NHfTscnsBPAGYpY/Wxzh
T8UYnFu6RzkixElTf2rseEav7rkdKkI3LAeIZy7B0HulKKsmqVQ7
-----END RSA PRIVATE KEY-----
`
	WrongPrivateKey = `
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAlrCgnEVgmNKCq7KPc+zUU5IrxPu1ClMNJS7RTsTPEkbwe5SB
p+6V6WtCbD/X/lDRRGbOENChh1Phulb7lViqgrdpHydgsrKoS5ah3DfSIxLFLE00
9Yo4TCYwgw6+s59j16ZAFVinaQ9l6Kmrb2ll136hMrz8QKh+qw+onOLd38WFgm+W
ZtUqSXf2LANzfzzy4OWFNyFqKaCAolSkPdTS9Nz+svtScvp002DQp8OdP1AgPO+l
o5N3M38Fftapwg0pCtJ5Zq0NRWIXEonXiTEMA6zy3gEZVOmDxoIFUWnmrqlMJLFy
5S6LDrHSdqJhCxDK6WRZj43X9j8spktk3eGhMwIDAQABAoIBAAem8ID/BOi9x+Tw
LFi2rhGQWqimH4tmrEQ3HGnjlKBY+d1MrUjZ1MMFr1nP5CgF8pqGnfA8p/c3Sz8r
K5tp5T6+EZiDZ2WrrOApxg5ox0MAsQKO6SGO40z6o3wEQ6rbbTaGOrraxaWQIpyu
AQanU4Sd6ZGqByVBaS1GnklZO+shCHqw73b7g1cpLEmFzcYnKHYHlUUIsstMe8E1
BaCY0CH7JbWBjcbiTnBVwIRZuu+EjGiQuhTilYL2OWqoMVg1WU0L2IFpR8lkf/2W
SBx5J6xhwbBGASOpM+qidiN580GdPzGhWYSqKGroHEzBm6xPSmV1tadNA26WFG4p
pthLiAECgYEA5BsPRpNYJAQLu5B0N7mj9eEp0HABVEgL/MpwiImjaKdAwp78HM64
IuPvJxs7r+xESiIz4JyjR8zrQjYOCKJsARYkmNlEuAz0SkHabCw1BdEBwUhjUGVB
efoERK6GxfAoNqmSDwsOvHFOtsmDIlbHmg7G2rUxNVpeou415BSB0B8CgYEAqR4J
YHKk2Ibr9rU+rBU33TcdTGw0aAkFNAVeqM9j0haWuFXmV3RArgoy09lH+2Ha6z/g
fTX2xSDAWV7QUlLOlBRIhurPAo2jO2yCrGHPZcWiugstrR2hTTInigaSnCmK3i7F
6sYmL3S7K01IcVNxSlWvGijtClT92Cl2WUCTfG0CgYAiEjyk4QtQTd5mxLvnOu5X
oqs5PBGmwiAwQRiv/EcRMbJFn7Oupd3xMDSflbzDmTnWDOfMy/jDl8MoH6TW+1PA
kcsjnYhbKWwvz0hN0giVdtOZSDO1ZXpzOrn6fEsbM7T9/TQY1SD9WrtUKCNTNL0Z
sM1ZC6lu+7GZCpW4HKwLJwKBgQCRT0yxQXBg1/UxwuO5ynV4rx2Oh76z0WRWIXMH
S0MyxdP1SWGkrS/SGtM3cg/GcHtA/V6vV0nUcWK0p6IJyjrTw2XZ/zGluPuTWJYi
9dvVT26Vunshrz7kbH7KuwEICy3V4IyQQHeY+QzFlR70uMS0IVFWAepCoWqHbIDT
CYhwNQKBgGPcLXmjpGtkZvggl0aZr9LsvCTckllSCFSI861kivL/rijdNoCHGxZv
dfDkLTLcz9Gk41rD9Gxn/3sqodnTAc3Z2PxFnzg1Q/u3+x6YAgBwI/g/jE2xutGW
H7CurtMwALQ/n/6LUKFmjRZjqbKX9SO2QSaC3grd6sY9Tu+bZjLe
-----END RSA PRIVATE KEY-----
`
)
