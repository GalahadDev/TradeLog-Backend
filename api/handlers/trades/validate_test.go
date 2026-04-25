package trades

import "testing"

func TestValidateScreenshotPaths(t *testing.T) {
	uid := "user-abc-123"

	tests := []struct {
		name    string
		paths   []string
		wantErr bool
	}{
		{"nil es ok", nil, false},
		{"vacío es ok", []string{}, false},
		{"un path válido", []string{uid + "/photo.jpg"}, false},
		{"tres paths válidos", []string{uid + "/a.jpg", uid + "/b.jpg", uid + "/c.jpg"}, false},
		{"cuatro paths — demasiados", []string{uid + "/a.jpg", uid + "/b.jpg", uid + "/c.jpg", uid + "/d.jpg"}, true},
		{"URL legacy https es permitida", []string{"https://example.com/photo.jpg"}, false},
		{"mix https legacy + path válido", []string{"https://old.com/x.jpg", uid + "/new.jpg"}, false},
		{"prefijo de otro usuario", []string{"other-user/photo.jpg"}, true},
		{"mix válido e inválido", []string{uid + "/a.jpg", "other/b.jpg"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateScreenshotPaths(tt.paths, uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateScreenshotPaths() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
