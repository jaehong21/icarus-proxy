package api

type NamespaceDto struct {
	Namespace string `json:"namespace" validate:"required"`
}

type NameDto struct {
	Name string `json:"name" validate:"required"`
}

type NamespaceNameDto struct {
	Namespace string `json:"namespace" validate:"required"`
	Name      string `json:"name" validate:"required"`
}
