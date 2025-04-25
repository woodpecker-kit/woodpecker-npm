package version_check

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSemverVersionMinimumSupport(t *testing.T) {
	// mock SemverVersionMinimumSupport
	type args struct {
		version     string
		lessVersion string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "empty version",
			args: args{
				version:     "",
				lessVersion: "1.0.0",
			},
			wantErr: fmt.Errorf("version is empty, please check"),
		},
		{
			name: "empty lessVersion",
			args: args{
				version:     "1.0.0",
				lessVersion: "",
			},
			wantErr: fmt.Errorf("lessVersion is empty, please check"),
		},
		{
			name: "not support version",
			args: args{
				version:     "Semantic Versioning",
				lessVersion: "1.0.0",
			},
			wantErr: fmt.Errorf("can not parse target version: Semantic Versioning err: Invalid Semantic Version"),
		},
		{
			name: "not support less version",
			args: args{
				version:     "1.0.0",
				lessVersion: "Semantic Versioning",
			},
			wantErr: fmt.Errorf("can not parse less version: Semantic Versioning err: improper constraint: >= Semantic Versioning"),
		},
		{
			name: "less",
			args: args{
				version:     "1.0.0",
				lessVersion: "2.0.0",
			},
			wantErr: fmt.Errorf("semver version: 1.0.0 not support, err: [1.0.0 is less than 2.0.0]"),
		},
		{
			name: "equal",
			args: args{
				version:     "2.0.0",
				lessVersion: "2.0.0",
			},
		},
		{
			name: "greater",
			args: args{
				version:     "2.1.1",
				lessVersion: "2.0.0",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			// do SemverVersionMinimumSupport
			gotErr := SemverVersionMinimumSupport(tc.args.version, tc.args.lessVersion)

			// verify SemverVersionMinimumSupport
			assert.Equal(t, tc.wantErr, gotErr)
			if tc.wantErr != nil {
				return
			}
		})
	}
}

func TestSemverVersionConstraint(t *testing.T) {
	// mock SemverVersionConstraint
	type args struct {
		version        string
		minimumVersion string
		maximumVersion string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "empty version",
			args: args{
				version:        "",
				minimumVersion: "1.0.0",
				maximumVersion: "2.0.0",
			},
			wantErr: fmt.Errorf("version is empty, please check"),
		},
		{
			name: "empty minimumVersion",
			args: args{
				version:        "1.0.0",
				minimumVersion: "",
				maximumVersion: "2.0.0",
			},
			wantErr: fmt.Errorf("minimum version is empty, please check"),
		},
		{
			name: "empty maximumVersion",
			args: args{
				version:        "1.0.0",
				minimumVersion: "1.0.0",
				maximumVersion: "",
			},
			wantErr: fmt.Errorf("maximum version is empty, please check"),
		},
		{
			name: "not support version",
			args: args{
				version:        "Semantic Versioning",
				minimumVersion: "1.0.0",
				maximumVersion: "2.0.0",
			},
			wantErr: fmt.Errorf("can not parse target version: Semantic Versioning err: Invalid Semantic Version"),
		},
		{
			name: "not support minimumVersion",
			args: args{
				version:        "1.0.0",
				minimumVersion: "Semantic Versioning",
				maximumVersion: "2.0.0",
			},
			wantErr: fmt.Errorf("can not parse constraint: <= 2.0.0, >= Semantic Versioning err: improper constraint: <= 2.0.0, >= Semantic Versioning"),
		},
		{
			name: "not support maximumVersion",
			args: args{
				version:        "1.0.0",
				minimumVersion: "1.0.0",
				maximumVersion: "Semantic Versioning",
			},
			wantErr: fmt.Errorf("can not parse constraint: <= Semantic Versioning, >= 1.0.0 err: improper constraint: <= Semantic Versioning, >= 1.0.0"),
		},
		{
			name: "less",
			args: args{
				version:        "1.0.0",
				minimumVersion: "2.0.0",
				maximumVersion: "3.0.0",
			},
			wantErr: fmt.Errorf("semver version: 1.0.0 not support, err: [1.0.0 is less than 2.0.0]"),
		},
		{
			name: "constraint",
			args: args{
				version:        "1.2.0",
				minimumVersion: "1.0.0",
				maximumVersion: "2.0.0",
			},
		},
		{
			name: "greater",
			args: args{
				version:        "2.1.1",
				minimumVersion: "1.0.0",
				maximumVersion: "2.0.0",
			},
			wantErr: fmt.Errorf("semver version: 2.1.1 not support, err: [2.1.1 is greater than 2.0.0]"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			// do SemverVersionConstraint
			gotErr := SemverVersionConstraint(tc.args.version, tc.args.minimumVersion, tc.args.maximumVersion)

			// verify SemverVersionConstraint
			assert.Equal(t, tc.wantErr, gotErr)
			if tc.wantErr != nil {
				return
			}
		})
	}
}

