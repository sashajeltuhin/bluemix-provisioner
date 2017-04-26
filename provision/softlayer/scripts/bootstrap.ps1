$domainPass = "@ppr3nda"

$LocalAdmin = "Administrator"

$objUser = [ADSI]"WinNT://localhost/$($LocalAdmin), user"

$objUser.psbase.Invoke("SetPassword", $domainPass)

Enable-PSRemoting -Force

Set-Item wsman:\localhost\client\trustedhosts *

Restart-Service WinRM
