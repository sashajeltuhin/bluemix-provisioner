$domainPass = "{{domainPass}}"

$dcip = "{{dcip}}"

$platformadmin = "{{platformadmin}}"

$platformsystem = "{{platformsystem}}"

$domainName = "{{domainName}}"

$domainSuf = "{{domainSuf}}"

$LocalAdmin = "Administrator"

$dsrmPassword = (ConvertTo-SecureString -AsPlainText -Force -String $domainPass)

$objUser = [ADSI]"WinNT://localhost/$($LocalAdmin), user"

$objUser.psbase.Invoke("SetPassword", $domainPass)

$password = $domainPass | ConvertTo-SecureString -asPlainText -Force

$username = "$($domainName)\Administrator"

$credential = New-Object System.Management.Automation.PSCredential($username,$password)

Set-DnsClientServerAddress -InterfaceAlias "Ethernet" -ServerAddresses $dcip

Add-Computer -DomainName "$($domainName).$($domainSuf)" -Credential $credential | Out-Null

New-ADUser -Credential $credential -Server "$($dcname).$($domainName).$($domainSuf)" -SamAccountName $platformadmin -AccountPassword $dsrmPassword -name "$($platformadmin)" -enabled $true -PasswordNeverExpires $true -ChangePasswordAtLogon $false | Out-Null

New-ADUser -Credential $credential -Server "$($dcname).$($domainName).$($domainSuf)" -SamAccountName $platformsystem -AccountPassword $dsrmPassword -name "$($platformsystem)" -enabled $true -PasswordNeverExpires $true -ChangePasswordAtLogon $false | Out-Null

Add-ADPrincipalGroupMembership -Identity "CN=$($platformadmin),CN=Users,DC=$($domainName),DC=$($domainSuf)" -MemberOf "CN=Domain Admins,CN=Users,DC=$($domainName),DC=$($domainSuf)"

Add-ADPrincipalGroupMembership -Identity "CN=$($platformsystem),CN=Users,DC=$($domainName),DC=$($domainSuf)" -MemberOf "CN=Domain Admins,CN=Users,DC=$($domainName),DC=$($domainSuf)"

Install-WindowsFeature -name AD-Domain-Services -IncludeManagementTools | Out-Null

Install-ADDSDomainController -Credential $credential -DomainName "$($domainName).$($domainSuf)" -InstallDns -SysvolPath -NoRebootOnCompletion -Force
