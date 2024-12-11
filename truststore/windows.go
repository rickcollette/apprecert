//go:build windows
// +build windows

package truststore

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"unsafe"

	"github.com/appremon/apprecert/config"
	"golang.org/x/sys/windows"
)

type windowsRootStore uintptr

// Windows crypt32.dll constants and functions.
const (
	X509_ASN_ENCODING         = 0x00000001
	CERT_STORE_ADD_REPLACE_EXISTING = 2
	CERT_FIND_SERIAL_NUMBER    = 5
)

// CRYPT_INTEGER_BLOB structure for serial numbers.
type CRYPT_INTEGER_BLOB struct {
	cbData uint32
	pbData *byte
}

var (
	crypt32                     = windows.NewLazySystemDLL("crypt32.dll")
	procCertOpenSystemStoreW    = crypt32.NewProc("CertOpenSystemStoreW")
	procCertCloseStore          = crypt32.NewProc("CertCloseStore")
	procCertAddEncodedCert      = crypt32.NewProc("CertAddEncodedCertificateToStore")
	procCertFindCertificate     = crypt32.NewProc("CertFindCertificateInStore")
	procCertDeleteCertificate   = crypt32.NewProc("CertDeleteCertificateFromStore")
	procCertFreeCertificateCtx  = crypt32.NewProc("CertFreeCertificateContext")
)

func installWindows(cfg *config.Config) error {
	certPath := filepath.Join(cfg.CAROOT, "rootCA.pem")
	certBytes, err := os.ReadFile(certPath)
	if err != nil {
		return err
	}

	block, _ := pem.Decode(certBytes)
	if block == nil || block.Type != "CERTIFICATE" {
		return fmt.Errorf("failed to decode PEM certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}

	store, err := openWindowsRootStore()
	if err != nil {
		return err
	}
	defer store.close()

	return store.addCert(cert.Raw)
}

func uninstallWindows(cfg *config.Config) error {
	certPath := filepath.Join(cfg.CAROOT, "rootCA.pem")
	certBytes, err := os.ReadFile(certPath)
	if err != nil {
		return err
	}

	block, _ := pem.Decode(certBytes)
	if block == nil || block.Type != "CERTIFICATE" {
		return fmt.Errorf("failed to decode PEM certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}

	store, err := openWindowsRootStore()
	if err != nil {
		return err
	}
	defer store.close()

	return store.deleteCertsWithSerial(cert.SerialNumber)
}

// openWindowsRootStore opens the Windows Root CA certificate store.
func openWindowsRootStore() (windowsRootStore, error) {
	rootStr, err := windows.UTF16PtrFromString("ROOT")
	if err != nil {
		return 0, fmt.Errorf("failed to create UTF16 string: %w", err)
	}
	store, _, err2 := procCertOpenSystemStoreW.Call(0, uintptr(unsafe.Pointer(rootStr)))
	if store == 0 {
		return 0, fmt.Errorf("failed to open Windows root store: %v", err2)
	}
	return windowsRootStore(store), nil
}

// close closes the Windows Root CA certificate store.
func (w windowsRootStore) close() error {
	ret, _, err := procCertCloseStore.Call(uintptr(w), 0)
	if ret == 0 {
		return fmt.Errorf("failed to close Windows root store: %v", err)
	}
	return nil
}

func (w windowsRootStore) addCert(cert []byte) error {
	ret, _, err := procCertAddEncodedCert.Call(
		uintptr(w),
		uintptr(X509_ASN_ENCODING),
		uintptr(unsafe.Pointer(&cert[0])),
		uintptr(len(cert)),
		uintptr(CERT_STORE_ADD_REPLACE_EXISTING),
		0,
	)
	if ret == 0 {
		return fmt.Errorf("failed to add certificate: %v", err)
	}
	return nil
}

func (w windowsRootStore) deleteCertsWithSerial(serialNumber *big.Int) error {
	snBytes := serialNumber.Bytes()
	blob := CRYPT_INTEGER_BLOB{
		cbData: uint32(len(snBytes)),
	}
	if len(snBytes) > 0 {
		blob.pbData = &snBytes[0]
	}

	certContext, _, err := procCertFindCertificate.Call(
		uintptr(w),
		uintptr(X509_ASN_ENCODING),
		0,
		uintptr(CERT_FIND_SERIAL_NUMBER),
		uintptr(unsafe.Pointer(&blob)),
		0,
	)
	if certContext == 0 {
		return fmt.Errorf("certificate not found: %v", err)
	}
	defer procCertFreeCertificateCtx.Call(certContext)

	ret, _, err2 := procCertDeleteCertificate.Call(certContext)
	if ret == 0 {
		return fmt.Errorf("failed to delete certificate: %v", err2)
	}
	return nil
}

