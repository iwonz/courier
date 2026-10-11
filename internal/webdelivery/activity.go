package webdelivery

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"

	"github.com/iwonz/courier/internal/delivery"
	"github.com/iwonz/courier/internal/fsx"
)

func writeDestinationReceipt(ctx context.Context, resource *Resource, receipt delivery.UploadReceipt) (resultErr error) {
	if resource == nil || resource.Backend == nil {
		return fmt.Errorf("%w: upload receipt destination is unavailable", delivery.ErrInvalid)
	}
	if err := receipt.Validate(); err != nil {
		return err
	}
	root := resource.Backend.Join(resource.Path, delivery.ActivityDirectory)
	records := resource.Backend.Join(root, delivery.ActivityRecordsDirectory)
	if err := ensureActivityDirectory(resource.Backend, root); err != nil {
		return err
	}
	if err := ensureActivityDirectory(resource.Backend, records); err != nil {
		return err
	}
	var encoded bytes.Buffer
	encoder := json.NewEncoder(&encoded)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(receipt) // UploadReceipt contains only infallible JSON value types.
	base := fmt.Sprintf("%020d-%s.json", receipt.ReceivedAt.UnixNano(), receipt.ID)
	final := resource.Backend.Join(records, base)
	stage := resource.Backend.Join(records, "."+string(receipt.ID)+".tmp")
	defer func() {
		if err := resource.Backend.RemoveAll(stage); err != nil && !errors.Is(err, fs.ErrNotExist) {
			resultErr = errors.Join(resultErr, err)
		}
	}()
	if err := ctx.Err(); err != nil {
		return err
	}
	if _, err := resource.Backend.Lstat(stage); err == nil {
		return fmt.Errorf("%w: upload receipt staging collision", delivery.ErrDuplicate)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	file, err := resource.Backend.Create(stage, 0o600)
	if err != nil {
		return err
	}
	closed := false
	defer func() {
		if !closed {
			resultErr = errors.Join(resultErr, file.Close())
		}
	}()
	if _, err := file.Write(encoded.Bytes()); err != nil {
		return err
	}
	if err := file.Sync(); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	closed = true
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := resource.Backend.CommitAbsent(stage, final); err != nil {
		return err
	}
	return resource.Backend.Chmod(final, 0o600)
}

func ensureActivityDirectory(backend fsx.Backend, name string) error {
	info, err := backend.Lstat(name)
	if errors.Is(err, fs.ErrNotExist) {
		if err := backend.MkdirAll(name, 0o700); err != nil {
			return err
		}
		info, err = backend.Lstat(name)
	}
	if err != nil {
		return err
	}
	if !info.IsDir() || info.Mode()&fs.ModeSymlink != 0 {
		return fmt.Errorf("%w: upload activity path is not a real directory", delivery.ErrInvalid)
	}
	return backend.Chmod(name, 0o700)
}
