package config

import "github.com/riferrei/srclient"

var schemaRegistryClient *srclient.SchemaRegistryClient

func ConfigSchemaRegister(){
	schemaRegistryClient = srclient.CreateSchemaRegistryClient("http://schema-registry:8081")
}

func GetSchemaRegister() *srclient.SchemaRegistryClient{
	return schemaRegistryClient
}