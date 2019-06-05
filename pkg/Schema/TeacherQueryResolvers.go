package Schema

import (
	"github.com/emendoza/classmanager/pkg/Auth"
	"github.com/emendoza/classmanager/pkg/Models"
	"github.com/graphql-go/graphql"
	"log"
)

var selectTeacherQuery = `
SELECT role, username, email
FROM users
WHERE id=$1;
`

var selectStudentFromClassStudentQuery = `
SELECT class_student.student_id, users.role, users.username, users.email
FROM class_student
INNER JOIN users
ON class_student.student_id=users.id
WHERE class_id=$1;
`

var listClassesByTeacher = func(params graphql.ResolveParams) (interface{}, error) {
	token := params.Context.Value("token").(string)
	if !Auth.VerifyToken(token, Models.Teacher) {
		return nil, permissionDenied
	}

	var classes []Models.Class

	rows, err := db.Query(`SELECT id, class_id FROM classes WHERE teacher_id=$1`,
		params.Args["teacherId"].(int))
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		var class Models.Class

		if err := rows.Scan(&class.ID, &class.ClassID); err != nil {
			log.Println(err)
		}

		{
			var teacher Models.User
			teacher.ID = params.Args["teacherId"].(int64)
			err := db.QueryRow(selectTeacherQuery, teacher.ID).Scan(&teacher.Role, &teacher.Username, &teacher.Email)
			if err != nil {
				log.Println(err)
			}
			class.Teacher = teacher
		}

		studentRows, err := db.Query(selectStudentFromClassStudentQuery, class.ID)
		if err != nil {
			log.Println(err)
		}

		for studentRows.Next() {
			var student Models.User

			if err := studentRows.Scan(&student.ID, &student.Role, &student.Username, &student.Email); err != nil {
				log.Println(err)
			}

			class.Students = append(class.Students, student)
		}
		classes = append(classes, class)
	}
	return classes, nil
}


