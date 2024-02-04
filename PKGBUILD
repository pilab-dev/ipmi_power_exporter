# Create a directory for the AUR package
mkdir -p ~/aur/ipmi_exporter
cd ~/aur/ipmi_exporter

# Create the PKGBUILD file
touch PKGBUILD# Maintainer: Your Name <your.email@example.com>
pkgname=ipmi_exporter
pkgver=0.1.0
pkgrel=1
pkgdesc="IPMI Exporter for Prometheus"
arch=('any')
url="https://github.com/yourusername/ipmi_exporter"
license=('MIT')
depends=('ipmitool')
makedepends=('go')
source=("$pkgname::git+https://github.com/yourusername/$pkgname.git")
sha256sums=('SKIP')

pkgver() {
  cd "$srcdir/$pkgname"
  printf "0.1.0-%s" "$(git rev-parse --short HEAD)"
}

build() {
  cd "$srcdir/$pkgname"
  go build -o ipmi_exporter
}

package() {
  install -Dm755 "$srcdir/$pkgname/ipmi_exporter" "$pkgdir/usr/bin/ipmi_exporter"
  install -Dm644 "$srcdir/$pkgname/ipmi_exporter.service" "$pkgdir/usr/lib/systemd/system/ipmi_exporter.service"
}