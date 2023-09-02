package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Applicant struct {
	firstName, lastName                                                          string
	ScorePhysics, ScoreChemistry, ScoreMath, ScoreCS, specialScore, displayScore float64
	prefs                                                                        [3]string
}

func getInputs() []Applicant {
	const (
		fileName = "data/applicants.txt"
		noPrefs  = 3
	)
	var a []Applicant
	var prefsArray [noPrefs]string
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ds := 0.00
		line := strings.Split(scanner.Text(), " ")
		lineScorePhysics, err := strconv.ParseFloat(line[2], 64)
		lineScoreChemistry, err := strconv.ParseFloat(line[3], 64)
		lineScoreMath, err := strconv.ParseFloat(line[4], 64)
		lineScoreCS, err := strconv.ParseFloat(line[5], 64)
		lineScoreSpecial, err := strconv.ParseFloat(line[6], 64)
		prefsArray[0] = line[7]
		prefsArray[1] = line[8]
		prefsArray[2] = line[9]
		if err != nil {
			fmt.Println(err)
		}
		// init applicants list
		a = append(a, Applicant{line[0],
			line[1],
			lineScorePhysics,
			lineScoreChemistry,
			lineScoreMath,
			lineScoreCS,
			lineScoreSpecial,
			ds,
			prefsArray,
		})
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	return a
}

func orderByDept(a []Applicant, department string) {
	getMean := func(data []float64) float64 {
		if len(data) == 0 {
			return 0
		}
		var sum float64
		for _, d := range data {
			sum += d
		}
		return sum / float64(len(data))
	}

	var optionalArray [2]float64
	//physics and math for the Physics department,
	//chemistry for the Chemistry department,
	//math for the Mathematics department,
	//computer science and math for the Engineering Department,
	//chemistry and physics for the Biotech department.
	for i := len(a) - 1; i >= 0; i-- {
		switch department {
		case "Biotech":
			optionalArray[0], optionalArray[1] = a[i].ScoreChemistry, a[i].ScorePhysics
			if getMean(optionalArray[:]) >= a[i].specialScore {
				a[i].displayScore = getMean(optionalArray[:])
			} else {
				a[i].displayScore = a[i].specialScore
			}
		case "Chemistry":
			if a[i].ScoreChemistry >= a[i].specialScore {
				a[i].displayScore = a[i].ScoreChemistry
			} else {
				a[i].displayScore = a[i].specialScore
			}
		case "Engineering":
			optionalArray[0], optionalArray[1] = a[i].ScoreCS, a[i].ScoreMath
			if getMean(optionalArray[:]) >= a[i].specialScore {
				a[i].displayScore = getMean(optionalArray[:])
			} else {
				a[i].displayScore = a[i].specialScore
			}
		case "Mathematics":
			if a[i].ScoreMath >= a[i].specialScore {
				a[i].displayScore = a[i].ScoreMath
			} else {
				a[i].displayScore = a[i].specialScore
			}
		case "Physics":
			optionalArray[0], optionalArray[1] = a[i].ScorePhysics, a[i].ScoreMath
			if getMean(optionalArray[:]) >= a[i].specialScore {
				a[i].displayScore = getMean(optionalArray[:])
			} else {
				a[i].displayScore = a[i].specialScore
			}
		}
	}

	sort.SliceStable(a, func(i, j int) bool {
		if a[i].displayScore == a[j].displayScore {
			return a[i].firstName+a[i].lastName > a[j].firstName+a[j].lastName
		}
		return a[i].displayScore < a[j].displayScore
	})
}

func admissionProcess(a []Applicant, N int) map[string][]Applicant {
	var departments = [5]string{"Biotech", "Chemistry", "Engineering", "Mathematics", "Physics"}
	var selectedApplicantsMap = make(map[string][]Applicant)
	// handle a new wave of participants, for each of the 3 prefs
	for wave := 0; wave <= 2; wave++ {
		for _, dept := range departments {
			orderByDept(a, dept)
			for i := len(a) - 1; i >= 0; i-- {
				foundApplicant := false
				if a[i].prefs[wave] == dept {
					foundApplicant = true
				}
				// if found an applicant in this wave, add it to the selected applicants.
				if foundApplicant && len(selectedApplicantsMap[dept]) < N {
					selectedApplicantsMap[dept] = append(selectedApplicantsMap[dept], a[i])
					//also remove from applicants
					a = append(a[:i], a[i+1:]...)
				}
			}
		}
	}
	return selectedApplicantsMap
}

func main() {
	var noApplicants int

	fmt.Scanf("%d", &noApplicants)
	applicants := getInputs()
	resultApplications := admissionProcess(applicants, noApplicants)

	keys := make([]string, 0, len(resultApplications))
	for k := range resultApplications {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		printOutput := applicants[:0]
		for _, rApplicant := range resultApplications[k] {
			printOutput = append(printOutput, rApplicant)
		}
		sort.SliceStable(printOutput, func(i, j int) bool {
			if printOutput[i].displayScore == printOutput[j].displayScore {
				return applicants[i].firstName+applicants[i].lastName < applicants[j].firstName+applicants[j].lastName
			}
			return printOutput[i].displayScore > printOutput[j].displayScore
		})
		file, err := os.Create("output/" + strings.ToLower(k) + ".txt")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		for _, rApplicant := range printOutput {
			fmt.Fprintf(file, "%s %s %.2f\n",
				rApplicant.firstName,
				rApplicant.lastName,
				rApplicant.displayScore,
			)
		}
	}
}
