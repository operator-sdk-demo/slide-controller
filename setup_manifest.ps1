$rootdir = Get-Location
$app = "operator"
$repository = "https://github.com/operator-sdk-demo/slide-config.git"
$CloneDir = ".tilt/manifests"
$Output = "deployment.yaml"
$namespace = "presentation-system"

$kustomizeDirectory = "$app"
$outputFile = "$CloneDir\deployment.yaml"

if (Test-path $CloneDir) {
  Set-Location -Path $CloneDir
  git pull
} else {
  git clone $repository $CloneDir
  Set-Location -Path $CloneDir
}


Set-Location -Path $kustomizeDirectory

# Run kustomize build and store the output in a file
kustomize edit set namespace $namespace
kustomize edit set image app="${namespace}-${app}:test"

kustomize build . > "$rootdir\$outputFile"
kustomize build

Get-Content "$rootdir\$outputFile"

# Change back to the original directory (optional)
Set-Location -Path $rootdir

