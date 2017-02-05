package router

import "encoding/json"

func (p *plugin) saveAddFilter(channelID string, f filter) error {
	fs, err := p.loadFilters(channelID)
	if err != nil {
		return err
	}

	fs = append(fs, &f)

	return p.saveFilters(channelID, fs)
}

func (p *plugin) saveDeleteFilter(channelID string, filterID string) error {
	fs, err := p.loadFilters(channelID)
	if err != nil {
		return err
	}

	removedFilters := make([]*filter, len(fs))
	if filterID != "all" {
		copy(removedFilters, fs)
		for i := range fs {
			if fs[i].ID == filterID {
				removedFilters = append(removedFilters[:i], removedFilters[i+1:]...)
			}
		}
	}

	return p.saveFilters(channelID, removedFilters)
}

func (p *plugin) loadFilters(channelID string) ([]*filter, error) {
	data, err := p.repository.Load(p.repositoryKey(channelID))
	if err != nil {
		return nil, err
	}
	fs := []*filter{}
	if data != "" {
		if err := json.Unmarshal([]byte(data), &fs); err != nil {
			return nil, err
		}
	}
	return fs, nil
}

func (p *plugin) saveFilters(channelID string, fs []*filter) error {
	writeData, err := json.Marshal(&fs)
	if err != nil {
		return err
	}
	if err := p.repository.Save(p.repositoryKey(channelID), string(writeData)); err != nil {
		return err
	}
	return nil
}
