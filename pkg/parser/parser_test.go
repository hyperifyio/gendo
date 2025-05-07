// Package parser_test contains test cases for the Gendo script parser.
// It includes comprehensive tests for node definitions, routing instructions,
// and various edge cases in script parsing. The tests verify correct
// handling of different script line formats and error conditions.
package parser

import (
	"reflect"
	"testing"
)

func TestParseLine(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		want   interface{}
		wantOk bool
	}{
		{
			name:   "Empty line",
			input:  "",
			want:   nil,
			wantOk: false,
		},
		{
			name:   "Comment line",
			input:  "# This is a comment",
			want:   nil,
			wantOk: false,
		},
		{
			name:  "Node definition with refs",
			input: "1 : 2 3 4",
			want: &NodeDefinition{
				ID:     1,
				RefIDs: []int{2, 3, 4},
			},
			wantOk: true,
		},
		{
			name:  "Node definition with prompt",
			input: "0 : : Extract the mathematical operation from the text",
			want: &NodeDefinition{
				ID:     0,
				Prompt: "Extract the mathematical operation from the text",
			},
			wantOk: true,
		},
		{
			name:  "Tool definition",
			input: "3 : tool math",
			want: &NodeDefinition{
				ID:     3,
				IsTool: true,
				Tool:   "math",
			},
			wantOk: true,
		},
		{
			name:  "Simple routing",
			input: "3 < 0",
			want: &RouteDefinition{
				Source: 0,
				Dest:   3,
			},
			wantOk: true,
		},
		{
			name:  "Routing with error handler",
			input: "2 ! 3 < 0 calculate 1 + 1",
			want: &RouteDefinition{
				Source:    0,
				Dest:      3,
				ErrorDest: 2,
				Input:     "calculate 1 + 1",
			},
			wantOk: true,
		},
		{
			name:  "Default error handler",
			input: "2 !",
			want: &RouteDefinition{
				ErrorDest: 2,
			},
			wantOk: true,
		},
		{
			name:  "Default output handler",
			input: "3 <",
			want: &RouteDefinition{
				Dest: 3,
			},
			wantOk: true,
		},
		{
			name:   "Invalid node ID",
			input:  "abc : some content",
			want:   nil,
			wantOk: false,
		},
		{
			name:  "Node with prompt containing colons",
			input: "1 : : Format result: add prefix and suffix: done",
			want: &NodeDefinition{
				ID:     1,
				Prompt: "Format result: add prefix and suffix: done",
			},
			wantOk: true,
		},
		{
			name:  "Node with prompt containing multiple colons",
			input: "1 : : Format result: add prefix and suffix: done",
			want: &NodeDefinition{
				ID:     1,
				Prompt: "Format result: add prefix and suffix: done",
			},
			wantOk: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotOk := ParseLine(tt.input)
			if gotOk != tt.wantOk {
				t.Errorf("ParseLine() ok = %v, want %v", gotOk, tt.wantOk)
				return
			}
			if !tt.wantOk {
				return
			}

			switch want := tt.want.(type) {
			case *NodeDefinition:
				got, ok := got.(*NodeDefinition)
				if !ok {
					t.Errorf("ParseLine() returned %T, want *NodeDefinition", got)
					return
				}
				if got.ID != want.ID {
					t.Errorf("ParseLine() ID = %v, want %v", got.ID, want.ID)
				}
				if !reflect.DeepEqual(got.RefIDs, want.RefIDs) {
					t.Errorf("ParseLine() RefIDs = %v, want %v", got.RefIDs, want.RefIDs)
				}
				if got.IsTool != want.IsTool {
					t.Errorf("ParseLine() IsTool = %v, want %v", got.IsTool, want.IsTool)
				}
				if got.Tool != want.Tool {
					t.Errorf("ParseLine() Tool = %v, want %v", got.Tool, want.Tool)
				}
				if got.Prompt != want.Prompt {
					t.Errorf("ParseLine() Prompt = %v, want %v", got.Prompt, want.Prompt)
				}
			case *RouteDefinition:
				got, ok := got.(*RouteDefinition)
				if !ok {
					t.Errorf("ParseLine() returned %T, want *RouteDefinition", got)
					return
				}
				if got.Source != want.Source {
					t.Errorf("ParseLine() Source = %v, want %v", got.Source, want.Source)
				}
				if got.Dest != want.Dest {
					t.Errorf("ParseLine() Dest = %v, want %v", got.Dest, want.Dest)
				}
				if got.ErrorDest != want.ErrorDest {
					t.Errorf("ParseLine() ErrorDest = %v, want %v", got.ErrorDest, want.ErrorDest)
				}
				if got.Input != want.Input {
					t.Errorf("ParseLine() Input = %v, want %v", got.Input, want.Input)
				}
			default:
				t.Errorf("Unknown type in test case: %T", want)
			}
		})
	}
}
