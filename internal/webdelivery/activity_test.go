package webdelivery

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/iwonz/courier/internal/delivery"
	"github.com/iwonz/courier/internal/fsx"
)

func activityReceipt() delivery.UploadReceipt {
	return delivery.UploadReceipt{
		SchemaVersion: delivery.UploadReceiptSchemaVersion,
		ID:            "00000000-0000-4000-8000-000000000091",
		ServerID:      testServerID,
		DeliveryID:    testDeliveryID,
		Route:         delivery.RouteWebToPath,
		Destination:   "./incoming",
		OriginalName:  "report.pdf",
		StoredName:    "report.pdf",
		Bytes:         7,
		SHA256:        strings.Repeat("a", 64),
		Auth:          delivery.AuthNone,
		Subject:       "anonymous",
		Peer:          "192.0.2.1",
		ReceivedAt:    webTestTime,
	}
}

func TestDestinationReceiptValidationAndDirectories(t *testing.T) {
	receipt := activityReceipt()
	if err := writeDestinationReceipt(context.Background(), nil, receipt); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("nil resource=%v", err)
	}
	if err := writeDestinationReceipt(context.Background(), &Resource{Backend: fsx.Local{}, Path: t.TempDir()}, delivery.UploadReceipt{}); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("invalid receipt=%v", err)
	}
	want := errors.New("failure")
	backend := overrideBackend{Backend: fsx.Local{}, lstat: func(string) (fs.FileInfo, error) { return nil, want }}
	if err := ensureActivityDirectory(backend, "ignored"); !errors.Is(err, want) {
		t.Fatalf("lstat=%v", err)
	}
	root := t.TempDir()
	regular := filepath.Join(root, "regular")
	if err := os.WriteFile(regular, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := ensureActivityDirectory(fsx.Local{}, regular); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("regular directory=%v", err)
	}
	if err := ensureActivityDirectory(fsx.Local{}, filepath.Join(root, "created")); err != nil {
		t.Fatalf("created directory=%v", err)
	}
	backend = overrideBackend{Backend: fsx.Local{}, mkdir: func(string, fs.FileMode) error { return want }}
	if err := ensureActivityDirectory(backend, filepath.Join(root, "mkdir-failure")); !errors.Is(err, want) {
		t.Fatalf("mkdir=%v", err)
	}
	blocked := filepath.Join(root, "blocked")
	if err := os.WriteFile(blocked, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := ensureActivityDirectory(fsx.Local{}, filepath.Join(blocked, "child")); err == nil {
		t.Fatal("blocked directory creation succeeded")
	}
}

func TestDestinationReceiptWriteFailures(t *testing.T) {
	receipt := activityReceipt()
	newResource := func(t *testing.T, backend fsx.Backend) *Resource {
		t.Helper()
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, delivery.ActivityDirectory, delivery.ActivityRecordsDirectory), 0o700); err != nil {
			t.Fatal(err)
		}
		return &Resource{Backend: backend, Path: root}
	}
	want := errors.New("failure")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := writeDestinationReceipt(ctx, &Resource{Backend: fsx.Local{}, Path: t.TempDir()}, receipt); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled=%v", err)
	}

	resource := newResource(t, fsx.Local{})
	stage := filepath.Join(resource.Path, delivery.ActivityDirectory, delivery.ActivityRecordsDirectory, "."+string(receipt.ID)+".tmp")
	if err := os.WriteFile(stage, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := writeDestinationReceipt(context.Background(), resource, receipt); !errors.Is(err, delivery.ErrDuplicate) {
		t.Fatalf("stage collision=%v", err)
	}
	resource = newResource(t, overrideBackend{Backend: fsx.Local{}, lstat: func(name string) (fs.FileInfo, error) {
		if strings.HasSuffix(name, ".tmp") {
			return nil, want
		}
		return os.Lstat(name)
	}})
	if err := writeDestinationReceipt(context.Background(), resource, receipt); !errors.Is(err, want) {
		t.Fatalf("stage inspection=%v", err)
	}

	for _, test := range []struct {
		name   string
		writer *errorWritable
		commit error
		remove error
	}{
		{"create", nil, nil, nil},
		{"write", &errorWritable{writeErr: want, closeErr: want}, nil, want},
		{"sync", &errorWritable{syncErr: want}, nil, nil},
		{"close", &errorWritable{closeErr: want}, nil, nil},
		{"commit", &errorWritable{}, want, nil},
		{"cleanup", &errorWritable{}, nil, want},
	} {
		t.Run(test.name, func(t *testing.T) {
			backend := overrideBackend{Backend: fsx.Local{}, commit: func(string, string) error { return test.commit }, remove: func(string) error { return test.remove }}
			backend.create = func(string, fs.FileMode) (fsx.Writable, error) {
				if test.writer == nil {
					return nil, want
				}
				return test.writer, nil
			}
			if err := writeDestinationReceipt(context.Background(), newResource(t, backend), receipt); err == nil {
				t.Fatal("expected receipt write failure")
			}
		})
	}

	ctx, cancel = context.WithCancel(context.Background())
	writer := &errorWritable{onClose: cancel}
	backend := overrideBackend{Backend: fsx.Local{}, create: func(string, fs.FileMode) (fsx.Writable, error) { return writer, nil }}
	if err := writeDestinationReceipt(ctx, newResource(t, backend), receipt); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancel after close=%v", err)
	}
}

func TestDestinationReceiptRejectsUnsafeActivityTree(t *testing.T) {
	receipt := activityReceipt()
	rootFile := filepath.Join(t.TempDir(), delivery.ActivityDirectory)
	if err := os.WriteFile(rootFile, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := writeDestinationReceipt(context.Background(), &Resource{Backend: fsx.Local{}, Path: filepath.Dir(rootFile)}, receipt); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("unsafe activity root=%v", err)
	}
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, delivery.ActivityDirectory), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, delivery.ActivityDirectory, delivery.ActivityRecordsDirectory), nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := writeDestinationReceipt(context.Background(), &Resource{Backend: fsx.Local{}, Path: root}, receipt); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("unsafe activity records=%v", err)
	}
}
