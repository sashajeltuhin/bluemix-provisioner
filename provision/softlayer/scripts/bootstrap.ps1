$domainPass = "@ppr3nda"

$LocalAdmin = "Administrator"

$objUser = [ADSI]"WinNT://localhost/$($LocalAdmin), user"

$objUser.psbase.Invoke("SetPassword", $domainPass)

Enable-PSRemoting -Force

Set-Item wsman:\localhost\client\trustedhosts *

winrm set winrm/config/service/auth @{Basic="true"}

winrm set winrm/config/service @{AllowUnencrypted="true"}

Restart-Service WinRM
