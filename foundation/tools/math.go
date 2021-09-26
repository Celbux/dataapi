package tools

// Min returns the minimum values off its 2 inputs
// The built in Math.Min() functions requires floats for its inputs and output
// which is shitty to work with
func Min(num1 int, num2 int) int {
	if num1 < num2 {
		return num1
	}
	return num2
}
