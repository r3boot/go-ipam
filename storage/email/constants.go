package email

const SignupTemplate string = `From: {{.SenderName}} <{{.Sender}}>
To: {{.Fullname}} <{{.Recipient}}>
Subject: Please activate your account on {{.NetworkName}}

Dear {{.Fullname}},

Thank you for signing up on {{.NetworkName}}, the premier tier-1 network in
your neighborhood. To complete your registration please follow the link below
in order to activate your account.

https://blah.blah/v1/activate/{{.Token}}

Once you have activated your account, you will be able to define your network
details.

Kind Regards,

{{.SenderName}}`
