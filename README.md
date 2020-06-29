# ![realCoverage](images/realcoverage_logo.png "realCoverage")

CLI tools which is used to generate PDF report for domains presence in Akamai property configurations and check if domain resolves to Akamai IP.

__Under development__
Please forgive us for some error or issue.

## Dependencies

* [Akamai Overview CLI](https://github.com/apiheat/akamai-cli-overview) - to get account overview data
* [Akamai Property CLI](https://github.com/akamai/cli-property) - to get property configuration data
* [Akamai Diagnostic Tools CLI](https://github.com/apiheat/akamai-cli-diagnostic-tools) - to check if IP belongs to Akamai
* dig shell command

## Akamai Permissions

* Read Akamai contract groups
* List Akamai Properties
* Read Akamai Properties
* Use Akamai Diagnostic Tools

## Usage

### Setup edgerc credentials location

```shell
> export AKAMAI_EDGERC_SECTION="default"
> export AKAMAI_EDGERC_CONFIG="~/.edgerc"
```

### Generate report

```shell
> ./realcoverage -f input.yml -l logo.png -o report.pdf -c "My Company" -t "+0 12345678" -s "mycompany.nl"

Usage of realcoverage
  -c string
      Company name.
  -f string
      YAML file to parse.
  -l string
      Company logo file.
  -o string
      Output file name.
  -s string
      Company web site.
  -t string
      Company phone.
```
