package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func DeleteStudentByID(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Authorization, ngrok-skip-browser-warning")

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	// Step 1: Nullify student_id in point_of_contact (ON DELETE SET NULL doesn't cascade)
	_, err = db.Exec(`UPDATE point_of_contact SET student_id = NULL WHERE student_id = ?`, id)
	if err != nil {
		http.Error(w, "Failed to update point_of_contact: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 2: Delete from student (cascades through specific_documentation, stu_accom, stu_dis, person, etc.)
	_, err = db.Exec(`DELETE FROM student WHERE id = ?`, id)
	if err != nil {
		http.Error(w, "Failed to delete student: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Student deleted successfully"})
}
