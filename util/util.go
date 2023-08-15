package util

func ChunkSlice[T any](slice []T, chunkSize int) [][]T {
	var chunks [][]T
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

func IsAwsProvider(provider string) bool {
	providers := []string{SsmProvider, SecretsManagerProvider}

	for _, p := range providers {
		if provider == p {
			return true
		}
	}

	return false
}
