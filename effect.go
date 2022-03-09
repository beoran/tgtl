package tgtl

// Flow deterimes how the flow of execution is
// affected by a command
type Flow int

// No effect
const NormalFlow Flow = 0

// Breaks out of the current block
const BreakFlow Flow = 1

// Breaks out of the current command
const ReturnFlow Flow = 2

// Error, breaks until rescue block is ofound
const FailFlow Flow = 4

// Every tgtl command evaluates to a value, which is the result
// of the command itself. But if this value also implements the Effect
// interface, then that value also has an effect on the flow of evaluation
// itself.
// Otherwise, if the value does not implement effect,
// this simply means "continue to the next command".
// But other effects may cause the flow of execution to change
// as per the Flow member.
// The unwrap member returns the Value that the effect was carrying wrapped in it
// and which is unwrapped when the effect has influenced the flow.
type Effect interface {
	Flow() Flow
	Unwrap() Value
}
