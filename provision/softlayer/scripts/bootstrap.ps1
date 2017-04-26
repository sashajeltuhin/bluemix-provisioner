$domainPass = "@ppr3nda"
$dsrmPassword = (ConvertTo-SecureString -AsPlainText -Force -String $domainPass)
$domainName = "acp"
$domainSuf = "local"
$LocalAdmin = "Administrator"
$objUser = [ADSI]"WinNT://localhost/$($LocalAdmin), user"
$objUser.psbase.Invoke("SetPassword", $domainPass)
Enable-PSRemoting –force
Set-Item wsman:\localhost\client\trustedhosts *
Restart-Service WinRM