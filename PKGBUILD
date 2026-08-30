# Maintainer: Bloby22
pkgname=andtls
pkgver=1.0.0
pkgrel=1
pkgdesc="TUI application andtls"
arch=('x86_64')
url="https://github.com/Bloby22/andtls"
license=('MIT')
makedepends=('go' 'git')
source=("$pkgname-$pkgver.tar.gz::https://github.com/Bloby22/andtls/archive/refs/tags/v$pkgver.tar.gz")
sha256sums=('SKIP')

build() {
    cd "$pkgname-$pkgver"
    export CGO_ENABLED=0
    export GOFLAGS="-buildmode=pie -trimpath -mod=readonly -modcacherw"
    go build -o "$pkgname" .
}

package() {
    cd "$pkgname-$pkgver"
    install -Dm755 "$pkgname" "$pkgdir/usr/bin/$pkgname"
}
