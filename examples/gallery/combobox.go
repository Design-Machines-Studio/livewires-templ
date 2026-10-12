package main

import "github.com/Design-Machines-Studio/livewires-templ/form"

func comboboxTimezoneOptions() []form.SelectOption {
	return []form.SelectOption{
		{Value: "UTC", Label: "UTC"},
		{Value: "Asia/Makassar", Label: "Makassar"},
		{Value: "Asia/Manila", Label: "Manila"},
		{Value: "Asia/Chatham", Label: "Chatham Islands"},
		{Value: "Europe/Paris", Label: "Paris"},
		{Value: "America/Sao_Paulo", Label: "São Paulo"},
		{Value: "America/New_York", Label: "New York"},
	}
}
