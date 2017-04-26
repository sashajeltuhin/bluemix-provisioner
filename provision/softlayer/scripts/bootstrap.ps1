$domainPass = "@ppr3nda" \n\r
$LocalAdmin = "Administrator" \n\r
$objUser = [ADSI]"WinNT://localhost/$($LocalAdmin), user" \n\r
$objUser.psbase.Invoke("SetPassword", $domainPass) \n\r
Enable-PSRemoting –force \n\r
Set-Item wsman:\localhost\client\trustedhosts * \n\r
Restart-Service WinRM \n\r