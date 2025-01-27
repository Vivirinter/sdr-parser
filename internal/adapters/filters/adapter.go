package filters

import (
    "fmt"
    "github.com/Vivirinter/sdr-parser/internal/domain"
    "github.com/Vivirinter/sdr-parser/internal/ports"
    "github.com/Vivirinter/sdr-parser/pkg/filter"
)

// FilterAdapter handles conversion between domain types and filter package types
type FilterAdapter struct {
    filter     filter.Filter
    filterType filter.FilterType
    config     filter.FilterConfig
    factory    *filter.FilterFactory
}

func NewFilterAdapter(filterType string) (ports.SignalFilter, error) {
    ft := filter.FilterType(filterType)
    if err := ft.Validate(); err != nil {
        return nil, fmt.Errorf("invalid filter type: %w", err)
    }

    return &FilterAdapter{
        filterType: ft,
        factory:    filter.NewFilterFactory(),
    }, nil
}

func (a *FilterAdapter) validateParams(params map[string]interface{}) error {
    var err error

    getFloat64 := func(params map[string]interface{}, key string) (float64, error) {
        value, ok := params[key]
        if !ok {
            return 0, fmt.Errorf("missing parameter: %s", key)
        }
        floatVal, ok := value.(float64)
        if !ok {
            return 0, fmt.Errorf("parameter %s must be a number", key)
        }
        return floatVal, nil
    }

    getInt := func(params map[string]interface{}, key string) (int, error) {
        value, ok := params[key]
        if !ok {
            return 0, fmt.Errorf("missing parameter: %s", key)
        }
        if floatVal, ok := value.(float64); ok {
            return int(floatVal), nil
        }
        intVal, ok := value.(int)
        if !ok {
            return 0, fmt.Errorf("parameter %s must be an integer", key)
        }
        return intVal, nil
    }

    config := filter.FilterConfig{
        Type: a.filterType,
    }

    switch a.filterType {
    case filter.MovingAverage, filter.Median:
        if config.WindowSize, err = getInt(params, "window_size"); err != nil {
            return err
        }
    case filter.Butterworth:
        if config.CutoffFreq, err = getFloat64(params, "cutoff_freq"); err != nil {
            return err
        }
        if config.Order, err = getInt(params, "order"); err != nil {
            return err
        }
    }

    a.config = config
    return nil
}

func (a *FilterAdapter) Configure(params map[string]interface{}) error {
    if err := a.validateParams(params); err != nil {
        return fmt.Errorf("invalid parameters: %w", err)
    }

    f, err := a.factory.CreateFilter(a.filterType, a.config)
    if err != nil {
        return fmt.Errorf("failed to create filter: %w", err)
    }

    a.filter = f
    return nil
}

func (a *FilterAdapter) Filter(signal *domain.Signal) (*domain.Signal, error) {
    if a.filter == nil {
        return nil, fmt.Errorf("filter not configured")
    }

    filtered, err := a.filter.Process(signal.Samples)
    if err != nil {
        return nil, fmt.Errorf("filter processing failed: %w", err)
    }

    return &domain.Signal{
        Samples:    filtered,
        SampleRate: signal.SampleRate,
    }, nil
}

func (a *FilterAdapter) GetFilterType() string {
    return string(a.filterType)
}

func (a *FilterAdapter) GetConfig() filter.FilterConfig {
    return a.config
}
