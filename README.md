# ase-3-6-attendance-tracking-web
## Install
Make sure you add the repository under 
`$GOPATH/src/github.com/johannesrohwer/ase-3-6-attendance-tracking-web`. 
~~In order to run the project locally add the required dependencies via `dep ensure`
([dep documentation](http://dep.com))~~ (use go get for now) and execute: 

`dev_appserver.py api/app.yaml`

`dev_appserver.py` is available as part of the `gcloud` (including `google-cloud-sdk-app-engine-go`)
package that is available [here](https://cloud.google.com/sdk/docs/). AFAIK `python2.7` is required to run it.

## Deploy
To deploy the application to app engine make sure that you initialized `gcloud` and run

`glcoud app deploy api/app.yaml`
## Routes
### /
Returns the `index.html` template.

### /api/version
Returns a JSON encoded version object that includes authors and version number.

## Frameworks and Libraries
### GorillaMux
[GorillaMux](http://www.gorillatoolkit.org/pkg/mux) adds extended functionality to URL routing
such as regular expressions.