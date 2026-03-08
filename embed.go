package zentxt

import "embed"

//go:embed all:frontend/dist
var Assets embed.FS

//go:embed migrations/*
var Migrations embed.FS
