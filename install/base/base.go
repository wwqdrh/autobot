package base

import "github.com/wwqdrh/autobot/internal/sshcmd/sshutil"

func Install(client *sshutil.SSH, ip string) error {
	return client.CmdAsync(ip, "rm -rf autobot; mkdir -p /usr/local/autobot; chmod +x /usr/local/autobot; cd /usr/local; rm -rf autobot-specs; git clone https://github.com/wwqdrh/autobot-specs.git; mv -f autobot-specs/src/* ./autobot; rm -rf ./autobot-specs")
}
