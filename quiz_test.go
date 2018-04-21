package main

import "testing"

func Test_isAnswerCorrect(t *testing.T) {
	type args struct {
		providedAnswer string
		expectedAnswer string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
	// TODO: Add test cases.
		{"Equal value", args{providedAnswer:"9", expectedAnswer:"9"}, true},
		{"Different value", args{providedAnswer:"9", expectedAnswer:"6"}, false},
		{"Equal with Capital expected", args{providedAnswer:"capital", expectedAnswer:"Capital"}, true},
		{"Equal with Capital provided", args{providedAnswer:"Capital", expectedAnswer:"capital"}, true},
		{"Equal with trailing spaces expected", args{providedAnswer:"spaces", expectedAnswer:"spaces "}, true},
		{"Equal with trailing spaces provided", args{providedAnswer:"spaces ", expectedAnswer:"spaces"}, true},
		{"Equal with starting spaces expected", args{providedAnswer:"spaces", expectedAnswer:" spaces"}, true},
		{"Equal with starting spaces provided", args{providedAnswer:" spaces", expectedAnswer:"spaces"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isAnswerCorrect(tt.args.providedAnswer, tt.args.expectedAnswer); got != tt.want {
				t.Errorf("isAnswerCorrect() = %v, want %v", got, tt.want)
			}
		})
	}
}
