# Maintainer: Nev3r Kn0wn <nev3rkn0wn@pm.me>
pkgname=ipmi_power_exporter
pkgver=0.1.0
pkgrel=1
pkgdesc="IPMI Exporter for Prometheus"
arch=('any')
url="https://github.com/devopshaven/ipmi_power_exporter"
license=('MIT')
depends=('ipmitool')
makedepends=('go')
source=("$pkgname::git+https://github.com/devopshaven/$pkgname.git")
sha256sums=('SKIP')

pkgver() {
  cd "$srcdir/$pkgname"
  printf "0.1.0-%s" "$(git rev-parse --short HEAD)"
}

build() {
  cd "$srcdir/$pkgname"
  go build -o ipmi_power_exporter
}

package() {
  install -Dm755 "$srcdir/$pkgname/ipmi_power_exporter" "$pkgdir/usr/bin/ipmi_power_exporter"
  install -Dm644 "$srcdir/$pkgname/ipmi_power_exporter.service" "$pkgdir/usr/lib/systemd/system/ipmi_power_exporter.service"
}