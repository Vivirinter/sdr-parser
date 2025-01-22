package processing

import (
	"math"
)

type AGC struct {
	attackTime  float64
	releaseTime float64
	target      float64
	maxGain     float64
	minGain     float64
	currentGain float64
}

func NewAGC(attackTime, releaseTime, target, maxGain, minGain float64) *AGC {
	return &AGC{
		attackTime:  attackTime,
		releaseTime: releaseTime,
		target:      target,
		maxGain:     maxGain,
		minGain:     minGain,
		currentGain: 1.0,
	}
}

func (a *AGC) Process(samples []float64, sampleRate float64) []float64 {
	output := make([]float64, len(samples))
	
	attackCoeff := math.Exp(-1.0 / (sampleRate * a.attackTime))
	releaseCoeff := math.Exp(-1.0 / (sampleRate * a.releaseTime))
	
	for i, sample := range samples {
		magnitude := math.Abs(sample)
		error := a.target/magnitude - a.currentGain
		
		if error > 0 {
			a.currentGain = attackCoeff*a.currentGain + (1-attackCoeff)*a.target/magnitude
		} else {
			a.currentGain = releaseCoeff*a.currentGain + (1-releaseCoeff)*a.target/magnitude
		}
		
		if a.currentGain > a.maxGain {
			a.currentGain = a.maxGain
		} else if a.currentGain < a.minGain {
			a.currentGain = a.minGain
		}
		
		output[i] = sample * a.currentGain
	}
	
	return output
}

func (a *AGC) GetCurrentGain() float64 {
	return a.currentGain
}

func (a *AGC) GetGainReduction() float64 {
	if a.currentGain > 0 {
		return 1.0 / a.currentGain
	}
	return math.MaxFloat64
}

func (a *AGC) GetCompressionDB() float64 {
	if a.currentGain > 0 {
		return 20 * math.Log10(a.currentGain)
	}
	return -math.MaxFloat64
}

func (a *AGC) Reset() {
	a.currentGain = 1.0
}

func (a *AGC) SetTarget(target float64) {
	a.target = target
}

func (a *AGC) SetAttackTime(attackTime float64) {
	a.attackTime = attackTime
}

func (a *AGC) SetReleaseTime(releaseTime float64) {
	a.releaseTime = releaseTime
}

func (a *AGC) SetGainLimits(minGain, maxGain float64) {
	a.minGain = minGain
	a.maxGain = maxGain
}
