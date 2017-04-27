$domainPass = "{{domainPass}}"
$dsrmPassword = (ConvertTo-SecureString -AsPlainText -Force -String $domainPass)
$domainName = "acp"
$domainSuf = "local"
$LocalAdmin = "Administrator"
$objUser = [ADSI]"WinNT://localhost/$($LocalAdmin), user"
$objUser.psbase.Invoke("SetPassword", $domainPass)
Install-WindowsFeature -name AD-Domain-Services -IncludeManagementTools | Out-Null
$dsrmPassword = (ConvertTo-SecureString -AsPlainText -Force -String $domainPass)
Install-ADDSForest -DomainName "$($domainName).$($domainSuf)" -InstallDNS -Force -SafeModeAdministratorPassword $dsrmPassword -ForestMode Win2012R2 -DomainMode Win2012R2 | Out-Null
Enable-PSRemoting –force
Set-Item wsman:\localhost\client\trustedhosts *
Restart-Service WinRM
