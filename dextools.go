package dextools

import (
	"github.com/barluscuda/dextools/envtools"
	"github.com/barluscuda/dextools/minecraft"
	"github.com/barluscuda/dextools/wenova"
)

func Env() envtools.EnvTools {
	return envtools.EnvTools{}
}

func MinecraftBE() minecraft.Bedrock {
	return minecraft.Bedrock{}
}

func WenovaAPI(token string) wenova.Wenova {
	return wenova.NewWenovaAPI(token)
}
