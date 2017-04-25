#ps1_sysnative
$domainPass = "@ppr3nda"
$dsrmPassword = (ConvertTo-SecureString -AsPlainText -Force -String $domainPass)
$domainName = "acp"
$domainSuf = "local"
$LocalAdmin = "Administrator"
$objUser = [ADSI]"WinNT://localhost/$($LocalAdmin), user"
$objUser.psbase.Invoke("SetPassword", $domainPass)
New-Item "C:\\builder\\temp" -type directory -force
$webClient = New-Object System.Net.WebClient
Install-WindowsFeature -name AD-Domain-Services -IncludeManagementTools | Out-Null
$dsrmPassword = (ConvertTo-SecureString -AsPlainText -Force -String $domainPass)
Install-ADDSForest -DomainName "$($domainName).$($domainSuf)" -InstallDNS -Force -SafeModeAdministratorPassword $dsrmPassword -ForestMode Win2012R2 -DomainMode Win2012R2 | Out-Null
