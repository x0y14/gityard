package security_test

import (
	"github.com/stretchr/testify/assert"
	"gityard-api/security"
	"testing"
)

func TestParseSSHKey(t *testing.T) {
	tests := []struct {
		name              string
		fulltext          string
		expectedAlgorithm string
		expectedKeyBody   string
		expectedComment   string
	}{
		{ // `ssh-keygen -t ed25519 -C "your_email@example.com"`: SHA256:C7gepW1lz1Nue4kHaWYrl/8/n2tyyFcpwXPfsZwqrj4
			"ed25519 with comment valid",
			"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIH0VbS/o1S7rE5UMgGO1Ze536UZVakIM39nsKKavJWMl your_email@example.com",
			"ssh-ed25519",
			"AAAAC3NzaC1lZDI1NTE5AAAAIH0VbS/o1S7rE5UMgGO1Ze536UZVakIM39nsKKavJWMl",
			"your_email@example.com",
		},
		{ // `ssh-keygen -t rsa  -C "your_email@example.com"`: SHA256:nkAUuttGxv/XVyMt4T9zWzjJENmxVkChsaddKuxFNcA
			"rsa with comment valid",
			"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDilUiBt0pQAvFyBLJnfQiIi5j7iVWTMJNxv0LyzbbYssd7B40c2B5/cS63hqMhtCuZza5UcygYNs6a0QkcHrmOGnpyXanP/qt/glUmkUvZqbnVlkRwl/uqvGOcM6C2vwAKESCSKyuWNzZ0A2lVqWFje8f6KYanwAVxGKwSUBGSJknqH4Qto+zmq6OvHAxP2MWjizomhUbQ+dZ9La64hVpVJ2W/Nd6NNTqe/XsnaLlzQ3CNfHN+J8dXO3mjwzKgf7ZRpNf/HTVDfI4C1SPU5X43neERbqkdQN6Q3WB955gu85K5KcZuNwApYYtZuvYlNCc05sCDaht8UT/dR6RtnIv6H4/hSqks7FQuCH8IIbXigfBP0YHmlmxm6L/hI2vIrB1EyTPj1CizpCGQwOxSP47+ZTldZadVC4/jdltLmZgjYQK4T1EUCGLBXM2eF+ILiGQrw/vgTMtxQMGX3JeeluxsLfj8U6b9XJd9UGsrmDc8reoYoLoT/rad0FTkxZcyChs= your_email@example.com",
			"ssh-rsa",
			"AAAAB3NzaC1yc2EAAAADAQABAAABgQDilUiBt0pQAvFyBLJnfQiIi5j7iVWTMJNxv0LyzbbYssd7B40c2B5/cS63hqMhtCuZza5UcygYNs6a0QkcHrmOGnpyXanP/qt/glUmkUvZqbnVlkRwl/uqvGOcM6C2vwAKESCSKyuWNzZ0A2lVqWFje8f6KYanwAVxGKwSUBGSJknqH4Qto+zmq6OvHAxP2MWjizomhUbQ+dZ9La64hVpVJ2W/Nd6NNTqe/XsnaLlzQ3CNfHN+J8dXO3mjwzKgf7ZRpNf/HTVDfI4C1SPU5X43neERbqkdQN6Q3WB955gu85K5KcZuNwApYYtZuvYlNCc05sCDaht8UT/dR6RtnIv6H4/hSqks7FQuCH8IIbXigfBP0YHmlmxm6L/hI2vIrB1EyTPj1CizpCGQwOxSP47+ZTldZadVC4/jdltLmZgjYQK4T1EUCGLBXM2eF+ILiGQrw/vgTMtxQMGX3JeeluxsLfj8U6b9XJd9UGsrmDc8reoYoLoT/rad0FTkxZcyChs=",
			"your_email@example.com",
		},
		{ // `ssh-keygen -t ed25519 -C "" `: SHA256:piKP5RDDLJ5oIxj/2xyd2Nlbo3/iokkEdsFMTlKHDNQ
			"ed25519 without comment valid",
			"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOCNKTGWWOng3vhpqK2F9y7iUBQbUGetixPwsrhvBnnt",
			"ssh-ed25519",
			"AAAAC3NzaC1lZDI1NTE5AAAAIOCNKTGWWOng3vhpqK2F9y7iUBQbUGetixPwsrhvBnnt",
			"",
		},
		{ // `ssh-keygen -t rsa -C ""`: SHA256:hqcu1KPcNXovozWP1FlF+g5vJRGHwZvqB7v1yxuait8
			"rsa without comment valid",
			"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC54sbCgv/lwzlEPu5XIDlQ+Vu/wQXbPJLarDEgBNlXI7ghxr8G29sJe3VPOcyDQeNIZr5SYygBo31grZIpp+YCAp3D6ezKkZG/QZrY9olTDX6Mp0STqcgEZa83sHWtE9BEQ9egX0I0CVDX6kv+/OwVQDftrPvwNTC+M1feByFhSNwurFS5IPcjUZfF6krtkfbe8eP7WCL0+9AJvoW9rOTTPFgN0jzZQNHNfjf6kr5JNDkUCfUh8+fYO2V3XhZQ1gaZYlVVWOIACE19KFDwVHgYsHFJAOLW4StmjNUTQopgttqVwe3xT5U0y57jiv6JC9v4gTNdsFel1YzC7G5OFRdlKyhCRSfhMdAyESw7WiyimgTyUMqLGR+Old4vquxodtl+mpRB2ONY6oeV8vUL+rwWyvUI3TwHevxSMq3Qq6PUED3DclgDMXcqhx2qbQ3SsxdkkOZVL9GUjBOl3tsPkpiF3anc0IYZtNVWT2lSkwYB5IHd/4roAPnfiAnpvw2kZ1U=",
			"ssh-rsa",
			"AAAAB3NzaC1yc2EAAAADAQABAAABgQC54sbCgv/lwzlEPu5XIDlQ+Vu/wQXbPJLarDEgBNlXI7ghxr8G29sJe3VPOcyDQeNIZr5SYygBo31grZIpp+YCAp3D6ezKkZG/QZrY9olTDX6Mp0STqcgEZa83sHWtE9BEQ9egX0I0CVDX6kv+/OwVQDftrPvwNTC+M1feByFhSNwurFS5IPcjUZfF6krtkfbe8eP7WCL0+9AJvoW9rOTTPFgN0jzZQNHNfjf6kr5JNDkUCfUh8+fYO2V3XhZQ1gaZYlVVWOIACE19KFDwVHgYsHFJAOLW4StmjNUTQopgttqVwe3xT5U0y57jiv6JC9v4gTNdsFel1YzC7G5OFRdlKyhCRSfhMdAyESw7WiyimgTyUMqLGR+Old4vquxodtl+mpRB2ONY6oeV8vUL+rwWyvUI3TwHevxSMq3Qq6PUED3DclgDMXcqhx2qbQ3SsxdkkOZVL9GUjBOl3tsPkpiF3anc0IYZtNVWT2lSkwYB5IHd/4roAPnfiAnpvw2kZ1U=",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			algo, body, comment, err := security.ParseSSHKey(tt.fulltext)
			//pk, comment, _, _, err := ssh.ParseAuthorizedKey([]byte(tt.fulltext))
			assert.Nil(t, err)
			assert.Equal(t, tt.expectedAlgorithm, algo)
			assert.Equal(t, tt.expectedKeyBody, body)
			assert.Equal(t, tt.expectedComment, comment)
		})
	}
}