func TestSemverVersionPrereleasePrefix(t *testing.T) {
	// mock SemverVersionPrereleasePrefix
	type args struct {
		version    string
		prerelease string
	}
	tests := []struct {
		name    string
		args    args
		wantRes string
		wantErr error
	}{
		{
			name: "empty version",
			args: args{
				version:    "",
				prerelease: "alpha",
			},
			wantErr: fmt.Errorf("version is empty, please check"),
		},
		{
			name: "not support version",
			args: args{
				version:    "Semantic Versioning",
				prerelease: "alpha",
			},
			wantErr: fmt.Errorf("can not parse semver version: Semantic Versioning err: Invalid Semantic Version"),
		},
		{
			name: "not match prerelease",
			args: args{
				version:    "1.2.3",
				prerelease: "alpha",
			},
			wantErr: fmt.Errorf("semver version want prerelease prefix [ alpha ] not start with [  ], please check"),
		},
		{
			name: "not match prerelease diff",
			args: args{
				version:    "1.2.3-alpha",
				prerelease: "beta",
			},
			wantErr: fmt.Errorf("semver version want prerelease prefix [ beta ] not start with [ alpha ], please check"),
		},
		{
			name: "contain prerelease alpha",
			args: args{
				version:    "1.2.3-alpha",
				prerelease: "alpha",
			},
		},
		{
			name: "contain prerelease complicated",
			args: args{
				version:    "1.2.3-alpha.1123",
				prerelease: "alpha",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			// do SemverVersionPrereleasePrefix
			gotErr := SemverVersionPrereleasePrefix(tc.args.version, tc.args.prerelease)

			// verify SemverVersionPrereleasePrefix
			assert.Equal(t, tc.wantErr, gotErr)
			if tc.wantErr != nil {
				return
			}
		})
	}
}

func TestSemverVersionPrereleaseInfoCheck(t *testing.T) {
	// mock SemverVersionPrerelease
	type args struct {
		version    string
		prerelease string
	}
	tests := []struct {
		name    string
		args    args
		wantRes string
		wantErr error
	}{
		{
			name: "empty version",
			args: args{
				version:    "",
				prerelease: "alpha",
			},
			wantErr: fmt.Errorf("version is empty, please check"),
		},
		{
			name: "not support version",
			args: args{
				version:    "Semantic Versioning",
				prerelease: "alpha",
			},
			wantErr: fmt.Errorf("can not parse semver version: Semantic Versioning err: Invalid Semantic Version"),
		},
		{
			name: "not match prerelease",
			args: args{
				version:    "1.2.3",
				prerelease: "alpha",
			},
			wantErr: fmt.Errorf("semver version want prerelease [ alpha ] not match [  ], please check"),
		},
		{
			name: "not match prerelease diff",
			args: args{
				version:    "1.2.3-alpha",
				prerelease: "beta",
			},
			wantErr: fmt.Errorf("semver version want prerelease [ beta ] not match [ alpha ], please check"),
		},
		{
			name: "match prerelease alpha",
			args: args{
				version:    "1.2.3-alpha",
				prerelease: "alpha",
			},
		},
		{
			name: "match prerelease complicated",
			args: args{
				version:    "1.2.3-alpha.1",
				prerelease: "alpha.1",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			// do SemverVersionPrerelease
			gotErr := SemverVersionPrereleaseInfoCheck(tc.args.version, tc.args.prerelease)

			// verify SemverVersionPrerelease
			assert.Equal(t, tc.wantErr, gotErr)
			if tc.wantErr != nil {
				return
			}
		})
	}
}
