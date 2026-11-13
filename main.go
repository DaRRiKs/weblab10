package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Student struct {
	ID     int
	Name   string
	Group  string
	Grades []int
}

func CalculateAverage(grades []int) float64 {
	if len(grades) == 0 {
		return 0
	}
	var total float64
	for _, grade := range grades {
		total += float64(grade)
	}
	return total / float64(len(grades))
}

func AddGrade(student *Student, grade int) error {
	if student == nil {
		return fmt.Errorf("студент не найден")
	}
	if grade < 0 || grade > 100 {
		return fmt.Errorf("оценка должна быть от 0 до 100")
	}
	student.Grades = append(student.Grades, grade)
	return nil
}

func UpdateGrade(student *Student, index int, newGrade int) error {
	if student == nil {
		return fmt.Errorf("студент не найден")
	}
	if index < 0 || index >= len(student.Grades) {
		return fmt.Errorf("оценки с таким номером нет")
	}
	if newGrade < 0 || newGrade > 100 {
		return fmt.Errorf("оценка должна быть от 0 до 100")
	}
	student.Grades[index] = newGrade
	return nil
}

func FindStudent(students []Student, indices map[int]int, id int) *Student {
	index, exists := indices[id]
	if !exists {
		return nil
	}
	return &students[index]
}

func PerformanceLevel(average float64) string {
	thresholds := [5]int{0, 50, 70, 85, 95}
	levels := [5]string{
		"Неудовлетворительно",
		"Удовлетворительно",
		"Хорошо",
		"Отлично",
		"Превосходно",
	}
	for i := len(thresholds) - 1; i >= 0; i-- {
		if average >= float64(thresholds[i]) {
			return levels[i]
		}
	}
	return levels[0]
}

func PrintStudent(student *Student) {
	fmt.Printf("ID: %d | ФИО: %s | Группа: %s\n", student.ID, student.Name, student.Group)
	if len(student.Grades) == 0 {
		fmt.Println("Оценки: нет | Средний балл: нет оценок | Уровень: не определён")
		return
	}
	average := CalculateAverage(student.Grades)
	fmt.Printf("Оценки: %v | Средний балл: %.2f | Уровень: %s\n",
		student.Grades, average, PerformanceLevel(average))
}

func PrintStatistics(students []Student) {
	fmt.Printf("Всего студентов: %d\n", len(students))
	if len(students) == 0 {
		fmt.Println("Список студентов пуст.")
		return
	}
	groupCounts := make(map[string]int)
	var totalAverage float64
	gradedCount, below50 := 0, 0
	var best *Student
	var bestAverage float64
	for i, student := range students {
		groupCounts[student.Group]++
		if len(student.Grades) == 0 {
			continue
		}
		average := CalculateAverage(student.Grades)
		totalAverage += average
		gradedCount++
		if average < 50 {
			below50++
		}
		if best == nil || average > bestAverage {
			best = &students[i]
			bestAverage = average
		}
	}
	if gradedCount == 0 {
		fmt.Println("Средний балл группы: нет оценок")
		fmt.Println("Лучший студент: невозможно определить (нет оценок).")
	} else {
		fmt.Printf("Средний балл группы: %.2f\n", totalAverage/float64(gradedCount))
		fmt.Printf("Лучший студент: ID %d, %s, группа %s, средний балл %.2f\n",
			best.ID, best.Name, best.Group, bestAverage)
	}
	fmt.Printf("Студентов со средним баллом ниже 50: %d\n", below50)
	fmt.Println("Количество студентов по академическим группам:")
	groups := make([]string, 0, len(groupCounts))
	for group := range groupCounts {
		groups = append(groups, group)
	}
	sort.Strings(groups)
	for _, group := range groups {
		fmt.Printf("%s: %d\n", group, groupCounts[group])
	}
}

func ReadLine(scanner *bufio.Scanner, prompt string) (string, bool) {
	fmt.Print(prompt)
	if !scanner.Scan() {
		return "", false
	}
	return strings.TrimSpace(scanner.Text()), true
}

func ReadInt(scanner *bufio.Scanner, prompt string) (int, bool) {
	for {
		line, ok := ReadLine(scanner, prompt)
		if !ok {
			return 0, false
		}
		value, err := strconv.Atoi(line)
		if err != nil {
			fmt.Println("Ошибка: введите целое число.")
			continue
		}
		return value, true
	}
}

func ReadText(scanner *bufio.Scanner, prompt string) (string, bool) {
	for {
		value, ok := ReadLine(scanner, prompt)
		if !ok {
			return "", false
		}
		if value == "" {
			fmt.Println("Ошибка: значение не должно быть пустым.")
			continue
		}
		return value, true
	}
}

func main() {
	students := make([]Student, 0)
	indices := make(map[int]int)
	scanner := bufio.NewScanner(os.Stdin)
	defer func() {
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Ошибка чтения ввода:", err)
		}
	}()
	for {
		fmt.Println("\n===== STUDENT GRADE MANAGER =====")
		fmt.Println("\n1. Добавить студента\n2. Добавить оценку\n3. Показать всех студентов\n4. Найти студента\n5. Изменить оценку\n6. Статистика успеваемости\n0. Выход")
		choice, ok := ReadInt(scanner, "Выберите пункт меню: ")
		if !ok {
			return
		}
		switch choice {
		case 0:
			return
		case 1:
			id, ok := ReadInt(scanner, "ID студента: ")
			if !ok {
				return
			}
			if _, exists := indices[id]; exists {
				fmt.Println("Ошибка: студент с таким ID уже существует.")
				continue
			}
			name, ok := ReadText(scanner, "ФИО: ")
			if !ok {
				return
			}
			group, ok := ReadText(scanner, "Академическая группа: ")
			if !ok {
				return
			}
			students = append(students, Student{ID: id, Name: name, Group: group})
			indices[id] = len(students) - 1
			fmt.Println("Студент добавлен.")
		case 2, 4, 5:
			id, ok := ReadInt(scanner, "ID студента: ")
			if !ok {
				return
			}
			student := FindStudent(students, indices, id)
			if student == nil {
				fmt.Println("Студент не найден.")
				continue
			}
			switch choice {
			case 2:
				grade, ok := ReadInt(scanner, "Оценка (0–100): ")
				if !ok {
					return
				}
				if err := AddGrade(student, grade); err != nil {
					fmt.Println("Ошибка:", err)
				} else {
					fmt.Println("Оценка добавлена.")
				}
			case 4:
				PrintStudent(student)
			case 5:
				if len(student.Grades) == 0 {
					fmt.Println("У студента нет оценок для изменения.")
					continue
				}
				fmt.Println("Оценки студента:")
				for i, grade := range student.Grades {
					fmt.Printf("%d. %d\n", i+1, grade)
				}
				number, ok := ReadInt(scanner, "Номер оценки: ")
				if !ok {
					return
				}
				if number < 1 || number > len(student.Grades) {
					fmt.Println("Ошибка: оценки с таким номером нет.")
					continue
				}
				newGrade, ok := ReadInt(scanner, "Новая оценка (0–100): ")
				if !ok {
					return
				}
				if err := UpdateGrade(student, number-1, newGrade); err != nil {
					fmt.Println("Ошибка:", err)
				} else {
					fmt.Println("Оценка изменена.")
				}
			}
		case 3:
			if len(students) == 0 {
				fmt.Println("Список студентов пуст.")
			}
			for i := range students {
				PrintStudent(&students[i])
			}
		case 6:
			PrintStatistics(students)
		default:
			fmt.Println("Ошибка: выберите пункт меню от 0 до 6.")
		}
	}
}
