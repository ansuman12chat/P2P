//go:generate mockgen -package mock -source internal/app/ioutil.go -destination internal/mock/ioutil.go
//go:generate mockgen -package mock -source internal/app/xdg.go -destination internal/mock/xdg.go
//go:generate mockgen -package mock -source internal/app/time.go -destination internal/mock/time.go
//go:generate mockgen -package mock -source internal/app/network.go -destination internal/mock/network.go
package pcp
