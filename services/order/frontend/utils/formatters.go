package utils

import (
	"fmt"
	"strings"
)

func Fullname(lastname, firstname, middlename string) string {
	return fmt.Sprintf(
		"%s %s %s",
		lastname,
		firstname,
		middlename,
	)
}

func NameInitials(fullname string) string {
    fullname = strings.TrimSpace(fullname)
    
    if fullname == "" {
        return ""
    }
    
    splittedName := strings.Split(fullname, " ")
    
    if len(splittedName) < 3 {
        return fullname
    }
    
    lastName := splittedName[0]
    firstName := splittedName[1]
    middleName := splittedName[2]
    
    if lastName == "" || firstName == "" || middleName == "" {
        return fullname
    }
    
    firstNameRunes := []rune(firstName)
    middleNameRunes := []rune(middleName)
    
    if len(firstNameRunes) == 0 || len(middleNameRunes) == 0 {
        return fullname
    }
    
    return fmt.Sprintf("%s %c.%c.", lastName, firstNameRunes[0], middleNameRunes[0])
}