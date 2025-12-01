package rules

import "github.com/danecwalker/otari/internal/definition"

type RuleError struct {
	Message string
}

func (e *RuleError) Error() string {
	return e.Message
}

type Rule interface {
	Validate(s *definition.Stack) []*RuleError
}

type RuleFunc func(s *definition.Stack) []*RuleError

func (f RuleFunc) Validate(s *definition.Stack) []*RuleError {
	return f(s)
}

func GetDefaultRules() []Rule {
	return []Rule{
		RuleFunc(ValidateContainerNames),
		RuleFunc(ValidateDuplicateEnvironmentVariables),
		RuleFunc(ValidateNetworkNames),
		RuleFunc(ValidateVolumeNames),
		RuleFunc(ValidateContainerNetworkExistence),
		RuleFunc(ValidateContainerVolumeExistence),
		RuleFunc(ValidatePortConflicts),
		RuleFunc(ValidateHostNetworkPortConflicts),
		RuleFunc(ValidateDuplicateVolumeMountsPerContainer),
		RuleFunc(ValidateDependencyExistence),
		RuleFunc(ValidateCircularDependencies),
	}
}

func Validate(stack *definition.Stack) []*RuleError {
	var allErrors []*RuleError
	rules := GetDefaultRules()
	for _, rule := range rules {
		errors := rule.Validate(stack)
		allErrors = append(allErrors, errors...)
	}
	return allErrors
}
