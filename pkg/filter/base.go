package filter

// BaseFilter provides common functionality for all filters
type BaseFilter struct {
    config FilterConfig
    stats  FilterStats
}

func (f *BaseFilter) Configure(config FilterConfig) error {
    if err := config.Validate(); err != nil {
        return err
    }
    f.config = config
    f.stats = FilterStats{}
    return nil
}

func (f *BaseFilter) GetStats() FilterStats {
    return f.stats
}

func (f *BaseFilter) GetConfig() FilterConfig {
    return f.config
}

func (f *BaseFilter) updateStats(input, output []float64) {
    f.stats.InputSamples += len(input)
    f.stats.OutputSamples += len(output)
    
    var sumIn, sumOut float64
    for _, v := range input {
        sumIn += v
    }
    for _, v := range output {
        sumOut += v
    }
    
    if len(input) > 0 {
        f.stats.InputMean = sumIn / float64(len(input))
    }
    if len(output) > 0 {
        f.stats.OutputMean = sumOut / float64(len(output))
    }
}
