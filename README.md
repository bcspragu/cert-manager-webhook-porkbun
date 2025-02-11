<p align="center">
  <img src="https://raw.githubusercontent.com/cert-manager/cert-manager/d53c0b9270f8cd90d908460d69502694e1838f5f/logo/logo-small.png" height="256" width="256" alt="cert-manager project logo" />
</p>

# Porkbun Webhook for cert-manager

#### ⚠️ Attention: Project No Longer Maintained

This repository is no longer actively maintained. A fork with an additional helm chart and more documentation can be found [here](https://github.com/Talinx/cert-manager-webhook-porkbun).    

An implementation of the cert-manager [`webhook.Solver` interface](https://pkg.go.dev/github.com/cert-manager/cert-manager@v1.12.3/pkg/acme/webhook#Solver) for [Porkbun](https://porkbun.com/). This is based on [cert-manager/webhook-example](https://github.com/cert-manager/webhook-example), with inspiration from [baarde/cert-manager-webhook-ovh](https://github.com/baarde/cert-manager-webhook-ovh)

Note: The test suite does work, but I straight up deleted `main_test.go` because the dependency on `github.com/cert-manager/cert-manager/test/acme` was giving me insane, impossible to resolve versioning conflicts. I'm sure these will be resolved by someone more knowledgeable updating the `go.mod` in the example webhook, at which point I'll add the tests back.

### Running the test suite

All DNS providers **must** run the DNS01 provider conformance testing suite,
else they will have undetermined behaviour when used with cert-manager.

**It is essential that you configure and run the test suite when creating a
DNS01 webhook.**

An example Go test file has been provided in [main_test.go](https://github.com/bcspragu/cert-manager-webhook-porkbun/blob/master/main_test.go).

You can run the test suite with:

```bash
$ TEST_ZONE_NAME=example.com. make test
```

The example file has a number of areas you must fill in and replace with your
own options in order for tests to pass.
