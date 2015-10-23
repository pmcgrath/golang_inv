<#
	usage
		# List all lib repos - will prompt for your stash password and will use your login name as the stash user name
		.\run.ps1 list lib

		# Clone all lib repos - will prompt for your stash password and will use your login name as the stash user name
		# Expects the root directory to exist and will skip where a sub directory with the same name as a repo exists - assumes already cloned
		(test-path c:\tmp\repos\lib) -or (mkdir c:\tmp\repos\lib)
		.\run.ps1 clone lib -rootDirectoryPath c:\tmp\repos

		# Remotes for all lib repos
		.\run.ps1 remote lib										# Using standard repo location convention
		.\run.ps1 remote lib -rootDirectoryPath c:\tmp\repos

		# Branches for all lib repos
		.\run.ps1 branch lib										# Using standard repo location convention
		.\run.ps1 branch lib -rootDirectoryPath c:\tmp\repos

		# Status for all lib repos
		.\run.ps1 status lib										# Using standard repo location convention
		.\run.ps1 status lib -rootDirectoryPath c:\tmp\repos
#>
param
(
	[string] $action = 'status',
	[string] $repoType = 'ser',
	[string] $remoteName = 'upstream',
	[string[]] $exclusions = $null,
	[string] $rootDirectoryPath = 'c:\repos\stash',
	[string] $url = 'http://stash'
)

# Strict mode and stop on errors
set-strictmode -version latest
$ErrorActionPreference = 'Stop'

# Powershell version pre-requisite
if ($host.Version.Major -lt 4) { throw 'Powershell version must be >= 4.0' }

# Danger - need to think of whatifs so we do not loose work
pushd .

try
{
	# Consistent version
	$action = $action.Trim().ToLower()

	if (@('clone', 'list') -contains $action)
	{
		# Just using secure string for the asterisks in the shell, not too bothered about secret as should be short lived
		# See http://www.stockenterprises.ca/post/Convert-a-SecureString-to-a-String-%28Plain-Text%29-in-PowerShell.aspx
		$env:REPO_PASSWORD = [System.Runtime.InteropServices.Marshal]::PtrToStringAuto([System.Runtime.InteropServices.Marshal]::SecureStringToBSTR((read-host 'Gimme your stash password' -assecurestring)));
	}

	$rootDirectoryPath = join-path $rootDirectoryPath $repoType

	# Assumes the current branch is tracking the correct remote branch
	switch ($action.Trim().ToLower())
	{
		'branch' 	{ ./repomgr.exe branch -projectsdirectorypath $rootDirectoryPath }
		'clone' 	{ ./repomgr.exe clone -provider stash -parentname $repoType -url $url -projectsdirectorypath $rootDirectoryPath }
		'fetch'  	{ ./repomgr.exe fetch -projectsdirectorypath $rootDirectoryPath -remotename $remoteName }
		'list' 		{ ./repomgr.exe list -verbose -provider stash -parentname $repoType -url $url }
		'pull'  	{ ./repomgr.exe pull -projectsdirectorypath $rootDirectoryPath -remotename $remoteName }
		'remote' 	{ ./repomgr.exe remote -projectsdirectorypath $rootDirectoryPath }
		'status' 	{ ./repomgr.exe status -projectsdirectorypath $rootDirectoryPath }
	}
}
finally
{
	popd
	$env:REPO_PASSWORD = $null
}

