/*
Package job handles conversions from configuration format into a collection of
jobs that can be executed by the driver in parallel.
*/
package job

import (
	"context"
	"os"
	"path/filepath"

	"github.com/friendly-fhir/fhenix/pkg/model"
	"github.com/friendly-fhir/fhenix/pkg/transform"
)

// Job represents a single job that can be executed by the driver.
type Job struct {
	outputPath string
	input      *input
	transform  *transform.Transform
}

// New creates a collection of jobs from a given transformation, where each
// job corresponds to a different output file-path for a set of input
// matched by the transformation.
func New(model *model.Model, outputPath string, transform *transform.Transform) ([]*Job, error) {
	// inputs is a mapping of output file path to the input types that can be
	// transformed.
	inputs := map[string]*input{}

	for _, t := range model.Types().All() {
		if transform.CanTransform(t) {
			out, err := transform.OutputPath(t)
			if err != nil {
				return nil, err
			}
			if !filepath.IsAbs(out) {
				out = filepath.Join(filepath.FromSlash(outputPath), out)
			}
			if _, ok := inputs[out]; !ok {
				inputs[out] = &input{}
			}
			inputs[out].StructureDefinitions = append(inputs[out].StructureDefinitions, t)
		}
	}
	for _, c := range model.CodeSystems() {
		if transform.CanTransform(c) {
			out, err := transform.OutputPath(c)
			if err != nil {
				return nil, err
			}
			if _, ok := inputs[out]; !ok {
				inputs[out] = &input{}
			}
			inputs[out].CodeSystems = append(inputs[out].CodeSystems, c)
		}
	}

	jobs := make([]*Job, 0, len(inputs))
	for out, in := range inputs {
		jobs = append(jobs, &Job{
			outputPath: out,
			input:      in,
			transform:  transform,
		})
	}

	return jobs, nil
}

// Execute runs the job and writes the output to the specified file.
func (j *Job) Execute(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		break
	}
	if err := os.MkdirAll(filepath.Dir(j.outputPath), 0755); err != nil {
		return err
	}

	file, err := os.Create(j.outputPath)
	if err != nil {
		return err
	}
	defer file.Close()
	return j.transform.Execute(file, j.input)
}

// OutputPath returns the output path for the job.
func (j *Job) OutputPath() string {
	return j.outputPath
}

// StructureDefinitions returns the structure definitions that should be
// transformed by this job.
func (j *Job) StructureDefinitions() []*model.Type {
	return j.input.StructureDefinitions
}

// CodeSystems returns the code systems that should be transformed by this job.
func (j *Job) CodeSystems() []*model.CodeSystem {
	return j.input.CodeSystems
}

// ValueSets returns the value sets that should be transformed by this job.
func (j *Job) ValueSets() []*struct{} {
	return j.input.ValueSets
}

type input struct {
	StructureDefinitions []*model.Type
	CodeSystems          []*model.CodeSystem
	ValueSets            []*struct{}
}
