package entities_test

import (
	"go-rest-user-api/entities"
	"go-rest-user-api/utils"
	"testing"
)

const (
	ALL_FIELDS_INVALID  = 3
	NONE_FIELDS_INVALID = 0
)

func TestUser(t *testing.T) {
	t.Run("should return that all fields are invalid", func(t *testing.T) {
		invalidUsers := []entities.User{
			{
				FirstName: "A",
				LastName:  "A",
				Biography: "0123456789012345678",
			},
			{
				FirstName: "Ç",
				LastName:  "Ç",
				Biography: "áááááááááá",
			},
			{
				FirstName: "012345678901234567891",
				LastName:  "012345678901234567891",
				Biography: "Um usuário administrador em um sistema é um perfil com privilégios elevados que vão além das permissões de um usuário comum. Ele possui acesso total para gerenciar configurações globais, criar, editar ou excluir contas de outros usuários, definir e modificar permissões de acesso, visualizar relatórios e registros sensíveis (logs), e realizar ações críticas como backups, restauração de dados e atualizações do sistema. Além disso, o administrador pode configurar integrações com sistemas externos, controlar fluxos de trabalho, e aprovar ou rejeitar solicitações que impactam o funcionamento geral da plataforma. Frequentemente, ele também tem acesso irrestrito a todos os módulos, inclusive aqueles restritos a operações sensíveis, como finanças, segurança e infraestrutura. Suas ações são monitoradas e auditáveis, garantindo que haja controle e rastreabilidade sobre as alterações realizadas no ambiente. Esse papel exige responsabilidade, pois qualquer mudança efetuada pode afetar o funcionamento e a segurança de todo o sistema.",
			},
		}

		for _, user := range invalidUsers {
			utils.Assert(t, ALL_FIELDS_INVALID, len(user.HasAnyFieldInvalid()))
		}
	})

	t.Run("should return that all fields are valid", func(t *testing.T) {
		invalidUsers := []entities.User{
			{
				FirstName: "Ad",
				LastName:  "Ad",
				Biography: "01234567890123456789",
			},
			{
				FirstName: "01234567890123456789",
				LastName:  "01234567890123456789",
				Biography: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Suspendisse turpis nisl, fermentum non convallis vitae, varius vel purus. Morbi vel imperdiet erat. Nam fringilla feugiat bibendum. Fusce eget magna accumsan quam aliquam convallis. Duis tempor dolor ipsum, nec accumsan nisl eleifend et. Etiam quis massa eu libero sollicitudin lacinia. Pellentesque tellus nisl, facilisis id purus eget, vulputate rutrum dolor. Quisque non tempor est posuerç.",
			},
		}

		for _, user := range invalidUsers {
			utils.Assert(t, NONE_FIELDS_INVALID, len(user.HasAnyFieldInvalid()))
		}
	})
}
