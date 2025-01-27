package filter

import "fmt"

type ConfigValidator interface {
    Validate(config FilterConfig) error
}

type DefaultConfigValidator struct{}

func (v *DefaultConfigValidator) Validate(config FilterConfig) error {
    if config.WindowSize <= 0 {
        return fmt.Errorf("window size must be positive")
    }
    return nil
}

// FilterFactory creates and configures filters based on type and configuration
type FilterFactory struct {
    validator ConfigValidator
    factories map[FilterType]func() Filter
}

func NewFilterFactory() *FilterFactory {
    return &FilterFactory{
        validator: &DefaultConfigValidator{},
        factories: map[FilterType]func() Filter{
            MovingAverage: func() Filter { return NewMovingAverageFilter() },
            Median:        func() Filter { return NewMedianFilter() },
            Butterworth:   func() Filter { return NewButterworthFilter() },
        },
    }
}

func (f *FilterFactory) CreateFilter(filterType FilterType, config FilterConfig) (Filter, error) {
    if err := f.validator.Validate(config); err != nil {
        return nil, fmt.Errorf("invalid configuration: %w", err)
    }

    factory, exists := f.factories[filterType]
    if !exists {
        return nil, fmt.Errorf("unknown filter type: %s", filterType)
    }

    filter := factory()
    if err := filter.Configure(config); err != nil {
        return nil, fmt.Errorf("failed to configure filter: %w", err)
    }

    return filter, nil
}
